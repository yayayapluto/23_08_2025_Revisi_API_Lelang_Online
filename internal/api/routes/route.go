package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
)

type RouteConfig struct {
	App                 *fiber.App
	ObjectTypeHandler   handlers.ObjectTypeHandler
	OrganizerHandler    handlers.OrganizerHandler
	ItemHandler         handlers.ItemHandler
	ItemDetailHandler   handlers.ItemDetailHandler
	ItemDocumentHandler handlers.ItemDocumentHandler
	ItemGradeHandler    handlers.ItemGradeHandler
}

func (r *RouteConfig) Setup() {
	r.ObjectType()
	r.Organizer()
	r.Item()
	r.ItemDetail()
	r.ItemDocument()
	r.ItemGrade()
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
}

func (r *RouteConfig) ItemDetail() {
	group := r.App.Group("/api/itemDetails")
	group.Get("/", r.ItemDetailHandler.List)
	group.Post("/", r.ItemDetailHandler.Create)
	group.Get("/:itemID", r.ItemDetailHandler.Get)
	group.Put("/:itemID", r.ItemDetailHandler.Update)
	group.Delete("/:itemID", r.ItemDetailHandler.Delete)
}

func (r *RouteConfig) ItemDocument() {
	group := r.App.Group("/api/itemDocuments")
	group.Get("/", r.ItemDocumentHandler.List)
	group.Post("/", r.ItemDocumentHandler.Create)
	group.Get("/:itemID", r.ItemDocumentHandler.Get)
	group.Put("/:itemID", r.ItemDocumentHandler.Update)
	group.Delete("/:itemID", r.ItemDocumentHandler.Delete)
}

func (r *RouteConfig) ItemGrade() {
	group := r.App.Group("/api/itemGrades")
	group.Get("/", r.ItemGradeHandler.List)
	group.Post("/", r.ItemGradeHandler.Create)
	group.Get("/:itemID", r.ItemGradeHandler.Get)
	group.Put("/:itemID", r.ItemGradeHandler.Update)
	group.Delete("/:itemID", r.ItemGradeHandler.Delete)
}
