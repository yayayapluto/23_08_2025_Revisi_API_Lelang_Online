package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/middleware"
)

type RouteConfig struct {
	App                  *fiber.App
	ObjectTypeHandler    handlers.ObjectTypeHandler
	OrganizerHandler     handlers.OrganizerHandler
	ItemHandler          handlers.ItemHandler
	ItemDetailHandler    handlers.ItemDetailHandler
	ItemDocumentHandler  handlers.ItemDocumentHandler
	ItemGradeHandler     handlers.ItemGradeHandler
	ItemThumbnailHandler handlers.ItemThumbnailHandler
	PIChandler           handlers.PICHandler
	AuctionHandler       handlers.AuctionHandler
	UserHandler          handlers.UserHandler
}

func (r *RouteConfig) Setup() {
	r.ObjectType()
	r.Organizer()
	r.Item()
	r.PIC()
	r.Auction()
	r.User()
	r.Auth()
}

func (r *RouteConfig) ObjectType() {
	group := r.App.Group("/api/objectTypes")
	group.Get("/", r.ObjectTypeHandler.List)
	group.Post("/", middleware.Protected(), r.ObjectTypeHandler.Create)
	group.Get("/:id", r.ObjectTypeHandler.Get)
	group.Put("/:id", middleware.Protected(), r.ObjectTypeHandler.Update)
	group.Delete("/:id", middleware.Protected(), r.ObjectTypeHandler.Delete)
}

func (r *RouteConfig) Organizer() {
	group := r.App.Group("/api/organizers")
	group.Get("/", r.OrganizerHandler.List)
	group.Post("/", middleware.Protected(), r.OrganizerHandler.Create)
	group.Get("/:id", r.OrganizerHandler.Get)
	group.Put("/:id", middleware.Protected(), r.OrganizerHandler.Update)
	group.Delete("/:id", middleware.Protected(), r.OrganizerHandler.Delete)
}

func (r *RouteConfig) Item() {
	group := r.App.Group("/api/items")
	group.Get("/", r.ItemHandler.List)
	group.Post("/", middleware.Protected(), r.ItemHandler.Create)
	group.Get("/:id", r.ItemHandler.Get)
	group.Put("/:id", middleware.Protected(), r.ItemHandler.Update)
	group.Delete("/:id", middleware.Protected(), r.ItemHandler.Delete)

	detail := group.Group("/:id/detail")
	detail.Get("/", r.ItemDetailHandler.Get)
	detail.Post("/", middleware.Protected(), r.ItemDetailHandler.Create)
	detail.Put("/", middleware.Protected(), r.ItemDetailHandler.Update)
	detail.Delete("/", middleware.Protected(), r.ItemDetailHandler.Delete)

	document := group.Group("/:id/document")
	document.Get("/", r.ItemDocumentHandler.Get)
	document.Post("/", middleware.Protected(), r.ItemDocumentHandler.Create)
	document.Put("/", middleware.Protected(), r.ItemDocumentHandler.Update)
	document.Delete("/", middleware.Protected(), r.ItemDocumentHandler.Delete)

	grade := group.Group("/:id/grade")
	grade.Get("/", r.ItemGradeHandler.Get)
	grade.Post("/", middleware.Protected(), r.ItemGradeHandler.Create)
	grade.Put("/", middleware.Protected(), r.ItemGradeHandler.Update)
	grade.Delete("/", middleware.Protected(), r.ItemGradeHandler.Delete)

	thumbnails := group.Group("/:id/thumbnails")
	thumbnails.Get("/", r.ItemThumbnailHandler.Get)
	thumbnails.Post("/", middleware.Protected(), r.ItemThumbnailHandler.Create)
	thumbnails.Put("/", middleware.Protected(), r.ItemThumbnailHandler.Update)
	thumbnails.Delete("/", middleware.Protected(), r.ItemThumbnailHandler.Delete)
}

func (r *RouteConfig) PIC() {
	group := r.App.Group("/api/pics")
	group.Get("/", r.PIChandler.List)
	group.Post("/", middleware.Protected(), r.PIChandler.Create)
	group.Get("/:id", r.PIChandler.Get)
	group.Put("/:id", middleware.Protected(), r.PIChandler.Update)
	group.Delete("/:id", middleware.Protected(), r.PIChandler.Delete)
}

func (r *RouteConfig) Auction() {
	group := r.App.Group("/api/auctions")
	group.Get("/", r.AuctionHandler.List)
	group.Post("/", middleware.Protected(), r.AuctionHandler.Create)
	group.Get("/:id", r.AuctionHandler.Get)
	group.Put("/:id", middleware.Protected(), r.AuctionHandler.Update)
	group.Delete("/:id", middleware.Protected(), r.AuctionHandler.Delete)
}

func (r *RouteConfig) User() {
	userGroup := r.App.Group("/api/users", middleware.Protected("admin"))
	userGroup.Get("/", r.UserHandler.List)
	userGroup.Post("/", r.UserHandler.Create)
	userGroup.Get("/:id", r.UserHandler.Get)
	userGroup.Put("/:id", r.UserHandler.Update)
	userGroup.Delete("/:id", r.UserHandler.Delete)
}

func (r *RouteConfig) Auth() {
	group := r.App.Group("/api/auth")
	group.Post("/login", r.UserHandler.Login)
	group.Get("/me", r.UserHandler.Me)
}
