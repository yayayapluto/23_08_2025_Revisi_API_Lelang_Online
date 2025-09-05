package pagination

type ResponseMetaData[T any] struct {
	CurrentPage    int     `json:"current_page"`
	CurrentPageUrl string  `json:"current_page_url"`
	Data           []T     `json:"data"`
	FirstPageUrl   string  `json:"first_page_url"`
	NextPageUrl    *string `json:"next_page_url"`
	PerPage        int     `json:"per_page"`
	PrevPageUrl    *string `json:"prev_page_url"`
	TotalPages     int     `json:"total_pages"`
}

func NewResponseMetaData[T any](currentPage int, currentPageUrl string, data []T, firstPageUrl string, nextPageUrl *string, perPage int, prevPageUrl *string, totalPages int) ResponseMetaData[T] {
	return ResponseMetaData[T]{
		CurrentPage:    currentPage,
		CurrentPageUrl: currentPageUrl,
		Data:           data,
		FirstPageUrl:   firstPageUrl,
		NextPageUrl:    nextPageUrl,
		PerPage:        perPage,
		PrevPageUrl:    prevPageUrl,
		TotalPages:     totalPages,
	}
}
