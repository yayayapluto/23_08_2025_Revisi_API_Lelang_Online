package domain

type UpdateRequestItemDocument struct {
	Bpkb             *bool `json:"bpkb,omitempty"`
	Stnk             *bool `json:"stnk,omitempty"`
	Facture          *bool `json:"facture,omitempty"`
	Receipt          *bool `json:"receipt,omitempty"`
	OwnershipRelease *bool `json:"ownership_release,omitempty"`
	Warranty         *bool `json:"warranty,omitempty"`
	Box              *bool `json:"box,omitempty"`
}
