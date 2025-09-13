package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/routes"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/auction"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/file"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/item"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemDetail"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemDocument"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemGrade"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemThumbnail"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/objectType"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/organizer"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/pic"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/user"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) (*fiber.App, error) {
	utils.InitValidator()
	//validator := utils.Validate

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	app.Use(logger.New())

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
	itemGradeRepository := itemGrade.NewRepository(db)
	itemThumbnailRepository := itemThumbnail.NewRepository(db)
	picRepository := pic.NewRepository(db)
	auctionRepository := auction.NewRepository(db)
	userRepository := user.NewRepository(db)

	// Services
	objectTypeService := objectType.NewService(objectTypeRepository)
	organizerService := organizer.NewService(organizerRepository)
	fileService := file.NewService(fileRepository)
	itemService := item.NewService(itemRepository, fileService)
	itemDetailService := itemDetail.NewService(itemDetailRepository)
	itemDocumentService := itemDocument.NewService(itemDocumentRepository)
	itemGradeService := itemGrade.NewService(itemGradeRepository)
	itemThumbnailService := itemThumbnail.NewService(itemThumbnailRepository, fileService)
	picService := pic.NewService(picRepository)
	auctionService := auction.NewService(auctionRepository)
	userService := user.NewService(userRepository)

	// Handlers
	objectTypeHandler := handlers.NewObjectTypeHandler(objectTypeService)
	organizerHandler := handlers.NewOrganizerHandler(organizerService)
	itemHandler := handlers.NewItemHandler(itemService)
	itemDetailHandler := handlers.NewItemDetailHandler(itemDetailService)
	itemDocumentHandler := handlers.NewItemDocumentHandler(itemDocumentService)
	itemGradeHandler := handlers.NewItemGradeHandler(itemGradeService)
	itemThumbnailHandler := handlers.NewItemThumbnailHandler(itemThumbnailService)
	picHandler := handlers.NewPICHandler(picService)
	auctionHandler := handlers.NewAuctionHandler(auctionService)
	userHandler := handlers.NewUserHandler(userService)

	routeConfig := routes.RouteConfig{
		App:                  app,
		ObjectTypeHandler:    objectTypeHandler,
		OrganizerHandler:     organizerHandler,
		ItemHandler:          itemHandler,
		ItemDetailHandler:    itemDetailHandler,
		ItemDocumentHandler:  itemDocumentHandler,
		ItemGradeHandler:     itemGradeHandler,
		ItemThumbnailHandler: itemThumbnailHandler,
		PIChandler:           picHandler,
		AuctionHandler:       auctionHandler,
		UserHandler:          userHandler,
	}
	routeConfig.Setup()

	return app, nil
}
