package domain

import "time"

type UpdateRequestItemDetail struct {
	PlateNumber   *string    `json:"plate_number,omitempty"`
	Brand         *string    `json:"brand,omitempty"`
	Series        *string    `json:"series,omitempty"`
	CC            *float64   `json:"cc,omitempty"`
	Type          *string    `json:"type,omitempty"`
	Transmission  *string    `json:"transmission,omitempty"`
	Model         *string    `json:"model,omitempty"`
	Year          *int       `json:"year,omitempty"`
	FrameNumber   *string    `json:"frame_number,omitempty"`
	MachineNumber *string    `json:"machine_number,omitempty"`
	Kilometer     *int       `json:"kilometer,omitempty"`
	Fuel          *string    `json:"fuel,omitempty"`
	Color         *string    `json:"color,omitempty"`
	DriveType     *string    `json:"drive_type,omitempty"`
	StnkDate      *time.Time `json:"stnk_date,omitempty"`
}
