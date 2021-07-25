package models

import "time"

type Appointment struct {
	ID            int64     `json:"id"`
	BeneficiaryID int64     `json:"beneficiarId"`
	Date          string    `json:"date"`
	TimeSlot      string    `json:"timeSlot"`
	Dose          string    `json:"dose"`
	VaccineCenter string    `json:"vaccineCenter"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
