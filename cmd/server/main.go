package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/joho/godotenv"

	"github.com/elginbrian/FILKOMPLAIN-BE/config"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/app"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/handler"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/repository"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/routes"
)

func main() {
	_ = godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&model.Report{}, &model.User{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}

	repo := repository.NewReportRepository(db)
	service := app.NewReportService(repo)
	reportHandler := handler.NewReportHandler(service)

	authRepo := repository.NewAuthRepository(db)
	authService := app.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	userRepo := repository.NewUserRepository(db)
	userService := app.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	
	handler.SetupSystemInfo(db)

	app := fiber.New()
	
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://filkomplain-fe.vercel.app",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(app, jwtSecret, authHandler, userHandler, reportHandler)

	log.Fatal(app.Listen(":3000"))
}