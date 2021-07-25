package routes

import (
	"net/http"
	app "vaccinationDrive/internals/services/appointment"
	"vaccinationDrive/models"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func appointment(router *httprouter.Router, indexHandlers alice.Chain) {
	router.POST("/bookappointment", wrapHandler(indexHandlers.ThenFunc(BookAppointment)))
	router.PUT("/updateappointment/:id", wrapHandler(indexHandlers.ThenFunc(UpdateAppointment)))
}

func BookAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentIns := models.Appointment{}

	rd := logAndGetContext(w, r)

	if !parseJSON(w, r.Body, &appointmentIns) {
		return
	}

	appIns := app.NewAppointmentData(rd.l, rd.dbConn)
	err := appIns.BookAppointment(appointmentIns)
	if err != nil {
		rd.l.Errorf("BookAppointment - ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	res := struct {
		Message string `json:"message"`
	}{
		"Your appointment has been Booked successfully..",
	}

	writeJSONStruct(res, http.StatusOK, rd)

}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {

	// ID, isErr := GetIDFromParams(w, r, "id")
	// if !isErr {
	// 	return
	// }

	appointmentIns := models.Appointment{}

	rd := logAndGetContext(w, r)

	if !parseJSON(w, r.Body, &appointmentIns) {
		return
	}

	appIns := app.NewAppointmentData(rd.l, rd.dbConn)
	err := appIns.BookAppointment(appointmentIns)
	if err != nil {
		rd.l.Errorf("BookAppointment - ", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	res := struct {
		Message string `json:"message"`
	}{
		"Your appointment has been Booked successfully..",
	}

	writeJSONStruct(res, http.StatusOK, rd)

}
