package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"one_click_2/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetCustomers lists all customers for the current agent's agency
func GetCustomers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT UserID, FirstName, LastName, Email, Role, AgencyID, LocationID, 
			Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, 
			NeedsFinancialFiling, AdditionalData FROM Users WHERE Role = 'Customer' AND AgencyID = $1`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		customers := []models.User{}
		for rows.Next() {
			var customer models.User
			if err := rows.Scan(&customer.UserID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Role,
				&customer.AgencyID, &customer.LocationID, &customer.Address, &customer.LicenseNumber, &customer.LicenseState,
				&customer.DateOfBirth, &customer.MobileNumber, &customer.HasPhysicalImpairment, &customer.NeedsFinancialFiling,
				&customer.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			customers = append(customers, customer)
		}

		c.JSON(http.StatusOK, customers)
	}
}

// CreateCustomer creates a new customer for the current agent's agency
func CreateCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var customer models.User
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a default password if not provided
		defaultPassword := "TempPassword123!"
		if customer.PasswordHash == "" {
			customer.PasswordHash = defaultPassword
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Set role to Customer and agency to current agent's agency
		customer.Role = models.RoleCustomer
		customer.AgencyID = &agencyID

		_, err = db.Exec(`INSERT INTO Users (FirstName, LastName, Email, PasswordHash, Role, AgencyID, LocationID, 
			Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, 
			NeedsFinancialFiling, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
			customer.FirstName, customer.LastName, customer.Email, string(hashedPassword), customer.Role, customer.AgencyID,
			customer.LocationID, customer.Address, customer.LicenseNumber, customer.LicenseState, customer.DateOfBirth,
			customer.MobileNumber, customer.HasPhysicalImpairment, customer.NeedsFinancialFiling, customer.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully"})
	}
}

// UpdateCustomer updates a customer within the current agent's agency
func UpdateCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		customerID := c.Param("id")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var customer models.User
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only if customer belongs to agent's agency
		result, err := db.Exec(`UPDATE Users SET FirstName = $1, LastName = $2, Email = $3, 
			Address = $4, LicenseNumber = $5, LicenseState = $6, DateOfBirth = $7, MobileNumber = $8, 
			HasPhysicalImpairment = $9, NeedsFinancialFiling = $10, AdditionalData = $11 
			WHERE UserID = $12 AND AgencyID = $13 AND Role = 'Customer'`,
			customer.FirstName, customer.LastName, customer.Email, customer.Address, customer.LicenseNumber,
			customer.LicenseState, customer.DateOfBirth, customer.MobileNumber, customer.HasPhysicalImpairment,
			customer.NeedsFinancialFiling, customer.AdditionalData, customerID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
	}
}

// DeleteCustomer deletes (soft delete) a customer within the current agent's agency
func DeleteCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		customerID := c.Param("id")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		// For now, we'll do a hard delete. In production, consider soft delete by adding a deleted_at field
		result, err := db.Exec("DELETE FROM Users WHERE UserID = $1 AND AgencyID = $2 AND Role = 'Customer'", customerID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
	}
}

// GetVehicles lists vehicles for a specific customer (with optional customer ID query param)
func GetVehicles(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		customerIDParam := c.Query("customer_id")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var query string
		var args []interface{}

		if customerIDParam != "" {
			// Verify customer belongs to agent's agency
			var customerAgencyID int
			err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1 AND Role = 'Customer'", customerIDParam).Scan(&customerAgencyID)
			if err != nil || customerAgencyID != agencyID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to customer's vehicles"})
				return
			}

			query = `SELECT VehicleID, CustomerID, VIN, Make, Model, Year, Type, PlateNumber, PlateType, 
				BodyStyle, VehicleUse, VehicleHistory, AnnualMileage, CurrentOdometerReading, 
				OdometerReadingDate, PurchasedInLast90Days, LengthOfOwnership, OwnershipStatus, 
				HasRacingEquipment, HasExistingDamage, CostNew, AdditionalData 
				FROM Vehicles WHERE CustomerID = $1`
			args = []interface{}{customerIDParam}
		} else {
			// Get all vehicles for customers in agent's agency
			query = `SELECT v.VehicleID, v.CustomerID, v.VIN, v.Make, v.Model, v.Year, v.Type, v.PlateNumber, v.PlateType, 
				v.BodyStyle, v.VehicleUse, v.VehicleHistory, v.AnnualMileage, v.CurrentOdometerReading, 
				v.OdometerReadingDate, v.PurchasedInLast90Days, v.LengthOfOwnership, v.OwnershipStatus, 
				v.HasRacingEquipment, v.HasExistingDamage, v.CostNew, v.AdditionalData 
				FROM Vehicles v 
				JOIN Users u ON v.CustomerID = u.UserID 
				WHERE u.AgencyID = $1 AND u.Role = 'Customer'`
			args = []interface{}{agencyID}
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		vehicles := []models.Vehicle{}
		for rows.Next() {
			var vehicle models.Vehicle
			if err := rows.Scan(&vehicle.VehicleID, &vehicle.CustomerID, &vehicle.VIN, &vehicle.Make, &vehicle.Model,
				&vehicle.Year, &vehicle.Type, &vehicle.PlateNumber, &vehicle.PlateType, &vehicle.BodyStyle,
				&vehicle.VehicleUse, &vehicle.VehicleHistory, &vehicle.AnnualMileage, &vehicle.CurrentOdometerReading,
				&vehicle.OdometerReadingDate, &vehicle.PurchasedInLast90Days, &vehicle.LengthOfOwnership,
				&vehicle.OwnershipStatus, &vehicle.HasRacingEquipment, &vehicle.HasExistingDamage, &vehicle.CostNew,
				&vehicle.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			vehicles = append(vehicles, vehicle)
		}

		c.JSON(http.StatusOK, vehicles)
	}
}

// CreateVehicle creates a new vehicle for a customer
func CreateVehicle(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var vehicle models.Vehicle
		if err := c.ShouldBindJSON(&vehicle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify customer belongs to agent's agency
		var customerAgencyID int
		err = db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1 AND Role = 'Customer'", vehicle.CustomerID).Scan(&customerAgencyID)
		if err != nil || customerAgencyID != agencyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to customer"})
			return
		}

		_, err = db.Exec(`INSERT INTO Vehicles (CustomerID, VIN, Make, Model, Year, Type, PlateNumber, PlateType, 
			BodyStyle, VehicleUse, VehicleHistory, AnnualMileage, CurrentOdometerReading, 
			OdometerReadingDate, PurchasedInLast90Days, LengthOfOwnership, OwnershipStatus, 
			HasRacingEquipment, HasExistingDamage, CostNew, AdditionalData) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`,
			vehicle.CustomerID, vehicle.VIN, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.Type,
			vehicle.PlateNumber, vehicle.PlateType, vehicle.BodyStyle, vehicle.VehicleUse, vehicle.VehicleHistory,
			vehicle.AnnualMileage, vehicle.CurrentOdometerReading, vehicle.OdometerReadingDate,
			vehicle.PurchasedInLast90Days, vehicle.LengthOfOwnership, vehicle.OwnershipStatus,
			vehicle.HasRacingEquipment, vehicle.HasExistingDamage, vehicle.CostNew, vehicle.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Vehicle created successfully"})
	}
}

// UpdateVehicle updates a vehicle
func UpdateVehicle(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		vehicleID := c.Param("id")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var vehicle models.Vehicle
		if err := c.ShouldBindJSON(&vehicle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only if vehicle's customer belongs to agent's agency
		result, err := db.Exec(`UPDATE Vehicles SET VIN = $1, Make = $2, Model = $3, Year = $4, Type = $5, 
			PlateNumber = $6, PlateType = $7, BodyStyle = $8, VehicleUse = $9, VehicleHistory = $10, 
			AnnualMileage = $11, CurrentOdometerReading = $12, OdometerReadingDate = $13, 
			PurchasedInLast90Days = $14, LengthOfOwnership = $15, OwnershipStatus = $16, 
			HasRacingEquipment = $17, HasExistingDamage = $18, CostNew = $19, AdditionalData = $20 
			WHERE VehicleID = $21 AND CustomerID IN (SELECT UserID FROM Users WHERE AgencyID = $22 AND Role = 'Customer')`,
			vehicle.VIN, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.Type, vehicle.PlateNumber,
			vehicle.PlateType, vehicle.BodyStyle, vehicle.VehicleUse, vehicle.VehicleHistory,
			vehicle.AnnualMileage, vehicle.CurrentOdometerReading, vehicle.OdometerReadingDate,
			vehicle.PurchasedInLast90Days, vehicle.LengthOfOwnership, vehicle.OwnershipStatus,
			vehicle.HasRacingEquipment, vehicle.HasExistingDamage, vehicle.CostNew, vehicle.AdditionalData,
			vehicleID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Vehicle updated successfully"})
	}
}

// AssignDriver assigns a driver to a vehicle
func AssignDriver(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var vehicleDriver models.VehicleDriver
		if err := c.ShouldBindJSON(&vehicleDriver); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify both vehicle's customer and driver belong to agent's agency
		var vehicleCustomerAgencyID, driverAgencyID int
		err = db.QueryRow(`SELECT u.AgencyID FROM Vehicles v 
			JOIN Users u ON v.CustomerID = u.UserID WHERE v.VehicleID = $1`, vehicleDriver.VehicleID).Scan(&vehicleCustomerAgencyID)
		if err != nil || vehicleCustomerAgencyID != agencyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to vehicle"})
			return
		}

		err = db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", vehicleDriver.UserID).Scan(&driverAgencyID)
		if err != nil || driverAgencyID != agencyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to driver"})
			return
		}

		_, err = db.Exec(`INSERT INTO VehicleDrivers (VehicleID, UserID, DriverType, RelationshipToInsured, 
			Gender, MaritalStatus, LicensingException, NeedsFinancialResponsibilityFiling, 
			HasUncompensatedImpairment, MobileNumber, AdditionalData) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			vehicleDriver.VehicleID, vehicleDriver.UserID, vehicleDriver.DriverType,
			vehicleDriver.RelationshipToInsured, vehicleDriver.Gender, vehicleDriver.MaritalStatus,
			vehicleDriver.LicensingException, vehicleDriver.NeedsFinancialResponsibilityFiling,
			vehicleDriver.HasUncompensatedImpairment, vehicleDriver.MobileNumber, vehicleDriver.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Driver assigned to vehicle successfully"})
	}
}

// AddDrivingHistory adds driving history for a user
func AddDrivingHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var drivingHistory models.DrivingHistory
		if err := c.ShouldBindJSON(&drivingHistory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify user belongs to agent's agency
		var targetUserAgencyID int
		err = db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", drivingHistory.UserID).Scan(&targetUserAgencyID)
		if err != nil || targetUserAgencyID != agencyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to user"})
			return
		}

		_, err = db.Exec(`INSERT INTO DrivingHistory (UserID, IncidentType, IncidentDate, ConvictionDate, 
			Amount, Description, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			drivingHistory.UserID, drivingHistory.IncidentType, drivingHistory.IncidentDate,
			drivingHistory.ConvictionDate, drivingHistory.Amount, drivingHistory.Description,
			drivingHistory.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Driving history added successfully"})
	}
}

// AddInsuranceHistory adds insurance history for a user
func AddInsuranceHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var insuranceHistory models.InsuranceHistory
		if err := c.ShouldBindJSON(&insuranceHistory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify user belongs to agent's agency
		var targetUserAgencyID int
		err = db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", insuranceHistory.UserID).Scan(&targetUserAgencyID)
		if err != nil || targetUserAgencyID != agencyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to user"})
			return
		}

		_, err = db.Exec(`INSERT INTO InsuranceHistory (UserID, InsuranceStatus, CurrentCarrier, 
			CurrentBodilyInjuryLimits, LengthWithCurrentCompany, ContinuousInsurance, 
			VehicleRegisteredToOther, DriverWithoutLicense, LicenseSuspendedRevoked, 
			DeclinedCancelledNonRenewed, MilitaryDeployment, PrimaryResidence, AdditionalData) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
			insuranceHistory.UserID, insuranceHistory.InsuranceStatus, insuranceHistory.CurrentCarrier,
			insuranceHistory.CurrentBodilyInjuryLimits, insuranceHistory.LengthWithCurrentCompany,
			insuranceHistory.ContinuousInsurance, insuranceHistory.VehicleRegisteredToOther,
			insuranceHistory.DriverWithoutLicense, insuranceHistory.LicenseSuspendedRevoked,
			insuranceHistory.DeclinedCancelledNonRenewed, insuranceHistory.MilitaryDeployment,
			insuranceHistory.PrimaryResidence, insuranceHistory.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Insurance history added successfully"})
	}
}

// GetAgentQuotes lists all quotes for the current agent
func GetAgentQuotes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		rows, err := db.Query(`SELECT QuoteID, AgentID, CustomerID, VehicleID, QuoteDate, Status, AdditionalData 
			FROM Quotes WHERE AgentID = $1`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		quotes := []models.Quote{}
		for rows.Next() {
			var quote models.Quote
			if err := rows.Scan(&quote.QuoteID, &quote.AgentID, &quote.CustomerID, &quote.VehicleID,
				&quote.QuoteDate, &quote.Status, &quote.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			quotes = append(quotes, quote)
		}

		c.JSON(http.StatusOK, quotes)
	}
}

// CreateQuote creates a new quote
func CreateQuote(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var quote models.Quote
		if err := c.ShouldBindJSON(&quote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set agent ID to current user
		quote.AgentID = userID.(int)

		// Verify customer belongs to agent's agency
		var customerAgencyID int
		err = db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1 AND Role = 'Customer'", quote.CustomerID).Scan(&customerAgencyID)
		if err != nil || customerAgencyID != agencyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to customer"})
			return
		}

		// Verify vehicle belongs to the customer (if provided)
		if quote.VehicleID != nil {
			var vehicleCustomerID int
			err = db.QueryRow("SELECT CustomerID FROM Vehicles WHERE VehicleID = $1", *quote.VehicleID).Scan(&vehicleCustomerID)
			if err != nil || vehicleCustomerID != quote.CustomerID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle does not belong to the specified customer"})
				return
			}
		}

		// Set default status to Draft
		quote.Status = models.QuoteStatusDraft
		quote.QuoteDate = time.Now()

		_, err = db.Exec(`INSERT INTO Quotes (AgentID, CustomerID, VehicleID, QuoteDate, Status, AdditionalData) 
			VALUES ($1, $2, $3, $4, $5, $6)`,
			quote.AgentID, quote.CustomerID, quote.VehicleID, quote.QuoteDate, quote.Status, quote.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Quote created successfully"})
	}
}

// UpdateQuote updates a quote (mainly status changes)
func UpdateQuote(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		quoteID := c.Param("id")

		var quote models.Quote
		if err := c.ShouldBindJSON(&quote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only quotes belonging to the current agent
		result, err := db.Exec("UPDATE Quotes SET Status = $1, AdditionalData = $2 WHERE QuoteID = $3 AND AgentID = $4",
			quote.Status, quote.AdditionalData, quoteID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Quote updated successfully"})
	}
}

// AddQuoteLineItem adds a coverage line item to a quote
func AddQuoteLineItem(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var lineItem models.QuoteLineItem
		if err := c.ShouldBindJSON(&lineItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify quote belongs to current agent
		var quoteAgentID int
		err = db.QueryRow("SELECT AgentID FROM Quotes WHERE QuoteID = $1", lineItem.QuoteID).Scan(&quoteAgentID)
		if err != nil || quoteAgentID != userID.(int) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to quote"})
			return
		}

		// Verify provider is accessible to agent's agency
		var providerAccess int
		err = db.QueryRow("SELECT COUNT(*) FROM AgencyProviderAccess WHERE AgencyID = $1 AND ProviderID = $2",
			agencyID, lineItem.ProviderID).Scan(&providerAccess)
		if err != nil || providerAccess == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Provider not accessible to your agency"})
			return
		}

		_, err = db.Exec(`INSERT INTO QuoteLineItems (QuoteID, ProviderID, CoverageID, Price, LimitAmount, 
			DeductibleAmount, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			lineItem.QuoteID, lineItem.ProviderID, lineItem.CoverageID, lineItem.Price,
			lineItem.LimitAmount, lineItem.DeductibleAmount, lineItem.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Quote line item added successfully"})
	}
}

// GetCoverages lists all available coverages
func GetCoverages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT CoverageID, CoverageName, Description, AdditionalData FROM Coverages")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		coverages := []models.Coverage{}
		for rows.Next() {
			var coverage models.Coverage
			if err := rows.Scan(&coverage.CoverageID, &coverage.CoverageName, &coverage.Description,
				&coverage.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			coverages = append(coverages, coverage)
		}

		c.JSON(http.StatusOK, coverages)
	}
}

// GetAgentProviders lists all accessible providers for the current agent's agency
func GetAgentProviders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT ip.ProviderID, ip.ProviderName, ip.ContactInfo, ip.AdditionalData 
			FROM InsuranceProviders ip 
			JOIN AgencyProviderAccess apa ON ip.ProviderID = apa.ProviderID 
			WHERE apa.AgencyID = $1`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		providers := []models.InsuranceProvider{}
		for rows.Next() {
			var provider models.InsuranceProvider
			if err := rows.Scan(&provider.ProviderID, &provider.ProviderName, &provider.ContactInfo,
				&provider.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			providers = append(providers, provider)
		}

		c.JSON(http.StatusOK, providers)
	}
}

// BindPolicy binds a quote to create a policy
func BindPolicy(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var policy models.Policy
		if err := c.ShouldBindJSON(&policy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify quote belongs to current agent
		var quoteAgentID, quoteCustomerID int
		err := db.QueryRow("SELECT AgentID, CustomerID FROM Quotes WHERE QuoteID = $1", policy.QuoteID).Scan(&quoteAgentID, &quoteCustomerID)
		if err != nil || quoteAgentID != userID.(int) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to quote"})
			return
		}

		// Begin transaction to update quote status and create policy atomically
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}
		defer tx.Rollback()

		// Update quote status to Bound
		_, err = tx.Exec("UPDATE Quotes SET Status = 'Bound' WHERE QuoteID = $1", policy.QuoteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quote status"})
			return
		}

		// Generate policy number if not provided
		if policy.PolicyNumber == "" {
			policy.PolicyNumber = fmt.Sprintf("POL-%d-%d", time.Now().Unix(), *policy.QuoteID)
		}

		// Set customer ID from quote
		policy.CustomerID = quoteCustomerID

		// Set default status to Active
		policy.Status = models.PolicyStatusActive

		// Insert policy
		_, err = tx.Exec(`INSERT INTO Policies (CustomerID, ProviderID, QuoteID, PolicyNumber, 
			EffectiveDate, ExpirationDate, Status, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			policy.CustomerID, policy.ProviderID, policy.QuoteID, policy.PolicyNumber,
			policy.EffectiveDate, policy.ExpirationDate, policy.Status, policy.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Policy bound successfully", "policy_number": policy.PolicyNumber})
	}
}

// GetPolicies lists all policies for customers in the current agent's agency
func GetPolicies(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get agent's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT p.PolicyID, p.CustomerID, p.ProviderID, p.QuoteID, p.PolicyNumber, 
			p.EffectiveDate, p.ExpirationDate, p.Status, p.AdditionalData 
			FROM Policies p 
			JOIN Users u ON p.CustomerID = u.UserID 
			WHERE u.AgencyID = $1`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		policies := []models.Policy{}
		for rows.Next() {
			var policy models.Policy
			if err := rows.Scan(&policy.PolicyID, &policy.CustomerID, &policy.ProviderID, &policy.QuoteID,
				&policy.PolicyNumber, &policy.EffectiveDate, &policy.ExpirationDate, &policy.Status,
				&policy.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			policies = append(policies, policy)
		}

		c.JSON(http.StatusOK, policies)
	}
}

// GetAgentReport generates agent-specific reports
func GetAgentReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get quotes stats for the agent
		var quotesReport struct {
			TotalQuotes     int `json:"total_quotes"`
			DraftQuotes     int `json:"draft_quotes"`
			PresentedQuotes int `json:"presented_quotes"`
			BoundQuotes     int `json:"bound_quotes"`
		}

		err := db.QueryRow(`SELECT COUNT(*) as total_quotes, 
			COUNT(CASE WHEN Status = 'Draft' THEN 1 END) as draft_quotes,
			COUNT(CASE WHEN Status = 'Presented' THEN 1 END) as presented_quotes,
			COUNT(CASE WHEN Status = 'Bound' THEN 1 END) as bound_quotes
			FROM Quotes WHERE AgentID = $1`, userID).Scan(
			&quotesReport.TotalQuotes, &quotesReport.DraftQuotes,
			&quotesReport.PresentedQuotes, &quotesReport.BoundQuotes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get customer count
		var customerCount int
		err = db.QueryRow(`SELECT COUNT(*) FROM Users u1 
			WHERE u1.Role = 'Customer' AND u1.AgencyID = (SELECT AgencyID FROM Users WHERE UserID = $1)`, userID).Scan(&customerCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		report := gin.H{
			"agent_id":       userID,
			"customer_count": customerCount,
			"quotes":         quotesReport,
		}

		c.JSON(http.StatusOK, report)
	}
}
