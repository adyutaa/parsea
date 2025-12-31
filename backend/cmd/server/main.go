package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/adyutaa/parsea/internal/handler"
	"github.com/adyutaa/parsea/internal/infrastructure/llm"
	"github.com/adyutaa/parsea/internal/infrastructure/vectordb"
	"github.com/adyutaa/parsea/internal/repository"
	"github.com/adyutaa/parsea/internal/service"
	"github.com/adyutaa/parsea/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Initialize database connection
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	fmt.Println("✅ Connected to Supabase PostgreSQL!")

	// Initialize Redis connection
	rdb, err := initRedis()
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}
	fmt.Println("✅ Connected to Redis Cloud!")

	// Initialize OpenAI client
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY not set in environment")
	}
	llmClient := llm.NewOpenAIClient(openaiKey)
	fmt.Println("✅ Connected to OpenAI!")

	// Initialize Qdrant (optional - will fallback if not configured)
	var contextService *service.ContextService
	qdrantClient, err := vectordb.NewQdrantClient()
	if err != nil {
		log.Printf("⚠️  Qdrant not available: %v (will use fallback context)\n", err)
		contextService = nil
	} else {
		fmt.Println("✅ Connected to Qdrant Cloud!")
		contextService = service.NewContextService(qdrantClient, llmClient)
	}

	// Create uploads directory
	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Initialize repositories
	docRepo := repository.NewDocumentRepository(db)
	evalRepo := repository.NewEvaluationRepository(db)

	// Initialize services
	docService := service.NewDocumentService(docRepo, uploadPath)
	evalService := service.NewEvaluationService(evalRepo, docRepo, rdb)

	// Initialize handlers
	docHandler := handler.NewDocumentHandler(docService)
	evalHandler := handler.NewEvaluationHandler(evalService)

	// Start background worker
	evalWorker := worker.NewEvaluationWorker(rdb, evalRepo, docRepo, llmClient, contextService)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	go evalWorker.Start(workerCtx)
	fmt.Println("✅ Background worker started!")

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "healthy",
			"database": "connected",
			"redis":    "connected",
			"llm":      "connected",
			"qdrant":   contextService != nil,
		})
	})

	// API routes
	r.POST("/upload", docHandler.Upload)
	r.POST("/evaluate", evalHandler.Evaluate)
	r.GET("/result", evalHandler.GetResult)
	r.GET("/queue/status", evalHandler.GetQueueStatus)

	// Debug endpoint
	r.GET("/debug/jobs", func(c *gin.Context) {
		var jobs []map[string]any
		err := db.Table("evaluation_jobs").Find(&jobs).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"total_jobs": len(jobs),
			"jobs":       jobs,
		})
	})

	// Start server
	port := getEnv("PORT", "8080")

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🚀 Server starting on http://localhost:%s\n", port)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n📋 Available endpoints:")
	fmt.Println("  GET    /health              - Health check")
	fmt.Println("  POST   /upload              - Upload CV and Project Report")
	fmt.Println("  POST   /evaluate            - Start evaluation job")
	fmt.Println("  GET    /result?id=job_id    - Get evaluation result")
	fmt.Println("  GET    /queue/status        - Get queue status")
	fmt.Println()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("\n👋 Shutting down server...")
		workerCancel()
		os.Exit(0) // Force exit after 1 second
	}()

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDatabase initializes PostgreSQL connection with proper settings
func initDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set in environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		PrepareStmt:                              true, // Enable prepared statements for better performance
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Reasonable connection pooling settings
	sqlDB.SetMaxIdleConns(10)                  // 10 idle connections
	sqlDB.SetMaxOpenConns(25)                  // 25 max open connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 30 minute connection lifetime
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // 5 minute idle timeout

	// Test connection
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return db, nil
}

// initRedis initializes Redis connection
func initRedis() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" || port == "" || password == "" {
		return nil, fmt.Errorf("Redis configuration not set")
	}

	if username == "" {
		username = "default"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Username: username,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
