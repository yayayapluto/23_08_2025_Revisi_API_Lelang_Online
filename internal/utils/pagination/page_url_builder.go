package pagination

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strconv"
	"strings"
)

func BuildPageURL(c *fiber.Ctx, search string, page, size int, sortBy, sortDir string) string {
	u, _ := url.Parse(c.BaseURL() + strings.TrimRight(c.Path(), "/"))

	q := u.Query()
	if search != "" {
		q.Set("search", search)
	}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if size > 0 {
		q.Set("size", strconv.Itoa(size))
	}
	if sortBy != "" {
		q.Set("sortBy", sortBy)
	}
	if sortDir != "" {
		q.Set("sortDir", sortDir)
	}

	u.RawQuery = q.Encode()
	return u.String()
}
