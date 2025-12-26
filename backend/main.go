package main

import (
	"log"
	"os"
	"task-management/config"
	"task-management/handlers"
	"task-management/middleware"
	"task-management/repositories"
	"task-management/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	db := config.InitDB()

	// Auto-migrate models (optional, for development)
	// config.AutoMigrate(db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	issueRepo := repositories.NewIssueRepository(db)
	assignmentRepo := repositories.NewAssignmentRepository(db)
	calendarRepo := repositories.NewCalendarRepository(db)
	statusRepo := repositories.NewStatusRepository(db)
	attachmentRepo := repositories.NewAttachmentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	orgService := services.NewOrganizationService(orgRepo)
	teamService := services.NewTeamService(teamRepo, userRepo)
	issueService := services.NewIssueService(issueRepo, statusRepo)
	assignmentService := services.NewAssignmentService(assignmentRepo, issueRepo, userRepo)
	calendarService := services.NewCalendarService(calendarRepo)
	permissionService := services.NewPermissionService(teamRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	orgHandler := handlers.NewOrganizationHandler(orgService, permissionService)
	teamHandler := handlers.NewTeamHandler(teamService, permissionService)
	issueHandler := handlers.NewIssueHandler(issueService, assignmentService, permissionService)
	calendarHandler := handlers.NewCalendarHandler(calendarService, permissionService)
	statusHandler := handlers.NewStatusHandler(statusRepo, permissionService)

	// Initialize storage service (optional - for file attachments)
	var attachmentHandler *handlers.AttachmentHandler
	storageService, err := services.NewStorageService()
	if err != nil {
		log.Printf("Warning: Storage service not configured: %v", err)
		log.Println("File attachment feature will be disabled.")
	} else {
		attachmentHandler = handlers.NewAttachmentHandler(attachmentRepo, storageService)
		log.Println("Storage service (Cloudflare R2) initialized successfully")
	}

	// Setup Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// Change password (needs auth)
	api.POST("/auth/change-password", authHandler.ChangePassword)
	{
		// Organizations
		orgs := api.Group("/organizations")
		{
			orgs.GET("", orgHandler.List)
			orgs.POST("", orgHandler.Create)
			orgs.GET("/:id", orgHandler.GetByID)
			orgs.PUT("/:id", orgHandler.Update)
			orgs.DELETE("/:id", orgHandler.Delete)
		}

		// Teams
		teams := api.Group("/teams")
		{
			teams.GET("", teamHandler.List)
			teams.POST("", teamHandler.Create)
			teams.GET("/:id", teamHandler.GetByID)
			teams.PUT("/:id", teamHandler.Update)
			teams.DELETE("/:id", teamHandler.Delete)
			teams.GET("/:id/members", teamHandler.GetMembers)
			teams.POST("/:id/members", teamHandler.AddMember)
			teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
		}

		// Issue Statuses
		statuses := api.Group("/statuses")
		{
			statuses.GET("", statusHandler.GetByOrganization)
			statuses.POST("", statusHandler.Create)
			statuses.PUT("/:id", statusHandler.Update)
			statuses.DELETE("/:id", statusHandler.Delete)
		}

		// Issues
		issues := api.Group("/issues")
		{
			issues.GET("", issueHandler.List)
			issues.POST("", issueHandler.Create)
			issues.GET("/:id", issueHandler.GetByID)
			issues.PUT("/:id", issueHandler.Update)
			issues.DELETE("/:id", issueHandler.Delete)
			issues.POST("/:id/assign", issueHandler.Assign)
			issues.POST("/:id/status", issueHandler.UpdateStatus)
			issues.POST("/:id/hold", issueHandler.Hold)
			issues.POST("/:id/resume", issueHandler.Resume)
			issues.GET("/:id/activities", issueHandler.GetActivities)
			issues.POST("/:id/worklog", issueHandler.LogWork)

			// Attachments (if storage service is configured)
			if attachmentHandler != nil {
				issues.GET("/:id/attachments", attachmentHandler.List)
				issues.POST("/:id/attachments", attachmentHandler.Upload)
			}
		}

		// Attachments (standalone routes)
		if attachmentHandler != nil {
			attachments := api.Group("/attachments")
			{
				attachments.GET("/:id/download", attachmentHandler.Download)
				attachments.DELETE("/:id", attachmentHandler.Delete)
			}
		}

		// Calendar
		calendar := api.Group("/calendar")
		{
			calendar.GET("", calendarHandler.GetCalendar)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
