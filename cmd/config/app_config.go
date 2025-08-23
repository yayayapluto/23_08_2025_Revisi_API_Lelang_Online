package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) (*fiber.App, error) {
	utils.InitValidator()
	//validator := utils.Validate

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	validator := utils.Validate

	// Repositories

	// Services

	// Handlers

	routeConfig := routes.RouteConfig{
		App: app,
	}
	routeConfig.Setup()

	return app, nil
}
