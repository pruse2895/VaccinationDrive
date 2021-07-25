package appointment

import (
	"errors"
	"time"

	"github.com/FenixAra/go-util/log"

	"vaccinationDrive/internals/daos"
	"vaccinationDrive/models"
	"vaccinationDrive/utils"

	"github.com/go-pg/pg"
)

type AppointmentData struct {
	dbConn         *pg.DB
	l              *log.Logger
	AppointmentDao daos.AppointmentDao
}

func NewAppointmentData(l *log.Logger, dbConn *pg.DB) *AppointmentData {
	return &AppointmentData{
		l:              l,
		dbConn:         dbConn,
		AppointmentDao: daos.NewAppointmentData(l, dbConn),
	}

}

func (a *AppointmentData) BookAppointment(app models.Appointment) error {

	today := time.Now()

	bookingDate, _ := utils.StringddMMyyyyToDate(app.Date)

	diff := today.Sub(bookingDate)

	numofdays := int(diff.Hours() / 24)

	if numofdays > 90 {
		a.l.Errorf("You cannot book before 90 days")
		return errors.New("You cannot book before 90 days")
	}

	timeAvailableFlag := a.AppointmentDao.CheckTimeSlotAvailable(app)

	if !timeAvailableFlag {
		a.l.Errorf("Slots are booked for selected time")
		return errors.New("Slots are booked for selected time")
	}

	totalVacineAvailableFlag := a.AppointmentDao.CheckTotalVacineAvailable(app)
	doseAvailableFlag := a.AppointmentDao.CheckDoseAvailable(app)

	if !totalVacineAvailableFlag || !doseAvailableFlag {
		a.l.Errorf("Vaccine are not available for selected day")
		return errors.New("Vaccine are not available for selected day")
	}

	maxSlot := a.AppointmentDao.CheckSlotsBooked(app)
	if !maxSlot {
		a.l.Errorf("you are reached the maximum slots")
		return errors.New("you are reached the maximum slots")
	}

	beneficiary, err := a.AppointmentDao.CheckDaysBetweenDoses(app)
	if err == nil {
		benDate, _ := utils.StringddMMyyyyToDate(beneficiary.Date)

		doseDateDiff := bookingDate.Sub(benDate)

		if doseDateDiff < 15 {
			a.l.Errorf("Book the slot after 15 days")
			return errors.New("Book the slot after 15 days")
		}
	}

	if beneficiary.ID > 0 {
		app.UpdatedAt = time.Now()
		uperr := a.AppointmentDao.UpdateAppointment(app)
		if uperr != nil {
			a.l.Errorf("BookAppointment Error : ", uperr)
			return uperr
		}

		return nil
	}

	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()

	err = a.AppointmentDao.SaveAppointment(app)
	if err != nil {
		a.l.Errorf("BookAppointment Error : ", err)
		return err
	}

	return nil

}
