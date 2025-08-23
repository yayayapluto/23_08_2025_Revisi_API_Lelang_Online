package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/handlers"
)

type RouteConfig struct {
	App               *fiber.App
	ObjectTypeHandler handlers.ObjectTypeHandler
	OrganizerHandler  handlers.OrganizerHandler
}

func (r *RouteConfig) Setup() {
	r.ObjectType()
	r.Organizer()
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
