package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/routes"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/file"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/item"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemDetail"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemDocument"
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
	itemDetailRepository := itemDetail.NewRepository(db)
	itemDocumentRepository := itemDocument.NewRepository(db)

	// Services
	objectTypeService := objectType.NewService(objectTypeRepository)
	organizerService := organizer.NewService(organizerRepository)
	fileService := file.NewService(fileRepository)
	itemService := item.NewService(itemRepository, fileService)
	itemDetailService := itemDetail.NewService(itemDetailRepository)
	itemDocumentService := itemDocument.NewService(itemDocumentRepository)

	// Handlers
	objectTypeHandler := handlers.NewObjectTypeHandler(objectTypeService)
	organizerHandler := handlers.NewOrganizerHandler(organizerService)
	itemHandler := handlers.NewItemHandler(itemService)
	itemDetailHandler := handlers.NewItemDetailHandler(itemDetailService)
	itemDocumentHandler := handlers.NewItemDocumentHandler(itemDocumentService)

	routeConfig := routes.RouteConfig{
		App:                 app,
		ObjectTypeHandler:   objectTypeHandler,
		OrganizerHandler:    organizerHandler,
		ItemHandler:         itemHandler,
		ItemDetailHandler:   itemDetailHandler,
		ItemDocumentHandler: itemDocumentHandler,
	}
	routeConfig.Setup()

	return app, nil
}
