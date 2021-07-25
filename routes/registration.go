package routes

import (
	"math"
	"net/http"
	"time"
	"vaccinationDrive/internals/services/userRegistration"
	"vaccinationDrive/models"
	"vaccinationDrive/utils"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func registration(router *httprouter.Router, indexHandlers alice.Chain) {
	router.POST("/user", wrapHandler(indexHandlers.ThenFunc(RegisterUser)))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	rd := logAndGetContext(w, r)

	user := models.User{}

	if !parseJSON(w, r.Body, &user) {
		return
	}

	if errs, err := user.Validate(); err != nil {
		rd.l.Error("Errors : ", errs)
		renderValidationError(w, http.StatusBadRequest, errs)
		return
	}

	today := time.Now()
	dateofbirth, _ := utils.StringddMMyyyyToDate(user.DOB)

	user.Age = math.Floor(today.Sub(dateofbirth).Hours() / 24 / 365)

	if user.Age < 45 {
		rd.l.Errorf("user age should be less than 45")
		writeJSONMessage("user age should be less than 45", ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	userIns := userRegistration.NewUserData(rd.l, rd.dbConn)
	err := userIns.RegisterUser(user)
	if err != nil {
		rd.l.Errorf("error in insert user", err.Error())
		writeJSONMessage(err.Error(), ERR_MSG, http.StatusBadRequest, rd)
		return
	}

	res := struct {
		Message string `json:"message"`
	}{
		"User Registered Successfully...",
	}

	writeJSONStruct(res, http.StatusOK, rd)

}
