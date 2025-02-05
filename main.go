package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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

// initializeDatabase reads and executes the init.sql file
func initializeDatabase(db *sql.DB) error {
	log.Println("Starting database initialization...")

	// Read init.sql file
	initSQL, err := os.ReadFile("init.sql")
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

func ensureStaffExists(userService user.Service) error {
	// Get staff credentials from environment variables or use defaults
	staffEmail := os.Getenv("STAFF_EMAIL")
	if staffEmail == "" {
		staffEmail = "staff@rvpark.com"
	}

	staffPassword := os.Getenv("STAFF_PASSWORD")
	if staffPassword == "" {
		// Generate a random password if not provided
		staffPassword = fmt.Sprintf("staff%d", time.Now().Unix())
	}

	// Check if staff exists by trying to get by email
	_, err := userService.GetUserByEmail(staffEmail)
	if err == nil {
		// Staff exists, just log the credentials
		log.Println("----------------------------------------")
		log.Println("Existing staff account found:")
		log.Printf("Email: %s\n", staffEmail)
		log.Printf("Password: %s\n", staffPassword)
		log.Println("Please save these credentials securely!")
		log.Println("----------------------------------------")
		return nil
	}

	// Create new staff user
	staffUser := user.User{
		ID:        uuid.New().String(),
		Email:     staffEmail,
		Username:  "staff",
		Role:      user.RoleStaff,
		CreatedAt: time.Now(),
	}

	err = userService.CreateUser(staffUser, staffPassword)
	if err != nil {
		return fmt.Errorf("failed to create staff user: %v", err)
	}

	// Print the credentials
	log.Println("----------------------------------------")
	log.Println("New staff account created successfully:")
	log.Printf("Email: %s\n", staffEmail)
	log.Printf("Password: %s\n", staffPassword)
	log.Println("Please save these credentials securely!")
	log.Println("----------------------------------------")

	return nil
}

func ensureAdminExists(userService user.Service) error {
	// Get admin credentials from environment variables or use defaults
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@rvpark.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		// Generate a random password if not provided
		adminPassword = fmt.Sprintf("admin%d", time.Now().Unix())
	}

	// Check if admin exists by trying to get by email
	_, err := userService.GetUserByEmail(adminEmail)
	if err == nil {
		// Admin exists, just log the credentials
		log.Println("----------------------------------------")
		log.Println("Existing admin account found:")
		log.Printf("Email: %s\n", adminEmail)
		log.Printf("Password: %s\n", adminPassword)
		log.Println("Please save these credentials securely!")
		log.Println("----------------------------------------")
		return nil
	}

	// Create new admin user
	adminUser := user.User{
		ID:        uuid.New().String(),
		Email:     adminEmail,
		Username:  "admin",
		Role:      user.RoleAdmin,
		CreatedAt: time.Now(),
	}

	err = userService.CreateUser(adminUser, adminPassword)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	// Print the credentials
	log.Println("----------------------------------------")
	log.Println("New admin account created successfully:")
	log.Printf("Email: %s\n", adminEmail)
	log.Printf("Password: %s\n", adminPassword)
	log.Println("Please save these credentials securely!")
	log.Println("----------------------------------------")

	return nil
}

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://dbuser:dbpassword@localhost:5433/RVParkDB?sslmode=disable"
	}

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

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

	// Ensure admin and staff users exist with known credentials
	if err := ensureAdminExists(userService); err != nil {
		log.Fatalf("Failed to ensure admin exists: %v", err)
	}

	if err := ensureStaffExists(userService); err != nil {
		log.Fatalf("Failed to ensure staff exists: %v", err)
	}

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Initialize server with all services including payment
	server := api.NewServer(
		userService,
		tenantService,
		spaceService,
		paymentService,
		authMiddleware,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Database initialization and checks completed successfully")
	log.Printf("Server running on http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, server.Mux))
}
