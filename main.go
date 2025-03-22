package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BodaciousX/RVParkBackend/api"
	"github.com/BodaciousX/RVParkBackend/middleware"
	"github.com/BodaciousX/RVParkBackend/payment"
	"github.com/BodaciousX/RVParkBackend/space"
	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/BodaciousX/RVParkBackend/user"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// getDBConfig returns the database connection string with appropriate SSL settings
func getDBConfig() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Always use SSL in production (Render requirement)
	if os.Getenv("GO_ENV") != "development" {
		if !strings.Contains(dbURL, "sslmode=") {
			if strings.Contains(dbURL, "?") {
				dbURL += "&sslmode=require"
			} else {
				dbURL += "?sslmode=require"
			}
		}
	}

	return dbURL
}

// initializeDatabase reads and executes the init.sql file
func initializeDatabase(db *sql.DB) error {
	log.Println("Starting database initialization...")

	// Read init.sql file from docker folder
	initSQL, err := os.ReadFile("docker/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql: %v", err)
	}

	// Execute the SQL statements
	_, err = db.Exec(string(initSQL))
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %v", err)
	}

	log.Println("Database initialization completed successfully")
	return nil
}

// checkDatabaseConnection attempts to connect to the database with retries
func checkDatabaseConnection(db *sql.DB) error {
	maxRetries := 30
	retryInterval := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		err := db.Ping()
		if err == nil {
			log.Printf("Successfully connected to database after %d attempts", i+1)
			return nil
		}
		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

// checkRequiredTables verifies that all required tables exist
func checkRequiredTables(db *sql.DB) error {
	requiredTables := []string{
		"users",
		"tokens",
		"sections",
		"spaces",
		"tenants",
		"payments",
	}

	for _, table := range requiredTables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			);
		`
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error checking table %s: %v", table, err)
		}
		if !exists {
			return fmt.Errorf("required table %s does not exist", table)
		}
	}
	log.Println("All required database tables are present")
	return nil
}

func ensureUserExists(userService user.Service) error {
	// Get user credentials from environment variables
	userEmail := os.Getenv("USER_EMAIL")
	if userEmail == "" {
		userEmail = "user@rvpark.com"
	}

	userPassword := os.Getenv("USER_PASSWORD")
	if userPassword == "" {
		userPassword = fmt.Sprintf("user%d", time.Now().Unix())
		log.Printf("WARNING: Generated random user password: %s", userPassword)
	}

	// Check if user exists
	_, err := userService.GetUserByEmail(userEmail)
	if err == nil {
		log.Printf("Existing user account found with email: %s\n", userEmail)
		return nil
	}

	// Create new user
	newUser := user.User{
		ID:        uuid.New().String(),
		Email:     userEmail,
		Username:  "user",
		CreatedAt: time.Now(),
	}

	err = userService.CreateUser(newUser, userPassword)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	log.Printf("New user account created with email: %s\n", userEmail)
	return nil
}

func main() {
	// Get database configuration
	dbURL := getDBConfig()
	log.Printf("Attempting to connect to database...")

	// Open database connection with adjusted settings for cloud environment
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Set connection pool settings
	db.SetMaxOpenConns(25) // Render's free tier limit
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check database connection with retries
	log.Println("Checking database connection...")
	if err := checkDatabaseConnection(db); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Initialize database schema
	log.Println("Initializing database schema...")
	if err := initializeDatabase(db); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		// Continue execution as tables might already exist
	}

	// Verify required tables exist
	log.Println("Verifying database tables...")
	if err := checkRequiredTables(db); err != nil {
		log.Fatalf("Database verification failed: %v", err)
	}

	// Initialize repositories
	userRepo := user.NewSQLRepository(db)
	tokenRepo := user.NewTokenRepository(db)
	tenantRepo := tenant.NewSQLRepository(db)
	spaceRepo := space.NewSQLRepository(db)
	paymentRepo := payment.NewSQLRepository(db)

	// Initialize services
	userService := user.NewService(userRepo, tokenRepo)
	tenantService := tenant.NewService(tenantRepo)
	spaceService := space.NewService(spaceRepo, tenantService)
	paymentService := payment.NewService(paymentRepo)

	// Ensure user exists
	if err := ensureUserExists(userService); err != nil {
		log.Fatalf("Failed to ensure user exists: %v", err)
	}

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Initialize server with all services
	server := api.NewServer(
		userService,
		tenantService,
		spaceService,
		paymentService,
		authMiddleware,
	)

	// Get PORT from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT environment variable not set - using default port %s", port)
	}

	// Listen on all interfaces
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Database initialization and checks completed successfully")
	log.Printf("Server starting on port %s", port)

	// Start the server
	if err := http.ListenAndServe(addr, server.Mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
