package handlers

import (
	"database/sql"
	"net/http"
	"one_click_2/auth"
	"one_click_2/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetAgencies retrieves all agencies
func GetAgencies(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM Agencies")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		agencies := []models.Agency{}
		for rows.Next() {
			var agency models.Agency
			if err := rows.Scan(&agency.AgencyID, &agency.AgencyName, &agency.AgentCode, &agency.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			agencies = append(agencies, agency)
		}

		c.JSON(http.StatusOK, agencies)
	}
}

// CreateAgency creates a new agency
func CreateAgency(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var agency models.Agency
		if err := c.ShouldBindJSON(&agency); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO Agencies (AgencyName, AgentCode, AdditionalData) VALUES ($1, $2, $3)", agency.AgencyName, agency.AgentCode, agency.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, agency)
	}
}

// UpdateAgency updates an existing agency
func UpdateAgency(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var agency models.Agency
		if err := c.ShouldBindJSON(&agency); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("UPDATE Agencies SET AgencyName = $1, AgentCode = $2, AdditionalData = $3 WHERE AgencyID = $4", agency.AgencyName, agency.AgentCode, agency.AdditionalData, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, agency)
	}
}

// DeleteAgency deletes an agency
func DeleteAgency(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM Agencies WHERE AgencyID = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Agency deleted successfully"})
	}
}

// GetProviders retrieves all insurance providers
func GetProviders(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM InsuranceProviders")
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

// CreateProvider creates a new insurance provider
func CreateProvider(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var provider models.InsuranceProvider
		if err := c.ShouldBindJSON(&provider); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO InsuranceProviders (ProviderName, ContactInfo, AdditionalData) VALUES ($1, $2, $3)", provider.ProviderName, provider.ContactInfo, provider.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, provider)
	}
}

// UpdateProvider updates an existing insurance provider
func UpdateProvider(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var provider models.InsuranceProvider
		if err := c.ShouldBindJSON(&provider); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("UPDATE InsuranceProviders SET ProviderName = $1, ContactInfo = $2, AdditionalData = $3 WHERE ProviderID = $4", provider.ProviderName, provider.ContactInfo, provider.AdditionalData, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, provider)
	}
}

// DeleteProvider deletes an insurance provider
func DeleteProvider(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM InsuranceProviders WHERE ProviderID = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Provider deleted successfully"})
	}
}

// GrantProviderAccess grants an agency access to an insurance provider
func GrantProviderAccess(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var access models.AgencyProviderAccess
		if err := c.ShouldBindJSON(&access); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO AgencyProviderAccess (AgencyID, ProviderID) VALUES ($1, $2)", access.AgencyID, access.ProviderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, access)
	}
}

// RevokeProviderAccess revokes an agency's access to an insurance provider
func RevokeProviderAccess(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var access models.AgencyProviderAccess
		if err := c.ShouldBindJSON(&access); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("DELETE FROM AgencyProviderAccess WHERE AgencyID = $1 AND ProviderID = $2", access.AgencyID, access.ProviderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Provider access revoked successfully"})
	}
}

// GetUsers retrieves all users
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT UserID, FirstName, LastName, Email, Role, AgencyID, LocationID, Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, NeedsFinancialFiling, AdditionalData FROM Users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		users := []models.User{}
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.AgencyID, &user.LocationID, &user.Address, &user.LicenseNumber, &user.LicenseState, &user.DateOfBirth, &user.MobileNumber, &user.HasPhysicalImpairment, &user.NeedsFinancialFiling, &user.AdditionalData); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

// CreateUser creates a new user
func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.PasswordHash = string(hashedPassword)

		_, err = db.Exec("INSERT INTO Users (FirstName, LastName, Email, PasswordHash, Role, AgencyID, LocationID, Address, LicenseNumber, LicenseState, DateOfBirth, MobileNumber, HasPhysicalImpairment, NeedsFinancialFiling, AdditionalData) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Role, user.AgencyID, user.LocationID, user.Address, user.LicenseNumber, user.LicenseState, user.DateOfBirth, user.MobileNumber, user.HasPhysicalImpairment, user.NeedsFinancialFiling, user.AdditionalData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// UpdateUser updates an existing user
func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("UPDATE Users SET FirstName = $1, LastName = $2, Email = $3, Role = $4, AgencyID = $5, LocationID = $6, Address = $7, LicenseNumber = $8, LicenseState = $9, DateOfBirth = $10, MobileNumber = $11, HasPhysicalImpairment = $12, NeedsFinancialFiling = $13, AdditionalData = $14 WHERE UserID = $15", user.FirstName, user.LastName, user.Email, user.Role, user.AgencyID, user.LocationID, user.Address, user.LicenseNumber, user.LicenseState, user.DateOfBirth, user.MobileNumber, user.HasPhysicalImpairment, user.NeedsFinancialFiling, user.AdditionalData, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// DeleteUser deletes a user
func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM Users WHERE UserID = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

// GetQuotesReport generates a report of quotes
func GetQuotesReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`SELECT COUNT(*) as total_quotes, 
			COUNT(CASE WHEN Status = 'Draft' THEN 1 END) as draft_quotes,
			COUNT(CASE WHEN Status = 'Presented' THEN 1 END) as presented_quotes,
			COUNT(CASE WHEN Status = 'Bound' THEN 1 END) as bound_quotes
			FROM Quotes`)
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

// GetPoliciesReport generates a report of policies
func GetPoliciesReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`SELECT COUNT(*) as total_policies, 
			COUNT(CASE WHEN Status = 'Active' THEN 1 END) as active_policies,
			COUNT(CASE WHEN Status = 'Expired' THEN 1 END) as expired_policies,
			COUNT(CASE WHEN Status = 'Cancelled' THEN 1 END) as cancelled_policies
			FROM Policies`)
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

// Login handles user login
func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user models.User
		err := db.QueryRow("SELECT UserID, PasswordHash, Role FROM Users WHERE Email = $1", creds.Email).Scan(&user.UserID, &user.PasswordHash, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := auth.GenerateJWT(user.UserID, string(user.Role))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		_, err = db.Exec("INSERT INTO UserSessions (UserID, SessionToken, ExpiresAt) VALUES ($1, $2, $3)", user.UserID, token, time.Now().Add(24*time.Hour))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}
		user.PasswordHash = ""
		c.JSON(http.StatusOK, gin.H{"token": token, "profile": user})
	}
}

// Logout handles user logout
func Logout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := db.Exec("DELETE FROM UserSessions WHERE SessionToken = $1", tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}

// ResetPassword handles password reset requests
func ResetPassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Check if the email exists in the database
		var userID int
		err := db.QueryRow("SELECT UserID FROM Users WHERE Email = $1", req.Email).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Generate a reset token
		resetToken, err := auth.GenerateResetToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
			return
		}

		// Store the reset token in the database
		_, err = db.Exec("INSERT INTO PasswordResetTokens (UserID, ResetToken, ExpiresAt) VALUES ($1, $2, $3)", userID, resetToken, time.Now().Add(1*time.Hour))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store reset token"})
			return
		}

		// TODO: Send the reset token to the user's email address
		// You can use a library like "net/smtp" or a third-party email service

		c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to your email address"})
	}
}

// UpdatePassword handles password updates
func UpdatePassword(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		var req struct {
			NewPassword string `json:"newPassword"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Verify the reset token
		var userID int
		err := db.QueryRow("SELECT UserID FROM PasswordResetTokens WHERE ResetToken = $1 AND ExpiresAt > $2 AND IsUsed = false", token, time.Now()).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired reset token"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Update the user's password in the database
		_, err = db.Exec("UPDATE Users SET PasswordHash = $1 WHERE UserID = $2", hashedPassword, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Mark the reset token as used instead of deleting it
		_, err = db.Exec("UPDATE PasswordResetTokens SET IsUsed = true WHERE ResetToken = $1", token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reset token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}
