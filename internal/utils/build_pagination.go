package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils/pagination"
	"math"
)

func BuildPagination[T any](ctx *fiber.Ctx, rm RequestMeta, data []T, total int64) pagination.ResponseMetaData[T] {
	currentPageUrl := pagination.BuildPageURL(ctx, rm.Search, rm.Page, rm.Size, rm.SortBy, rm.SortDir)
	firstPageUrl := pagination.BuildPageURL(ctx, rm.Search, 1, rm.Size, rm.SortBy, rm.SortDir)

	var nextPageUrl, prevPageUrl *string
	if rm.Offset+len(data) > int(total) {
		url := pagination.BuildPageURL(ctx, rm.Search, rm.Page+1, rm.Size, rm.SortBy, rm.SortDir)
		nextPageUrl = &url
	}

	if rm.Page > 1 {
		url := pagination.BuildPageURL(ctx, rm.Search, rm.Page-1, rm.Size, rm.SortBy, rm.SortDir)
		prevPageUrl = &url
	}

	totalPage := int(math.Floor(float64(total) / float64(rm.Size)))
	paginationRes := pagination.NewResponseMetaData[T](rm.Page, currentPageUrl, data, firstPageUrl, nextPageUrl, rm.Size, prevPageUrl, totalPage)
	return paginationRes
}
