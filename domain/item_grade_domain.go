package domain

type UpdateRequestItemGrade struct {
	Interior *string `json:"interior,omitempty"`
	Exterior *string `json:"exterior,omitempty"`
	Frame    *string `json:"frame,omitempty"`
	Machine  *string `json:"machine,omitempty"`
}
