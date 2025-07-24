package main

import (
	"fmt"
	"one_click_2/database"
	"one_click_2/handlers"
	"one_click_2/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	// connStr := "user=postgres password=postgres dbname=insurance sslmode=disable"
	db := database.Connect()

	defer db.Close()

	fmt.Println("Successfully connected to database!")

	// Initialize Gin router
	router := gin.Default()

	// Master Admin routes
	masterAdmin := router.Group("/api/master-admin")
	masterAdmin.Use(middleware.AuthMiddleware())
	masterAdmin.Use(middleware.MasterAdminMiddleware())
	{
		masterAdmin.GET("/agencies", handlers.GetAgencies(db))
		masterAdmin.POST("/agencies", handlers.CreateAgency(db))
		masterAdmin.PUT("/agencies/:id", handlers.UpdateAgency(db))
		masterAdmin.DELETE("/agencies/:id", handlers.DeleteAgency(db))
		masterAdmin.GET("/providers", handlers.GetProviders(db))
		masterAdmin.POST("/providers", handlers.CreateProvider(db))
		masterAdmin.PUT("/providers/:id", handlers.UpdateProvider(db))
		masterAdmin.DELETE("/providers/:id", handlers.DeleteProvider(db))
		masterAdmin.POST("/provider-access", handlers.GrantProviderAccess(db))
		masterAdmin.DELETE("/provider-access", handlers.RevokeProviderAccess(db))
		masterAdmin.GET("/users", handlers.GetUsers(db))
		masterAdmin.POST("/users", handlers.CreateUser(db))
		masterAdmin.PUT("/users/:id", handlers.UpdateUser(db))
		masterAdmin.DELETE("/users/:id", handlers.DeleteUser(db))
		masterAdmin.GET("/reports/quotes", handlers.GetQuotesReport(db))
		masterAdmin.GET("/reports/policies", handlers.GetPoliciesReport(db))
		masterAdmin.POST("/auth/logout", handlers.Logout(db))
		// Profile management for master admin
		masterAdmin.GET("/profiles/users", handlers.GetUsersByRole(db))
		masterAdmin.GET("/profiles/user/:id", handlers.GetUserProfile(db))
	}

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", handlers.Login(db))
		auth.POST("/reset-password", handlers.ResetPassword(db))
		auth.PUT("/reset-password/:token", handlers.UpdatePassword(db))
	}

	// Profile routes (requires authentication)
	profile := router.Group("/api/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		// Get current user's profile
		profile.GET("/me", handlers.GetMyProfile(db))
		// Get specific user's profile (with role-based access control)
		profile.GET("/user/:id", handlers.GetUserProfile(db))
		// Get users by role (with role-based access control)
		profile.GET("/users", handlers.GetUsersByRole(db))
	}

	// Agency Admin routes
	agencyAdmin := router.Group("/api/agency-admin")
	agencyAdmin.Use(middleware.AuthMiddleware())
	agencyAdmin.Use(middleware.AgencyAdminMiddleware())
	{
		agencyAdmin.GET("/locations", handlers.GetLocations(db))
		agencyAdmin.POST("/locations", handlers.CreateLocation(db))
		agencyAdmin.PUT("/locations/:id", handlers.UpdateLocation(db))
		agencyAdmin.DELETE("/locations/:id", handlers.DeleteLocation(db))
		agencyAdmin.GET("/agents", handlers.GetAgents(db))
		agencyAdmin.POST("/agents", handlers.CreateAgent(db))
		agencyAdmin.PUT("/agents/:id", handlers.UpdateAgent(db))
		agencyAdmin.DELETE("/agents/:id", handlers.DeleteAgent(db))
		agencyAdmin.GET("/users", handlers.GetAgencyUsers(db))
		agencyAdmin.POST("/users", handlers.CreateAgencyUser(db))
		agencyAdmin.PUT("/users/:id", handlers.UpdateAgencyUser(db))
		agencyAdmin.GET("/providers", handlers.GetAgencyProviders(db))
		agencyAdmin.GET("/reports/quotes", handlers.GetAgencyQuotesReport(db))
		agencyAdmin.GET("/reports/policies", handlers.GetAgencyPoliciesReport(db))
		agencyAdmin.POST("/auth/logout", handlers.Logout(db))
		// Profile management for agency admin
		agencyAdmin.GET("/profiles/users", handlers.GetUsersByRole(db))
		agencyAdmin.GET("/profiles/user/:id", handlers.GetUserProfile(db))
	}

	// Agent routes
	agent := router.Group("/api/agent")
	agent.Use(middleware.AuthMiddleware())
	agent.Use(middleware.AgentMiddleware())
	{
		agent.GET("/customers", handlers.GetCustomers(db))
		agent.POST("/customers", handlers.CreateCustomer(db))
		agent.PUT("/customers/:id", handlers.UpdateCustomer(db))
		agent.DELETE("/customers/:id", handlers.DeleteCustomer(db))
		agent.GET("/vehicles", handlers.GetVehicles(db))
		agent.POST("/vehicles", handlers.CreateVehicle(db))
		agent.PUT("/vehicles/:id", handlers.UpdateVehicle(db))
		agent.POST("/drivers", handlers.AssignDriver(db))
		agent.POST("/driving-history", handlers.AddDrivingHistory(db))
		agent.POST("/insurance-history", handlers.AddInsuranceHistory(db))
		agent.GET("/quotes", handlers.GetAgentQuotes(db))
		agent.POST("/quotes", handlers.CreateQuote(db))
		agent.PUT("/quotes/:id", handlers.UpdateQuote(db))
		agent.POST("/quote-line-items", handlers.AddQuoteLineItem(db))
		agent.GET("/coverages", handlers.GetCoverages(db))
		agent.GET("/providers", handlers.GetAgentProviders(db))
		agent.POST("/policies", handlers.BindPolicy(db))
		agent.GET("/policies", handlers.GetPolicies(db))
		agent.GET("/reports", handlers.GetAgentReport(db))
		agent.POST("/auth/logout", handlers.Logout(db))
		// Profile management for agent
		agent.GET("/profiles/users", handlers.GetUsersByRole(db))
		agent.GET("/profiles/user/:id", handlers.GetUserProfile(db))
	}

	// Start the server
	router.Run(":8081")
}
