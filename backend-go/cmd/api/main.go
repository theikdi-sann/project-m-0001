package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theikdi-sann/qr-restaurant-api/internal/delivery/http"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
	"github.com/theikdi-sann/qr-restaurant-api/internal/usecase"
)

func main() {
	// 1. Configuration
	dbSource := os.Getenv("DATABASE_URL")
	if dbSource == "" {
		dbSource = "postgresql://postgres:postgres@127.0.0.1:54328/postgres"
	}

	// 2. Database Connection
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer connPool.Close()

	if err := connPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to Database")

	// 3. Init Layers
	queries := db.New(connPool)
	
	sessionRepo := postgres.NewDiningSessionRepository(queries)
	sessionTypeRepo := postgres.NewSessionTypeRepository(queries)

	sessionUsecase := usecase.NewDiningSessionUsecase(sessionRepo, sessionTypeRepo)
	sessionHandler := http.NewDiningSessionHandler(sessionUsecase)

	// 4. HTTP Server
	r := gin.Default()
	
	// Simple Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	sessionHandler.RegisterRoutes(r)

	// 5. Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}