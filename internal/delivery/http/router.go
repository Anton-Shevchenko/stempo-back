package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stempo/backend/internal/infrastructure/middleware"
	"github.com/stempo/backend/internal/usecase"
)

func SetupRouter(
	authUsecase usecase.AuthUsecase,
	businessUsecase usecase.BusinessUsecase,
	programUsecase usecase.BonusProgramUsecase,
	cardUsecase usecase.LoyaltyCardUsecase,
	categoryUsecase usecase.CategoryUsecase,
	cityUsecase usecase.CityUsecase,
	employeeUsecase usecase.EmployeeUsecase,
	redemptionUsecase usecase.BonusRedemptionUsecase,
	qrCodeUsecase usecase.QRCodeUsecase,
) *gin.Engine {
	router := gin.New()
	
	// Add logger middleware
	router.Use(gin.Logger())
	
	// Add recovery middleware with JSON error response
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		c.Abort()
	}))
	
	// Ensure all responses are JSON
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	})

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	api := router.Group("/api")
	{
		authHandler := NewAuthHandler(authUsecase)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google", authHandler.GoogleLogin)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
			auth.PUT("/me", middleware.AuthMiddleware(), authHandler.UpdateProfile)
			auth.PUT("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword)
			auth.POST("/set-password", authHandler.SetPassword)
		}

		businessHandler := NewBusinessHandler(businessUsecase)
		businesses := api.Group("/businesses")
		{
			businesses.GET("", businessHandler.GetAll)
			businesses.GET("/featured", businessHandler.GetFeatured)
			businesses.GET("/newest", businessHandler.GetNewest)
			businesses.GET("/my-business", middleware.AuthMiddleware(), businessHandler.GetMyBusiness)
			// Specific route must come before parameterized route
			businesses.GET("/by-category/:category", businessHandler.GetByCategory)
			businesses.GET("/:id", businessHandler.GetByID)
			businesses.POST("", middleware.AuthMiddleware(), businessHandler.Create)
			businesses.PUT("/:id", middleware.AuthMiddleware(), businessHandler.Update)
			businesses.DELETE("/:id", middleware.AuthMiddleware(), businessHandler.Delete)
		}

		programHandler := NewBonusProgramHandler(programUsecase, businessUsecase)
		programs := api.Group("/bonus-programs")
		{
			// Specific routes must come before parameterized routes
			programs.GET("/public/by-business/:id", programHandler.GetPublicByBusinessID)
			programs.GET("/by-business/:id", middleware.AuthMiddleware(), programHandler.GetByBusinessID)
			programs.GET("", programHandler.GetAll)
			programs.GET("/:id", programHandler.GetByID)
			programs.POST("", middleware.AuthMiddleware(), programHandler.Create)
			programs.PUT("/:id", middleware.AuthMiddleware(), programHandler.Update)
			programs.DELETE("/:id", middleware.AuthMiddleware(), programHandler.Delete)
		}

		cardHandler := NewLoyaltyCardHandler(cardUsecase)
		cards := api.Group("/loyalty-cards")
		{
			cards.GET("/by-business/:businessId", middleware.AuthMiddleware(), cardHandler.GetByUserIDAndBusinessID)
			cards.GET("", middleware.AuthMiddleware(), cardHandler.GetByUserID)
			cards.GET("/:id", cardHandler.GetByID)
			cards.POST("", middleware.AuthMiddleware(), cardHandler.Create)
			cards.POST("/:id/stamps", middleware.AuthMiddleware(), cardHandler.AddStamp)
			cards.PUT("/:id", middleware.AuthMiddleware(), cardHandler.Update)
		}

		employeeHandler := NewEmployeeHandler(employeeUsecase)
		employees := api.Group("/employees")
		{
			employees.GET("/my-businesses", middleware.AuthMiddleware(), employeeHandler.GetMyBusinesses)
			employees.GET("/business/:businessId", middleware.AuthMiddleware(), employeeHandler.GetEmployees)
			employees.POST("/business/:businessId", middleware.AuthMiddleware(), employeeHandler.AddEmployee)
			employees.DELETE("/business/:businessId/:employeeId", middleware.AuthMiddleware(), employeeHandler.RemoveEmployee)
		}

		redemptionHandler := NewBonusRedemptionHandler(redemptionUsecase)
		redemptions := api.Group("/bonus-redemptions")
		{
			redemptions.POST("/card/:cardId/generate", middleware.AuthMiddleware(), redemptionHandler.GenerateCode)
			redemptions.POST("/redeem", middleware.AuthMiddleware(), redemptionHandler.Redeem)
			redemptions.GET("/card/:cardId", middleware.AuthMiddleware(), redemptionHandler.GetByCard)
		}

		categoryHandler := NewCategoryHandler(categoryUsecase)
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:slug", categoryHandler.GetBySlug)
		}

		cityHandler := NewCityHandler(cityUsecase)
		cities := api.Group("/cities")
		{
			cities.GET("", cityHandler.GetAll)
			cities.GET("/:id", cityHandler.GetByID)
			cities.POST("", middleware.AuthMiddleware(), cityHandler.Create)
			cities.PUT("/:id", middleware.AuthMiddleware(), cityHandler.Update)
			cities.DELETE("/:id", middleware.AuthMiddleware(), cityHandler.Delete)
		}

		qrCodeHandler := NewQRCodeHandler(qrCodeUsecase)
		qrCodes := api.Group("/qr-codes")
		{
			qrCodes.POST("/generate", middleware.AuthMiddleware(), qrCodeHandler.Generate)
			qrCodes.GET("/validate/:code", qrCodeHandler.Validate)
			qrCodes.GET("/program/:programId", middleware.AuthMiddleware(), qrCodeHandler.GetByProgramID)
			qrCodes.GET("/business/:businessId", middleware.AuthMiddleware(), qrCodeHandler.GetByBusinessID)
			qrCodes.DELETE("/:id", middleware.AuthMiddleware(), qrCodeHandler.Delete)
		}

		adminHandler := NewAdminHandler(businessUsecase, programUsecase)
		admin := api.Group("/admin")
		{
			admin.GET("/businesses/pending", adminHandler.GetPendingBusinesses)
			admin.POST("/businesses/:id/approve", adminHandler.ApproveBusiness)
			admin.POST("/businesses/:id/reject", adminHandler.RejectBusiness)
			admin.PUT("/businesses/:id", adminHandler.UpdateBusiness)
			admin.GET("/bonus-programs/pending", adminHandler.GetPendingPrograms)
			admin.POST("/bonus-programs/:id/approve", adminHandler.ApproveProgram)
			admin.POST("/bonus-programs/:id/reject", adminHandler.RejectProgram)
		}
	}

	// Handle 404 errors with JSON response
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found", "path": c.Request.URL.Path})
	})

	return router
}
