package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
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
}

func (r *RouteConfig) Setup() {
	r.ObjectType()
	r.Organizer()
	r.Item()
	r.PIC()
	r.Auction()
}

func (r *RouteConfig) ObjectType() {
	group := r.App.Group("/api/objectTypes")
	group.Get("/", r.ObjectTypeHandler.List)
	group.Post("/", r.ObjectTypeHandler.Create)
	group.Get("/:id", r.ObjectTypeHandler.Get)
	group.Put("/:id", r.ObjectTypeHandler.Update)
	group.Delete("/:id", r.ObjectTypeHandler.Delete)
}

func (r *RouteConfig) Organizer() {
	group := r.App.Group("/api/organizers")
	group.Get("/", r.OrganizerHandler.List)
	group.Post("/", r.OrganizerHandler.Create)
	group.Get("/:id", r.OrganizerHandler.Get)
	group.Put("/:id", r.OrganizerHandler.Update)
	group.Delete("/:id", r.OrganizerHandler.Delete)
}

func (r *RouteConfig) Item() {
	group := r.App.Group("/api/items")
	group.Get("/", r.ItemHandler.List)
	group.Post("/", r.ItemHandler.Create)
	group.Get("/:id", r.ItemHandler.Get)
	group.Put("/:id", r.ItemHandler.Update)
	group.Delete("/:id", r.ItemHandler.Delete)

	detail := group.Group("/:id/detail")
	detail.Get("/", r.ItemDetailHandler.Get)
	detail.Post("/", r.ItemDetailHandler.Create)
	detail.Put("/", r.ItemDetailHandler.Update)
	detail.Delete("/", r.ItemDetailHandler.Delete)

	document := group.Group("/:id/document")
	document.Get("/", r.ItemDocumentHandler.Get)
	document.Post("/", r.ItemDocumentHandler.Create)
	document.Put("/", r.ItemDocumentHandler.Update)
	document.Delete("/", r.ItemDocumentHandler.Delete)

	grade := group.Group("/:id/grade")
	grade.Get("/", r.ItemGradeHandler.Get)
	grade.Post("/", r.ItemGradeHandler.Create)
	grade.Put("/", r.ItemGradeHandler.Update)
	grade.Delete("/", r.ItemGradeHandler.Delete)

	thumbnails := group.Group("/:id/thumbnails")
	thumbnails.Get("/", r.ItemThumbnailHandler.Get)
	thumbnails.Post("/", r.ItemThumbnailHandler.Create)
	thumbnails.Put("/", r.ItemThumbnailHandler.Update)
	thumbnails.Delete("/", r.ItemThumbnailHandler.Delete)
}

func (r *RouteConfig) PIC() {
	group := r.App.Group("/api/pics")
	group.Get("/", r.PIChandler.List)
	group.Post("/", r.PIChandler.Create)
	group.Get("/:id", r.PIChandler.Get)
	group.Put("/:id", r.PIChandler.Update)
	group.Delete("/:id", r.PIChandler.Delete)
}

func (r *RouteConfig) Auction() {
	group := r.App.Group("/api/auctions")
	group.Get("/", r.AuctionHandler.List)
	group.Post("/", r.AuctionHandler.Create)
	group.Get("/:id", r.AuctionHandler.Get)
	group.Put("/:id", r.AuctionHandler.Update)
	group.Delete("/:id", r.AuctionHandler.Delete)
}
