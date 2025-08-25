package utils

import (
	"github.com/gofiber/fiber/v2"
	"math"
)

type RequestMeta struct {
	Search  string
	Page    int
	Size    int
	Offset  int
	SortBy  string
	SortDir string
}

func GetRequestMeta(ctx *fiber.Ctx) RequestMeta {
	var rm RequestMeta

	size := ctx.QueryInt("size")
	size = int(math.Min(math.Max(float64(size), 10), 100)) // min 10, max 100

	rm.Search = ctx.Query("search")
	rm.Page = ctx.QueryInt("page", 1)
	rm.Size = size
	rm.Offset = (rm.Page - 1) * rm.Size
	rm.SortBy = ctx.Query("sortBy", "id")
	rm.SortDir = ctx.Query("sortDir", "asc")

	return rm
}
