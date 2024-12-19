package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BodaciousX/RVParkBackend/api"
	"github.com/BodaciousX/RVParkBackend/middleware"
	"github.com/BodaciousX/RVParkBackend/space"
	"github.com/BodaciousX/RVParkBackend/tenant"
	"github.com/BodaciousX/RVParkBackend/user"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://dbuser:dbpassword@localhost:5433/RVParkDB?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize repositories
	userRepo := user.NewSQLRepository(db)
	tokenRepo := user.NewTokenRepository(db)
	tenantRepo := tenant.NewSQLRepository(db)
	spaceRepo := space.NewSQLRepository(db)

	// Initialize services
	userService := user.NewService(userRepo, tokenRepo)
	tenantService := tenant.NewService(tenantRepo)
	spaceService := space.NewService(spaceRepo, tenantService)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Initialize server
	server := api.NewServer(userService, tenantService, spaceService, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Server running on http://localhost%s\n", addr)

	log.Fatal(http.ListenAndServe(addr, server.Mux))
}
