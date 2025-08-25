package utils

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"net/url"
)

func BuildFileURL(ctx *fiber.Ctx, file *entities.File) (*string, error) {
	baseUrl := ctx.Protocol() + "://" + ctx.Hostname()
	newUrlPath, err := url.JoinPath(baseUrl, file.Path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to build file url: %d\n", err))
	}
	return &newUrlPath, nil
}
