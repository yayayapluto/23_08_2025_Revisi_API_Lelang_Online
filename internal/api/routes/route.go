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
	group.Post("/", middleware.Protected("admin"), r.ObjectTypeHandler.Create)
	group.Get("/:id", r.ObjectTypeHandler.Get)
	group.Put("/:id", middleware.Protected("admin"), r.ObjectTypeHandler.Update)
	group.Delete("/:id", middleware.Protected("admin"), r.ObjectTypeHandler.Delete)
}

func (r *RouteConfig) Organizer() {
	group := r.App.Group("/api/organizers")
	group.Get("/", r.OrganizerHandler.List)
	group.Post("/", middleware.Protected("admin"), r.OrganizerHandler.Create)
	group.Get("/:id", r.OrganizerHandler.Get)
	group.Put("/:id", middleware.Protected("admin"), r.OrganizerHandler.Update)
	group.Delete("/:id", middleware.Protected("admin"), r.OrganizerHandler.Delete)
}

func (r *RouteConfig) Item() {
	group := r.App.Group("/api/items")
	group.Get("/", r.ItemHandler.List)
	group.Post("/", middleware.Protected("admin", "organizer"), r.ItemHandler.Create)
	group.Get("/:id", r.ItemHandler.Get)
	group.Put("/:id", middleware.Protected("admin", "organizer"), r.ItemHandler.Update)
	group.Delete("/:id", middleware.Protected("admin", "organizer"), r.ItemHandler.Delete)

	detail := group.Group("/:id/detail")
	detail.Get("/", r.ItemDetailHandler.Get)
	detail.Post("/", middleware.Protected("admin", "organizer"), r.ItemDetailHandler.Create)
	detail.Put("/", middleware.Protected("admin", "organizer"), r.ItemDetailHandler.Update)
	detail.Delete("/", middleware.Protected("admin", "organizer"), r.ItemDetailHandler.Delete)

	document := group.Group("/:id/document")
	document.Get("/", r.ItemDocumentHandler.Get)
	document.Post("/", middleware.Protected("admin", "organizer"), r.ItemDocumentHandler.Create)
	document.Put("/", middleware.Protected("admin", "organizer"), r.ItemDocumentHandler.Update)
	document.Delete("/", middleware.Protected("admin", "organizer"), r.ItemDocumentHandler.Delete)

	grade := group.Group("/:id/grade")
	grade.Get("/", r.ItemGradeHandler.Get)
	grade.Post("/", middleware.Protected("admin", "organizer"), r.ItemGradeHandler.Create)
	grade.Put("/", middleware.Protected("admin", "organizer"), r.ItemGradeHandler.Update)
	grade.Delete("/", middleware.Protected("admin", "organizer"), r.ItemGradeHandler.Delete)

	thumbnails := group.Group("/:id/thumbnails")
	thumbnails.Get("/", middleware.Protected("admin", "organizer"), r.ItemThumbnailHandler.Get)
	thumbnails.Post("/", middleware.Protected("admin", "organizer"), r.ItemThumbnailHandler.Create)
	thumbnails.Put("/", middleware.Protected("admin", "organizer"), r.ItemThumbnailHandler.Update)
	thumbnails.Delete("/", middleware.Protected("admin", "organizer"), r.ItemThumbnailHandler.Delete)
}

func (r *RouteConfig) PIC() {
	group := r.App.Group("/api/pics")
	group.Get("/", r.PIChandler.List)
	group.Post("/", middleware.Protected("admin", "organizer"), r.PIChandler.Create)
	group.Get("/:id", r.PIChandler.Get)
	group.Put("/:id", middleware.Protected("admin", "organizer"), r.PIChandler.Update)
	group.Delete("/:id", middleware.Protected("admin", "organizer"), r.PIChandler.Delete)
}

func (r *RouteConfig) Auction() {
	group := r.App.Group("/api/auctions")
	group.Get("/", r.AuctionHandler.List)
	group.Post("/", middleware.Protected("admin", "organizer"), r.AuctionHandler.Create)
	group.Get("/:id", r.AuctionHandler.Get)
	group.Put("/:id", middleware.Protected("admin", "organizer"), r.AuctionHandler.Update)
	group.Delete("/:id", middleware.Protected("admin", "organizer"), r.AuctionHandler.Delete)
}

func (r *RouteConfig) User() {
	userGroup := r.App.Group("/api/users", middleware.Protected("admin"))
	userGroup.Get("/", r.UserHandler.List)
	userGroup.Post("/", r.UserHandler.Create)
	userGroup.Post("/organizer", r.UserHandler.CreateOrganizerUser)
	userGroup.Get("/:id", r.UserHandler.Get)
	userGroup.Put("/:id", r.UserHandler.Update)
	userGroup.Delete("/:id", r.UserHandler.Delete)
}

func (r *RouteConfig) Auth() {
	group := r.App.Group("/api/auth")
	group.Post("/login", r.UserHandler.Login)
	group.Post("/register", r.UserHandler.Register)
	group.Get("/me", middleware.Protected(), r.UserHandler.Me)
}
