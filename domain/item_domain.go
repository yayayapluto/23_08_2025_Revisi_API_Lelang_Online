package domain

type UpdateRequestItem struct {
	ObjectTypeID *uint    `json:"object_type_id,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	DepositPrice *float64 `json:"deposit_price,omitempty"`
	Description  *string  `json:"description,omitempty"`
	FileID       *uint    `json:"file_id,omitempty"`
}
