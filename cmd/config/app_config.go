package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/routes"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/objectType"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) (*fiber.App, error) {
	utils.InitValidator()
	//validator := utils.Validate

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// Repositories
	objectTypeRepository := objectType.NewRepository(db)

	// Services
	objectTypeService := objectType.NewService(objectTypeRepository)

	// Handlers
	objectTypeHandler := handlers.NewObjectTypeHandler(objectTypeService)

	routeConfig := routes.RouteConfig{
		App:               app,
		ObjectTypeHandler: objectTypeHandler,
	}
	routeConfig.Setup()

	return app, nil
}
