package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/routes"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/file"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/item"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/objectType"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/organizer"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) (*fiber.App, error) {
	utils.InitValidator()
	//validator := utils.Validate

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Repositories
	objectTypeRepository := objectType.NewRepository(db)
	organizerRepository := organizer.NewRepository(db)
	fileRepository := file.NewRepository(db)
	itemRepository := item.NewRepository(db)

	// Services
	objectTypeService := objectType.NewService(objectTypeRepository)
	organizerService := organizer.NewService(organizerRepository)
	fileService := file.NewService(fileRepository)
	itemService := item.NewService(itemRepository, fileService)

	// Handlers
	objectTypeHandler := handlers.NewObjectTypeHandler(objectTypeService)
	organizerHandler := handlers.NewOrganizerHandler(organizerService)
	itemHandler := handlers.NewItemHandler(itemService)

	routeConfig := routes.RouteConfig{
		App:               app,
		ObjectTypeHandler: objectTypeHandler,
		OrganizerHandler:  organizerHandler,
		ItemHandler:       itemHandler,
	}
	routeConfig.Setup()

	return app, nil
}
