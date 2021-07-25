package daos

import (
	"github.com/FenixAra/go-util/log"

	"vaccinationDrive/models"

	"github.com/go-pg/pg"
)

type AppointmentObj struct {
	l      *log.Logger
	dbConn *pg.DB
}

func NewAppointmentData(l *log.Logger, dbConn *pg.DB) *AppointmentObj {
	return &AppointmentObj{
		l:      l,
		dbConn: dbConn,
	}
}

type AppointmentDao interface {
	SaveAppointment(Appointment models.Appointment) error
	CheckTimeSlotAvailable(Appointment models.Appointment) bool
	CheckTotalVacineAvailable(Appointment models.Appointment) bool
	CheckDoseAvailable(Appointment models.Appointment) bool
	CheckDaysBetweenDoses(Appointment models.Appointment) (*models.Appointment, error)
	CheckSlotsBooked(Appointment models.Appointment) bool
	UpdateAppointment(Appointment models.Appointment) error
}

func (a *AppointmentObj) SaveAppointment(Appointment models.Appointment) error {

	err := a.dbConn.Insert(&Appointment)
	if err != nil {
		a.l.Errorf("SaveAppointment Error ", err)
		return err
	}
	return nil
}

func (a *AppointmentObj) UpdateAppointment(Appointment models.Appointment) error {
	_, err := a.dbConn.Model(&Appointment).Column("date", "time_slot", "dose", "vaccine_center", "updated_at").Where("id = ? ", Appointment.ID).Returning("*").Update()
	if err != nil {
		a.l.Errorf("UpdateAppointment Error ", err)
		return err
	}
	return nil
}

func (a *AppointmentObj) CheckTimeSlotAvailable(Appointment models.Appointment) bool {

	c, _ := a.dbConn.Model(&Appointment).Where("date = ? AND time_slot = ? AND vaccine_center = ? ", Appointment.Date, Appointment.TimeSlot, Appointment.VaccineCenter).Count()
	if c > 10 {
		return false
	}
	return true
}

func (a *AppointmentObj) CheckTotalVacineAvailable(Appointment models.Appointment) bool {
	c, _ := a.dbConn.Model(&Appointment).Where("date = ? AND vaccine_center = ? ", Appointment.Date, Appointment.VaccineCenter).Count()
	if c > 30 {
		return false
	}
	return true
}

func (a *AppointmentObj) CheckDoseAvailable(Appointment models.Appointment) bool {
	c, _ := a.dbConn.Model(&Appointment).Where("date = ? AND vaccine_center = ? AND dose = ?", Appointment.Date, Appointment.VaccineCenter, Appointment.Dose).Count()
	if c > 15 {
		return false
	}
	return true
}

func (a *AppointmentObj) CheckDaysBetweenDoses(Appointment models.Appointment) (*models.Appointment, error) {
	err := a.dbConn.Model(&Appointment).Where("beneficiary_id = ?", Appointment.BeneficiaryID).Select()
	if err != nil {
		return nil, err
	}
	return &Appointment, nil
}

func (a *AppointmentObj) CheckSlotsBooked(Appointment models.Appointment) bool {
	c, _ := a.dbConn.Model(&Appointment).Where("beneficiary_id = ?", Appointment.BeneficiaryID).Count()
	if c > 2 {
		return false
	}
	return true
}
