package pagination

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type RequestMetaData struct {
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	SortBy  string `json:"sort_by"`
	SortDir string `json:"sort_dir"`
}

func NewRequestMetaData(c *fiber.Ctx) RequestMetaData {
	meta := RequestMetaData{
		Page:    1,
		Size:    10,
		SortBy:  "id",
		SortDir: "asc",
	}

	meta.Page = ToInt(c.Query("page"))
	meta.Size = DefaultTake(ToInt(c.Query("size")))
	sortBy := c.Query("sortBy")
	sortDir := c.Query("sortDir")

	if sortBy != "" {
		meta.SortBy = sortBy
	}

	if sortDir != "" {
		meta.SortDir = sortDir
	}

	return meta
}

func DefaultTake(i int) int {
	if i <= 0 {
		return 10
	}
	return i
}

func ToInt(i string) int {
	res, err := strconv.Atoi(i)
	if err != nil {
		return 0
	}
	return res
}
