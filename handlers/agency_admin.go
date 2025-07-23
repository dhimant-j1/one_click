package handlers

import (
	"database/sql"
	"net/http"
	"one_click_2/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetLocations lists all locations for the current agency
func GetLocations(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query("SELECT LocationID, AgencyID, Address, AdditionalData FROM AgencyLocations WHERE AgencyID = $1", agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		locations := []models.AgencyLocation{}
		for rows.Next() {
			var location models.AgencyLocation
			if err := rows.Scan(&location.LocationID, &location.AgencyID, &location.Address, &location.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			locations = append(locations, location)
		}

		c.JSON(http.StatusOK, locations)
	}
}

// CreateLocation creates a new location for the current agency
func CreateLocation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var location models.AgencyLocation
		if err := c.ShouldBindJSON(&location); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set the agency ID to the current user's agency
		location.AgencyID = agencyID

		_, err = db.Exec("INSERT INTO AgencyLocations (AgencyID, Address, AdditionalData) VALUES ($1, $2, $3)",
			location.AgencyID, location.Address, location.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Location created successfully"})
	}
}

// UpdateLocation updates a location within the current agency
func UpdateLocation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		locationID := c.Param("id")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var location models.AgencyLocation
		if err := c.ShouldBindJSON(&location); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only if location belongs to user's agency
		result, err := db.Exec("UPDATE AgencyLocations SET Address = $1, AdditionalData = $2 WHERE LocationID = $3 AND AgencyID = $4",
			location.Address, location.AdditionalData, locationID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Location updated successfully"})
	}
}

// DeleteLocation deletes a location within the current agency
func DeleteLocation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		locationID := c.Param("id")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		// Delete only if location belongs to user's agency
		result, err := db.Exec("DELETE FROM AgencyLocations WHERE LocationID = $1 AND AgencyID = $2", locationID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Location deleted successfully"})
	}
}

// GetAgents lists all agents within the current agency
func GetAgents(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT UserID, FirstName, LastName, Email, Role, AgencyID, LocationID, 
			Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, 
			NeedsFinancialFiling, AdditionalData FROM Users WHERE Role = 'Agent' AND AgencyID = $1`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		agents := []models.User{}
		for rows.Next() {
			var agent models.User
			if err := rows.Scan(&agent.UserID, &agent.FirstName, &agent.LastName, &agent.Email, &agent.Role,
				&agent.AgencyID, &agent.LocationID, &agent.Address, &agent.LicenseNumber, &agent.LicenseState,
				&agent.DateOfBirth, &agent.MobileNumber, &agent.HasPhysicalImpairment, &agent.NeedsFinancialFiling,
				&agent.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			agents = append(agents, agent)
		}

		c.JSON(http.StatusOK, agents)
	}
}

// CreateAgent creates a new agent within the current agency
func CreateAgent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var agent models.User
		if err := c.ShouldBindJSON(&agent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(agent.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Set role to Agent and agency to current user's agency
		agent.Role = models.RoleAgent
		agent.AgencyID = &agencyID

		_, err = db.Exec(`INSERT INTO Users (FirstName, LastName, Email, PasswordHash, Role, AgencyID, LocationID, 
			Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, 
			NeedsFinancialFiling, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
			agent.FirstName, agent.LastName, agent.Email, string(hashedPassword), agent.Role, agent.AgencyID,
			agent.LocationID, agent.Address, agent.LicenseNumber, agent.LicenseState, agent.DateOfBirth,
			agent.MobileNumber, agent.HasPhysicalImpairment, agent.NeedsFinancialFiling, agent.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Agent created successfully"})
	}
}

// UpdateAgent updates an agent within the current agency
func UpdateAgent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		agentID := c.Param("id")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var agent models.User
		if err := c.ShouldBindJSON(&agent); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only if agent belongs to user's agency
		result, err := db.Exec(`UPDATE Users SET FirstName = $1, LastName = $2, Email = $3, LocationID = $4, 
			Address = $5, LicenseNumber = $6, LicenseState = $7, DateOfBirth = $8, MobileNumber = $9, 
			HasPhysicalImpairment = $10, NeedsFinancialFiling = $11, AdditionalData = $12 
			WHERE UserID = $13 AND AgencyID = $14 AND Role = 'Agent'`,
			agent.FirstName, agent.LastName, agent.Email, agent.LocationID, agent.Address, agent.LicenseNumber,
			agent.LicenseState, agent.DateOfBirth, agent.MobileNumber, agent.HasPhysicalImpairment,
			agent.NeedsFinancialFiling, agent.AdditionalData, agentID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Agent updated successfully"})
	}
}

// DeleteAgent deletes an agent within the current agency
func DeleteAgent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		agentID := c.Param("id")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		// Delete only if agent belongs to user's agency
		result, err := db.Exec("DELETE FROM Users WHERE UserID = $1 AND AgencyID = $2 AND Role = 'Agent'", agentID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
	}
}

// GetAgencyUsers lists all users (LocationAdmins and Agents) within the current agency
func GetAgencyUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT UserID, FirstName, LastName, Email, Role, AgencyID, LocationID, 
			Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, 
			NeedsFinancialFiling, AdditionalData FROM Users WHERE AgencyID = $1 AND Role IN ('LocationAdmin', 'Agent')`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		users := []models.User{}
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Role,
				&user.AgencyID, &user.LocationID, &user.Address, &user.LicenseNumber, &user.LicenseState,
				&user.DateOfBirth, &user.MobileNumber, &user.HasPhysicalImpairment, &user.NeedsFinancialFiling,
				&user.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

// CreateAgencyUser creates a new LocationAdmin or Agent within the current agency
func CreateAgencyUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate role
		if user.Role != models.RoleLocationAdmin && user.Role != models.RoleAgent {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Only LocationAdmin and Agent roles are allowed"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Set agency to current user's agency
		user.AgencyID = &agencyID

		_, err = db.Exec(`INSERT INTO Users (FirstName, LastName, Email, PasswordHash, Role, AgencyID, LocationID, 
			Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, 
			NeedsFinancialFiling, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
			user.FirstName, user.LastName, user.Email, string(hashedPassword), user.Role, user.AgencyID,
			user.LocationID, user.Address, user.LicenseNumber, user.LicenseState, user.DateOfBirth,
			user.MobileNumber, user.HasPhysicalImpairment, user.NeedsFinancialFiling, user.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

// UpdateAgencyUser updates a user within the current agency
func UpdateAgencyUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		targetUserID := c.Param("id")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update only if user belongs to the same agency
		result, err := db.Exec(`UPDATE Users SET FirstName = $1, LastName = $2, Email = $3, LocationID = $4, 
			Address = $5, LicenseNumber = $6, LicenseState = $7, DateOfBirth = $8, MobileNumber = $9, 
			HasPhysicalImpairment = $10, NeedsFinancialFiling = $11, AdditionalData = $12 
			WHERE UserID = $13 AND AgencyID = $14`,
			user.FirstName, user.LastName, user.Email, user.LocationID, user.Address, user.LicenseNumber,
			user.LicenseState, user.DateOfBirth, user.MobileNumber, user.HasPhysicalImpairment,
			user.NeedsFinancialFiling, user.AdditionalData, targetUserID, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found or access denied"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

// GetAgencyProviders lists all accessible providers for the current agency
func GetAgencyProviders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
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
			if err := rows.Scan(&provider.ProviderID, &provider.ProviderName, &provider.ContactInfo, &provider.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			providers = append(providers, provider)
		}

		c.JSON(http.StatusOK, providers)
	}
}

// GetAgencyQuotesReport generates quotes report for the current agency
func GetAgencyQuotesReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT COUNT(*) as total_quotes, 
			COUNT(CASE WHEN Status = 'Draft' THEN 1 END) as draft_quotes,
			COUNT(CASE WHEN Status = 'Presented' THEN 1 END) as presented_quotes,
			COUNT(CASE WHEN Status = 'Bound' THEN 1 END) as bound_quotes
			FROM Quotes q 
			JOIN Users u ON q.AgentID = u.UserID 
			WHERE u.AgencyID = $1`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var report struct {
			TotalQuotes     int `json:"total_quotes"`
			DraftQuotes     int `json:"draft_quotes"`
			PresentedQuotes int `json:"presented_quotes"`
			BoundQuotes     int `json:"bound_quotes"`
		}

		if rows.Next() {
			if err := rows.Scan(&report.TotalQuotes, &report.DraftQuotes, &report.PresentedQuotes, &report.BoundQuotes); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, report)
	}
}

// GetAgencyPoliciesReport generates policies report for the current agency
func GetAgencyPoliciesReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		// Get user's agency ID
		var agencyID int
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", userID).Scan(&agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user agency"})
			return
		}

		rows, err := db.Query(`SELECT COUNT(*) as total_policies, 
			COUNT(CASE WHEN Status = 'Active' THEN 1 END) as active_policies,
			COUNT(CASE WHEN Status = 'Expired' THEN 1 END) as expired_policies,
			COUNT(CASE WHEN Status = 'Cancelled' THEN 1 END) as cancelled_policies
			FROM Policies p 
			JOIN Users u ON p.CustomerID = u.UserID 
			WHERE u.AgencyID = $1`, agencyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var report struct {
			TotalPolicies     int `json:"total_policies"`
			ActivePolicies    int `json:"active_policies"`
			ExpiredPolicies   int `json:"expired_policies"`
			CancelledPolicies int `json:"cancelled_policies"`
		}

		if rows.Next() {
			if err := rows.Scan(&report.TotalPolicies, &report.ActivePolicies, &report.ExpiredPolicies, &report.CancelledPolicies); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, report)
	}
}
