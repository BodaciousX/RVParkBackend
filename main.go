package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BodaciousX/RVParkBackend/api"
	"github.com/BodaciousX/RVParkBackend/middleware"
	"github.com/BodaciousX/RVParkBackend/payment"
	"github.com/BodaciousX/RVParkBackend/space"
	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/BodaciousX/RVParkBackend/user"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		dbUser = "dbuser"
	}

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		dbPassword = "dbpassword"
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "RVParkDB"
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5433"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize server with all services
	server, err := initializeServices(db)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func initializeServices(db *sql.DB) (*api.Server, error) {
	// Initialize repositories
	userRepo := user.NewSQLRepository(db)
	tokenRepo := user.NewTokenRepository(db)
	spaceRepo := space.NewSQLRepository(db)
	tenantRepo := tenant.NewSQLRepository(db)
	paymentRepo := payment.NewSQLRepository(db)

	// Initialize services (note the changed order)
	userService := user.NewService(userRepo, tokenRepo)

	// Create space service first, with no dependency on tenant
	spaceService := space.NewService(spaceRepo)

	// Create tenant service with dependency on space
	tenantService := tenant.NewService(tenantRepo, spaceService)

	paymentService := payment.NewService(paymentRepo)

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Create server with all services
	server := api.NewServer(
		userService,
		tenantService,
		spaceService,
		paymentService,
		authMiddleware,
	)

	return server, nil
}
