package domain

type UpdateRequestOrganizer struct {
	Name          *string `json:"name"`
	Address       *string `json:"address"`
	BankName      *string `json:"bank_name"`
	AccountNumber *string `json:"account_number"`
	AccountName   *string `json:"account_name"`
}
