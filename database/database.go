package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect initializes a connection to the database
func Connect() *sql.DB {
	// Get database configuration from environment variables
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "example")

	// First connect to postgres database to check if our database exists
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'insurance')").Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	// If database doesn't exist, create it
	if !exists {
		_, err = db.Exec("CREATE DATABASE insurance")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created database: insurance")
		db.Close()
	} else {
		db.Close()
	}

	// Connect to the insurance database
	connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=insurance sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database!")

	// Initialize tables if they don't exist
	initializeTables(db)

	// Seed initial data (create default MasterAdmin if none exists)
	seedInitialData(db)

	return db
}

// initializeTables checks if tables exist and creates them if they don't
func initializeTables(db *sql.DB) {
	// Check if the main tables exist by checking for the Users table
	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users')").Scan(&tableExists)
	if err != nil {
		log.Printf("Error checking if tables exist: %v", err)
		return
	}

	if !tableExists {
		fmt.Println("Tables don't exist. Initializing database schema...")

		// Read the schema.sql file
		schemaPath := filepath.Join("database", "schema.sql")
		schemaBytes, err := os.ReadFile(schemaPath)
		if err != nil {
			log.Printf("Error reading schema.sql: %v", err)
			return
		}

		// Execute the schema
		_, err = db.Exec(string(schemaBytes))
		if err != nil {
			log.Printf("Error executing schema: %v", err)
			return
		}

		fmt.Println("Database schema initialized successfully!")
	} else {
		fmt.Println("Database tables already exist.")

		// Check if ENUM types exist and create them if they don't
		ensureEnumTypes(db)
	}
}

// ensureEnumTypes creates ENUM types if they don't exist
func ensureEnumTypes(db *sql.DB) {
	enumQueries := []string{
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN CREATE TYPE user_role AS ENUM ('MasterAdmin', 'AgencyAdmin', 'LocationAdmin', 'Agent', 'Customer'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'quote_status') THEN CREATE TYPE quote_status AS ENUM ('Draft', 'Presented', 'Bound'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'policy_status') THEN CREATE TYPE policy_status AS ENUM ('Active', 'Expired', 'Cancelled'); END IF; END $$;",
	}

	for _, query := range enumQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating ENUM type: %v", err)
		}
	}
}

// seedInitialData creates initial required data like default MasterAdmin account
func seedInitialData(db *sql.DB) {
	// Check if any MasterAdmin users exist
	var masterAdminExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE Role = 'MasterAdmin')").Scan(&masterAdminExists)
	if err != nil {
		log.Printf("Error checking for MasterAdmin users: %v", err)
		return
	}

	if !masterAdminExists {
		// Create default MasterAdmin account
		fmt.Println("No MasterAdmin account found. Creating default MasterAdmin...")

		// Get admin credentials from environment variables with fallback defaults
		email := getEnvOrDefault("ADMIN_EMAIL", "admin@insurance.com")
		password := getEnvOrDefault("ADMIN_PASSWORD", "admin123")
		firstName := getEnvOrDefault("ADMIN_FIRST_NAME", "System")
		lastName := getEnvOrDefault("ADMIN_LAST_NAME", "Administrator")

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return
		}

		// Insert the default MasterAdmin user
		_, err = db.Exec(`
			INSERT INTO Users (FirstName, LastName, Email, PasswordHash, Role) 
			VALUES ($1, $2, $3, $4, $5)
		`, firstName, lastName, email, string(hashedPassword), "MasterAdmin")

		if err != nil {
			log.Printf("Error creating default MasterAdmin: %v", err)
			return
		}

		fmt.Printf("✅ Default MasterAdmin account created!\n")
		fmt.Printf("   Email: %s\n", email)
		fmt.Printf("   Password: %s\n", password)
		fmt.Printf("   ⚠️  IMPORTANT: Change this password after first login!\n")
	} else {
		fmt.Println("MasterAdmin account(s) already exist.")
	}
}
