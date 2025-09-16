package domain

type (
	CreateRequestItem struct {
		ObjectTypeID uint    `json:"object_type_id"`
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		DepositPrice float64 `json:"deposit_price"`
		Description  string  `json:"description"`
		FileID       uint64  `json:"file_id"`
	}

	UpdateRequestItem struct {
		ObjectTypeID *uint    `json:"object_type_id,omitempty"`
		Name         *string  `json:"name,omitempty"`
		Price        *float64 `json:"price,omitempty"`
		DepositPrice *float64 `json:"deposit_price,omitempty"`
		Description  *string  `json:"description,omitempty"`
		FileID       *uint    `json:"file_id,omitempty"`
	}
)
