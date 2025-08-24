package utils

import (
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"slices"
)

func BuildOrderQuery(validSortBy, validSortDir []string, sortBy, sortDir *string) (*string, error) {
	defSortBy := validSortBy[0]
	if sortBy != nil {
		if !slices.Contains(validSortBy, *sortBy) {
			return nil, domain.ErrInvalidSortByColumn
		}
		defSortBy = *sortBy
	}

	defSortDir := validSortDir[0]
	if sortDir != nil {
		if !slices.Contains(validSortDir, *sortDir) {
			return nil, domain.ErrInvalidSortDir
		}
		defSortDir = *sortDir
	}

	orderStr := fmt.Sprintf("%s %s", defSortBy, defSortDir)
	return &orderStr, nil
}
