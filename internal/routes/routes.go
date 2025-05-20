package routes

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/handler"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	jwtSecret string,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	reportHandler *handler.ReportHandler,
) {
	app.Get("/", handler.WelcomeHandler)
	
	apiV1 := app.Group("/api/v1")
	
	apiV1.Post("/register", authHandler.Register)
	apiV1.Post("/login", authHandler.Login)
	apiV1.Post("/admin/register", authHandler.RegisterAdmin)
	apiV1.Post("/admin/login", authHandler.LoginAdmin)

	protected := apiV1.Group("/", middleware.JWTProtected(jwtSecret))
	protected.Post("/refresh-token", authHandler.RefreshToken)

	protected.Get("/profile", userHandler.GetProfile)
	protected.Put("/profile", userHandler.UpdateProfile)

	protected.Get("/reports", reportHandler.GetReports)
	protected.Get("/reports/:id", reportHandler.GetReport)
	protected.Post("/reports", reportHandler.CreateReport)
	protected.Put("/reports/:id", reportHandler.UpdateReport)
	protected.Delete("/reports/:id", reportHandler.DeleteReport)
	
	admin := protected.Group("/", middleware.AdminOnly)
	admin.Patch("/reports/:id/status", reportHandler.ResolveReportStatus)
	admin.Patch("/reports/:id/reply", reportHandler.ReplyReport)
}
