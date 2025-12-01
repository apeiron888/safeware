package main

import (
	"log"
	"os"

	"github.com/a2sv/safeware/internal/audit"
	"github.com/a2sv/safeware/internal/auth"
	"github.com/a2sv/safeware/internal/config"
	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/email"
	"github.com/a2sv/safeware/internal/handlers"
	"github.com/a2sv/safeware/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	database.ConnectMongoDB(cfg.Database.URI, cfg.Database.Database)
	defer database.Close()

	// Set Gin mode
	if cfg.Server.GinMode != "" {
		gin.SetMode(cfg.Server.GinMode)
	}

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.RefreshSecret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	emailService := email.NewEmailService(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPass,
		cfg.Email.SMTPFrom,
		cfg.Email.FrontendURL,
	)
	auditService := audit.NewAuditService(cfg.JWT.Secret) // Using JWT secret as encryption key for MVP

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(jwtService, emailService)
	warehouseHandler := handlers.NewWarehouseHandler()
	itemHandler := handlers.NewItemHandler()
	// roleHandler := handlers.NewRoleHandler()
	// permissionHandler := handlers.NewPermissionHandler()
	// userHandler := handlers.NewUserHandler()
	managerHandler := handlers.NewManagerHandler()
	auditHandler := handlers.NewAuditHandler(auditService)

	// Seed default permissions (run once on startup)
	if err := database.SeedDefaultPermissions(); err != nil {
		log.Printf("Warning: Failed to seed permissions: %v", err)
	}

	// Initialize router
	router := gin.Default()

	// Global Middleware
	router.Use(middleware.AuditMiddleware(auditService))

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"message":  "SIMS API with MongoDB is running!",
			"database": "MongoDB Atlas",
			"version":  "1.0.0-mvp",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// Common routes (all authenticated users)
			protected.GET("/users/me", authHandler.GetProfile)
			protected.POST("/auth/logout", authHandler.Logout)

			// MANAGER ROUTES (Full Access)
			manager := protected.Group("/manager")
			manager.Use(middleware.RoleEnforcementMiddleware("Manager"))
			{
				// Employee Management
				manager.POST("/staff/create", managerHandler.CreateEmployee("Staff"))
				manager.POST("/supervisor/create", managerHandler.CreateEmployee("Supervisor"))
				manager.POST("/auditor/create", managerHandler.CreateEmployee("Auditor"))

				// Warehouse Management
				manager.POST("/warehouse/create", warehouseHandler.Create)
				manager.PATCH("/warehouse/update/:id", warehouseHandler.Update)
				manager.DELETE("/warehouse/delete/:id", warehouseHandler.Delete)
				manager.GET("/summary/warehouses", warehouseHandler.List) // List all

				// Item Management (Global)
				manager.POST("/item/create", itemHandler.Create)
				manager.PUT("/item/update/:id", itemHandler.Update) // Using PUT as per spec
				manager.DELETE("/item/remove/:id", itemHandler.Delete)
				manager.GET("/items/all", itemHandler.List)
				manager.GET("/items/warehouse/:id", itemHandler.List) // Filter by warehouse
			}

			// SUPERVISOR ROUTES (Warehouse Bound + Time Restricted)
			supervisor := protected.Group("/supervisor")
			supervisor.Use(middleware.TimeEnforcementMiddleware())
			supervisor.Use(middleware.RoleEnforcementMiddleware("Supervisor"))
			supervisor.Use(middleware.WarehouseEnforcementMiddleware())
			{
				// Staff Management (Own Warehouse)
				supervisor.POST("/staff/create", managerHandler.CreateEmployee("Staff")) // Reusing logic, but need to ensure warehouse restriction
				// Note: managerHandler.CreateEmployee might need adjustment to allow Supervisors to create Staff for THEIR warehouse only

				// Item Management (Own Warehouse)
				supervisor.POST("/item/add", itemHandler.Create)
				supervisor.PUT("/item/update/:id", itemHandler.Update)
				supervisor.DELETE("/item/remove/:id", itemHandler.Delete)
				supervisor.GET("/items", itemHandler.List)
				supervisor.GET("/item/:id", itemHandler.Get)
			}

			// STAFF ROUTES (Warehouse Bound + Time Restricted)
			staff := protected.Group("/staff")
			staff.Use(middleware.TimeEnforcementMiddleware())
			staff.Use(middleware.RoleEnforcementMiddleware("Staff"))
			staff.Use(middleware.WarehouseEnforcementMiddleware())
			{
				staff.POST("/item/add", itemHandler.Create)
				staff.PUT("/item/update/:id", itemHandler.Update)
				staff.DELETE("/item/remove/:id", itemHandler.Delete)
				staff.GET("/items", itemHandler.List)
				staff.GET("/item/:id", itemHandler.Get)
			}

			// AUDITOR ROUTES (Read Only + Time Restricted)
			auditor := protected.Group("/auditor")
			auditor.Use(middleware.TimeEnforcementMiddleware())
			auditor.Use(middleware.RoleEnforcementMiddleware("Auditor"))
			{
				auditor.GET("/warehouses", warehouseHandler.List)
				auditor.GET("/items/warehouse/:id", itemHandler.List)
				auditor.GET("/items/all", itemHandler.List)
				auditor.GET("/audit-logs", auditHandler.List) // View audit logs
			}

			// Manager Audit Logs
			manager.GET("/audit-logs", auditHandler.List)
		}
	}

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Println("==========================================")
	log.Println("🚀 SIMS Backend Starting...")
	log.Printf("📊 Database: MongoDB Atlas")
	log.Printf("🌐 Server running on port %s", port)
	log.Printf("✅ API available at: http://localhost:%s/api/v1", port)
	log.Printf("💚 Health check: http://localhost:%s/health", port)
	log.Printf("📝 Total endpoints: 27")
	log.Println("==========================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
		os.Exit(1)
	}
}
