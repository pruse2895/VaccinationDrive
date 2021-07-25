package models

import (
	"errors"
	"time"
	"vaccinationDrive/dbcon"
	validator "vaccinationDrive/validators"
)

const (
	UserMobilePattern = `^((\+)?(\d{2}[-])?(\d{10}){1})?(\d{11}){0,1}?$`
)

type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	DOB         string    `json:"dob"`
	Age         float64   `json:"age"`
	AadharNo    string    `json:"aadharNo" validate:"required" sql:",notnull"`
	PhoneNumber string    `json:"phoneNumber" validate:"required" sql:",notnull"`
	CreatedAt   time.Time `json:"-" sql:",default:now()"`
	UpdatedAt   time.Time `json:"-" sql:",default:now()"`
}

//Validate is validation for User fields
func (us User) Validate() (validator.Errors, error) {
	db := dbcon.Get()
	v := validator.New("User")

	count, _ := db.Model(&us).Where("aadhar_no = ?", us.AadharNo).Count()
	if count > 0 {
		v.AddError("aadharNo", errors.New("AadharNo is already present"))

	}

	//aadhar no length
	if us.AadharNo != "" {
		if aadharNumberLength := len(us.AadharNo); aadharNumberLength != 15 {
			v.AddError("aadharNo", errors.New("AadharNo should be length of 15 digits"))
		}

	}

	//mobile no length
	if us.PhoneNumber != "" {
		if phoneNumberLength := len(us.PhoneNumber); phoneNumberLength != 10 {
			v.AddError("PhoneNumber", errors.New("Phone Number should be length of 10 digits"))
		}

		//valid mobile no
		v.ValidateField("PhoneNumber", us.PhoneNumber, []validator.Tag{
			{Name: "regexp", Fn: validator.Regex, Param: UserMobilePattern},
		})
	}

	return v.Validate(us)
}

// BeforeInsert func
func (us *User) BeforeInsert() {
	us.CreatedAt = time.Now().UTC()
	us.UpdatedAt = time.Now().UTC()
}
