package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stempo/backend/internal/delivery/http"
	"github.com/stempo/backend/internal/infrastructure/database"
	"github.com/stempo/backend/internal/infrastructure/email"
	"github.com/stempo/backend/internal/infrastructure/oauth"
	"github.com/stempo/backend/internal/repository/postgres"
	"github.com/stempo/backend/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			runMigrations()
			return
		case "seed":
			runSeeds()
			return
		}
	}

	runServer()
}

func runMigrations() {
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Migrations completed successfully")
}

func runSeeds() {
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}
}

func runServer() {
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userRepo := postgres.NewUserRepository(db)
	businessRepo := postgres.NewBusinessRepository(db)
	programRepo := postgres.NewBonusProgramRepository(db)
	cardRepo := postgres.NewLoyaltyCardRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	cityRepo := postgres.NewCityRepository(db)
	employeeRepo := postgres.NewEmployeeRepository(db)
	redemptionRepo := postgres.NewBonusRedemptionRepository(db)
	qrCodeRepo := postgres.NewQRCodeRepository(db)

	emailSvc := email.NewMailjetService()

	googleVerifier := oauth.NewGoogleTokenVerifier()
	authUsecase := usecase.NewAuthUsecase(userRepo, googleVerifier)
	businessUsecase := usecase.NewBusinessUsecase(businessRepo)
	programUsecase := usecase.NewBonusProgramUsecase(programRepo, businessRepo)
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo, businessRepo, userRepo, emailSvc)
	cardUsecase := usecase.NewLoyaltyCardUsecase(cardRepo, businessRepo, employeeUsecase)
	redemptionUsecase := usecase.NewBonusRedemptionUsecase(redemptionRepo, cardRepo, businessRepo, employeeUsecase)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	cityUsecase := usecase.NewCityUsecase(cityRepo)
	qrCodeUsecase := usecase.NewQRCodeUsecase(qrCodeRepo, programRepo, businessRepo)

	router := http.SetupRouter(authUsecase, businessUsecase, programUsecase, cardUsecase, categoryUsecase, cityUsecase, employeeUsecase, redemptionUsecase, qrCodeUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
