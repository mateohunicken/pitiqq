package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mateo/finanzas-app/config"
	"github.com/mateo/finanzas-app/database"
	"github.com/mateo/finanzas-app/handlers"
	"github.com/mateo/finanzas-app/middleware"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a base de datos
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✓ Connected to PostgreSQL")

	// Crear tablas si no existen
	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	log.Println("✓ Database schema initialized")

	// Configurar Gin
	router := setupRouter(db, cfg)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 Finanzas App Backend running on http://localhost:%s\n", port)
	fmt.Println("📍 API Documentation available at /api/docs")
	fmt.Println("")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(db *database.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Middleware global
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"app":    "Finanzas Domésticas",
		})
	})

	// API v1
	v1 := router.Group("/api/v1")

	// Rutas públicas
	authHandlers := handlers.NewAuthHandler(db)
	v1.POST("/auth/register", authHandlers.Register)
	v1.POST("/auth/login", authHandlers.Login)

	// Rutas protegidas
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Transacciones
	transactionHandlers := handlers.NewTransactionHandler(db)
	protected.GET("/households/:householdId/transactions", transactionHandlers.GetTransactions)
	protected.POST("/households/:householdId/transactions", transactionHandlers.CreateTransaction)
	protected.GET("/households/:householdId/transactions/:transactionId", transactionHandlers.GetTransaction)
	protected.PUT("/households/:householdId/transactions/:transactionId", transactionHandlers.UpdateTransaction)
	protected.DELETE("/households/:householdId/transactions/:transactionId", transactionHandlers.DeleteTransaction)

	// Divisas
	currencyHandlers := handlers.NewCurrencyHandler(db)
	protected.GET("/households/:householdId/currency-exchanges", currencyHandlers.GetExchanges)
	protected.POST("/households/:householdId/currency-exchanges", currencyHandlers.CreateExchange)
	protected.DELETE("/households/:householdId/currency-exchanges/:exchangeId", currencyHandlers.DeleteExchange)

	// Servicios e Impuestos
	serviceHandlers := handlers.NewServiceHandler(db)
	protected.GET("/households/:householdId/recurring-expenses", serviceHandlers.GetRecurringExpenses)
	protected.POST("/households/:householdId/recurring-expenses", serviceHandlers.CreateRecurringExpense)
	protected.PUT("/households/:householdId/recurring-expenses/:expenseId", serviceHandlers.UpdateRecurringExpense)
	protected.DELETE("/households/:householdId/recurring-expenses/:expenseId", serviceHandlers.DeleteRecurringExpense)

	// Categorías
	categoryHandlers := handlers.NewCategoryHandler(db)
	protected.GET("/households/:householdId/categories", categoryHandlers.GetCategories)
	protected.POST("/households/:householdId/categories", categoryHandlers.CreateCategory)
	protected.PUT("/households/:householdId/categories/:categoryId", categoryHandlers.UpdateCategory)
	protected.DELETE("/households/:householdId/categories/:categoryId", categoryHandlers.DeleteCategory)

	// Reportes
	reportHandlers := handlers.NewReportHandler(db)
	protected.GET("/households/:householdId/reports/summary", reportHandlers.GetSummary)
	protected.GET("/households/:householdId/reports/income-expenses", reportHandlers.GetIncomeExpensesReport)
	protected.GET("/households/:householdId/reports/currency-balance", reportHandlers.GetCurrencyBalance)
	protected.GET("/households/:householdId/reports/monthly-breakdown", reportHandlers.GetMonthlyBreakdown)

	// Usuarios
	userHandlers := handlers.NewUserHandler(db)
	protected.GET("/households/:householdId/users", userHandlers.GetHouseholdUsers)
	protected.GET("/me", userHandlers.GetCurrentUser)

	return router
}
