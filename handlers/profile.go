package handlers

import (
	"database/sql"
	"net/http"
	"one_click_2/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProfileResponse represents the user profile response without sensitive data
type ProfileResponse struct {
	UserID                int             `json:"user_id"`
	FirstName             string          `json:"first_name"`
	LastName              string          `json:"last_name"`
	Email                 string          `json:"email"`
	Role                  models.UserRole `json:"role"`
	AgencyID              *int            `json:"agency_id,omitempty"`
	LocationID            *int            `json:"location_id,omitempty"`
	Address               *string         `json:"address,omitempty"`
	LicenseNumber         *string         `json:"license_number,omitempty"`
	LicenseState          *string         `json:"license_state,omitempty"`
	DateOfBirth           *string         `json:"date_of_birth,omitempty"`
	MobileNumber          *string         `json:"mobile_number,omitempty"`
	HasPhysicalImpairment *bool           `json:"has_physical_impairment,omitempty"`
	NeedsFinancialFiling  *bool           `json:"needs_financial_filing,omitempty"`
	AdditionalData        *string         `json:"additional_data,omitempty"`
}

// GetMyProfile returns the profile of the currently logged-in user
func GetMyProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		profile, err := getUserProfile(db, userID.(int))
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, profile)
	}
}

// GetUserProfile returns the profile of a specific user (role-based access control)
func GetUserProfile(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetUserID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		currentUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		currentRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		// Check if user has permission to view this profile
		hasPermission, err := checkProfileViewPermission(db, currentUserID.(int), currentRole.(string), targetUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view this profile"})
			return
		}

		profile, err := getUserProfile(db, targetUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, profile)
	}
}

// GetUsersByRole returns users filtered by role (role-based access control)
func GetUsersByRole(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetRole := c.Query("role")
		if targetRole == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role parameter is required"})
			return
		}

		currentUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		currentRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			return
		}

		// Check if user has permission to view users with this role
		hasPermission := checkRoleViewPermission(currentRole.(string), targetRole)
		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to view users with this role"})
			return
		}

		// Build query based on current user's role and permissions
		query, args, err := buildUsersByRoleQuery(db, currentUserID.(int), currentRole.(string), targetRole)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		profiles := []ProfileResponse{}
		for rows.Next() {
			var profile ProfileResponse
			var dateOfBirth sql.NullString

			err := rows.Scan(
				&profile.UserID, &profile.FirstName, &profile.LastName, &profile.Email,
				&profile.Role, &profile.AgencyID, &profile.LocationID, &profile.Address,
				&profile.LicenseNumber, &profile.LicenseState, &dateOfBirth,
				&profile.MobileNumber, &profile.HasPhysicalImpairment,
				&profile.NeedsFinancialFiling, &profile.AdditionalData,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if dateOfBirth.Valid {
				profile.DateOfBirth = &dateOfBirth.String
			}

			profiles = append(profiles, profile)
		}

		c.JSON(http.StatusOK, profiles)
	}
}

// getUserProfile retrieves a user's profile by ID
func getUserProfile(db *sql.DB, userID int) (*ProfileResponse, error) {
	var profile ProfileResponse
	var dateOfBirth sql.NullString

	query := `
		SELECT UserID, FirstName, LastName, Email, Role, AgencyID, LocationID, 
		       Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, 
		       HasPhysicalImpairment, NeedsFinancialFiling, AdditionalData 
		FROM Users 
		WHERE UserID = $1
	`

	err := db.QueryRow(query, userID).Scan(
		&profile.UserID, &profile.FirstName, &profile.LastName, &profile.Email,
		&profile.Role, &profile.AgencyID, &profile.LocationID, &profile.Address,
		&profile.LicenseNumber, &profile.LicenseState, &dateOfBirth,
		&profile.MobileNumber, &profile.HasPhysicalImpairment,
		&profile.NeedsFinancialFiling, &profile.AdditionalData,
	)

	if err != nil {
		return nil, err
	}

	if dateOfBirth.Valid {
		profile.DateOfBirth = &dateOfBirth.String
	}

	return &profile, nil
}

// checkProfileViewPermission checks if current user can view target user's profile
func checkProfileViewPermission(db *sql.DB, currentUserID int, currentRole string, targetUserID int) (bool, error) {
	// Users can always view their own profile
	if currentUserID == targetUserID {
		return true, nil
	}

	// Get target user's details
	var targetRole string
	var targetAgencyID sql.NullInt64
	err := db.QueryRow("SELECT Role, AgencyID FROM Users WHERE UserID = $1", targetUserID).Scan(&targetRole, &targetAgencyID)
	if err != nil {
		return false, err
	}

	switch currentRole {
	case "MasterAdmin":
		// Master admin can see profile of all users
		return true, nil

	case "AgencyAdmin", "LocationAdmin":
		// Agency admin can see profile of agent, agency admin, and users in their agency
		if targetRole == "Agent" || targetRole == "AgencyAdmin" || targetRole == "LocationAdmin" || targetRole == "Customer" {
			// Check if they belong to the same agency
			var currentAgencyID sql.NullInt64
			err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", currentUserID).Scan(&currentAgencyID)
			if err != nil {
				return false, err
			}

			if currentAgencyID.Valid && targetAgencyID.Valid && currentAgencyID.Int64 == targetAgencyID.Int64 {
				return true, nil
			}
		}
		return false, nil

	case "Agent":
		// Agent can see profile of users (customers) and other agents in their agency
		if targetRole == "Customer" || targetRole == "Agent" {
			// Check if they belong to the same agency
			var currentAgencyID sql.NullInt64
			err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", currentUserID).Scan(&currentAgencyID)
			if err != nil {
				return false, err
			}

			if currentAgencyID.Valid && targetAgencyID.Valid && currentAgencyID.Int64 == targetAgencyID.Int64 {
				return true, nil
			}
		}
		return false, nil

	default:
		return false, nil
	}
}

// checkRoleViewPermission checks if current user can view users with target role
func checkRoleViewPermission(currentRole, targetRole string) bool {
	switch currentRole {
	case "MasterAdmin":
		// Master admin can view all roles
		return true

	case "AgencyAdmin", "LocationAdmin":
		// Agency admin can view agents, agency admins, and customers
		return targetRole == "Agent" || targetRole == "AgencyAdmin" || targetRole == "LocationAdmin" || targetRole == "Customer"

	case "Agent":
		// Agent can view customers and other agents
		return targetRole == "Customer" || targetRole == "Agent"

	default:
		return false
	}
}

// buildUsersByRoleQuery builds the appropriate query based on user role and permissions
func buildUsersByRoleQuery(db *sql.DB, currentUserID int, currentRole, targetRole string) (string, []interface{}, error) {
	baseQuery := `
		SELECT UserID, FirstName, LastName, Email, Role, AgencyID, LocationID, 
		       Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, 
		       HasPhysicalImpairment, NeedsFinancialFiling, AdditionalData 
		FROM Users 
		WHERE Role = $1
	`

	switch currentRole {
	case "MasterAdmin":
		// Master admin can see all users of target role
		return baseQuery, []interface{}{targetRole}, nil

	case "AgencyAdmin", "LocationAdmin":
		// Agency admin can only see users in their agency
		var currentAgencyID sql.NullInt64
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", currentUserID).Scan(&currentAgencyID)
		if err != nil {
			return "", nil, err
		}

		if !currentAgencyID.Valid {
			return "", nil, nil
		}

		query := baseQuery + " AND AgencyID = $2"
		return query, []interface{}{targetRole, currentAgencyID.Int64}, nil

	case "Agent":
		// Agent can only see users in their agency
		var currentAgencyID sql.NullInt64
		err := db.QueryRow("SELECT AgencyID FROM Users WHERE UserID = $1", currentUserID).Scan(&currentAgencyID)
		if err != nil {
			return "", nil, err
		}

		if !currentAgencyID.Valid {
			return "", nil, nil
		}

		query := baseQuery + " AND AgencyID = $2"
		return query, []interface{}{targetRole, currentAgencyID.Int64}, nil

	default:
		return "", nil, nil
	}
}
