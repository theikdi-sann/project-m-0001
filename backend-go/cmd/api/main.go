package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theikdi-sann/qr-restaurant-api/internal/delivery/http"
	"github.com/theikdi-sann/qr-restaurant-api/internal/delivery/http/middleware"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres"
	"github.com/theikdi-sann/qr-restaurant-api/internal/repository/postgres/db"
	"github.com/theikdi-sann/qr-restaurant-api/internal/usecase"
	"github.com/theikdi-sann/qr-restaurant-api/internal/worker"
)

func main() {
	// 1. Configuration
	dbSource := os.Getenv("DATABASE_URL")
	if dbSource == "" {
		// testing
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

	sessionRepo := postgres.NewDiningSessionRepository(connPool)
	sessionTypeRepo := postgres.NewSessionTypeRepository(queries)
	menuRepo := postgres.NewMenuItemRepository(queries)
	orderRepo := postgres.NewOrderRepository(connPool)

	sessionUsecase := usecase.NewDiningSessionUsecase(sessionRepo, sessionTypeRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, sessionRepo, menuRepo)

	sessionHandler := http.NewDiningSessionHandler(sessionUsecase)
	orderHandler := http.NewOrderHandler(orderUsecase)
	menuHandler := http.NewMenuHandler(menuRepo)

	// 4. HTTP Server
	r := gin.Default()

	// CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Simple Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	sessionHandler.RegisterPublicRoutes(r)
	menuHandler.RegisterPublicRoutes(r)

	// Protected Routes
	protected := r.Group("/")

	jwksURL := os.Getenv("SUPABASE_JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://127.0.0.1:54321/auth/v1/.well-known/jwks.json"
	}
	protected.Use(middleware.NewAuthMiddleware(jwksURL))

	sessionHandler.RegisterRoutes(protected)
	orderHandler.RegisterRoutes(protected)
	menuHandler.RegisterProtectedRoutes(protected)

	// Background Worker
	cleanupWorker := worker.NewSessionCleanupWorker(sessionRepo)
	cleanupWorker.Start(context.Background(), 1*time.Minute)

	// 5. Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s | current time: %s", port, time.Now())
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
