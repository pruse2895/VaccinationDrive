package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	lg "log"
	"net/http"
	"strconv"
	"time"
	"vaccinationDrive/dbcon"
	"vaccinationDrive/models"
	validator "vaccinationDrive/validators"

	"github.com/FenixAra/go-util/log"
	"github.com/go-pg/pg"
	"github.com/julienschmidt/httprouter"
)

const (
	ERR_MSG = "ERROR_MESSAGE"
	MSG     = "MESSAGE"
)

type ResStruct struct {
	Status   string `json:"status" example:"SUCCESS" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"200" example:"500"`
	Message  string `json:"message" example:"pong" example:"could not connect to db"`
}

type Res500Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"500"`
	Message  string `json:"message" example:"could not connect to db"`
}

type Res400Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"400"`
	Message  string `json:"message" example:"Invalid param"`
}

type RequestData struct {
	l      *log.Logger
	dbConn *pg.DB
	Start  time.Time
	w      http.ResponseWriter
	r      *http.Request
}

type RenderData struct {
	Data  interface{}
	Paths []string
}

type TemplateData struct {
	Data interface{}
}

func (t *TemplateData) SetConstants() {

}

func logAndGetContext(w http.ResponseWriter, r *http.Request) *RequestData {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-Frame-Options", "DENY")
	//Set config according to the use case..
	cfg := log.NewConfig("")
	//cfg.SetRemoteConfig(conf.LOG_REMOTE_URL, "", "")
	cfg.SetLevelStr("")
	cfg.SetFilePathSizeStr("")
	cfg.SetReference(r.Header.Get("ReferenceID"))
	l := log.New(cfg)
	db := dbcon.Get()
	//dbConn := new(db.DBConn)
	//dbConn.Init(l)
	//pgdbConn := new(pgsqldb.Conn)
	//pgdbConn.Init(l)
	start := time.Now()
	//l.LogAPIInfo(r, 0, 0)
	return &RequestData{
		l:      l,
		dbConn: db,
		Start:  start,
		r:      r,
		w:      w,
	}
}

func jsonifyMessage(msg string, msgType string, httpCode int) ([]byte, int) {
	var data []byte
	var Obj struct {
		Status   string `json:"status"`
		HTTPCode int    `json:"code"`
		Message  string `json:"message"`
		Err      error  `json:"error"`
	}
	Obj.Message = msg
	Obj.HTTPCode = httpCode
	switch msgType {
	case ERR_MSG:
		Obj.Status = "FAILED"

	case MSG:
		Obj.Status = "SUCCESS"
	}
	data, _ = json.Marshal(Obj)
	return data, httpCode
}

func writeJSONMessage(msg string, msgType string, httpCode int, rd *RequestData) {
	d, code := jsonifyMessage(msg, msgType, httpCode)
	writeJSONResponse(d, code, rd)
}

func writeJSONStruct(v interface{}, code int, rd *RequestData) {
	d, err := json.Marshal(v)
	if err != nil {
		writeJSONMessage("Unable to marshal data. Err: "+err.Error(), ERR_MSG, http.StatusInternalServerError, rd)
		return
	}
	writeJSONResponse(d, code, rd)
}

func writeJSONResponse(d []byte, code int, rd *RequestData) {
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), code)
	if code == http.StatusInternalServerError {
		rd.l.Info(rd.r.URL.String(), "Status Code:", code, ", Response time:", time.Since(rd.Start), " Response:", string(d))
	} else {
		rd.l.Info(rd.r.URL.String(), "Status Code:", code, ", Response time:", time.Since(rd.Start))
	}
	rd.w.Header().Set("Access-Control-Allow-Origin", "*")
	rd.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rd.w.WriteHeader(code)
	rd.w.Write(d)
}

func writeJSONMessageWithData(msg string, msgType string, httpCode int, rd *RequestData, functionName string, requestData string) {
	d, code := jsonifyMessage(msg, msgType, httpCode)
	writeJSONResponseWithData(d, code, rd, functionName, requestData)
}

func writeJSONResponseWithData(d []byte, code int, rd *RequestData, functionName string, requestData string) {
	rd.l.LogAPIInfo(rd.r, time.Since(rd.Start).Seconds(), code)

	rd.l.Info("Service name : ", functionName, ", Request Data : ", requestData,
		", Response data : ", string(d), ", Status Code : ", code, ", Response time : ", time.Since(rd.Start))

	rd.w.Header().Set("Access-Control-Allow-Origin", "*")
	rd.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rd.w.WriteHeader(code)
	rd.w.Write(d)
}

func writeJSONStructWithData(v interface{}, code int, rd *RequestData, functionName string, requestData string) {
	d, err := json.Marshal(v)
	if err != nil {
		writeJSONMessageWithData("Unable to marshal data. Err: "+err.Error(), ERR_MSG, http.StatusInternalServerError, rd, functionName, requestData)
		return
	}
	writeJSONResponseWithData(d, code, rd, functionName, requestData)
}

func renderJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		lg.Printf("ERROR: renderJson - %q\n", err)
	}
}

//renderValidationError to render error msg validation
func renderValidationError(w http.ResponseWriter, status int, errs validator.Errors) {

	res := struct {
		Message string           `json:"message"`
		Errors  validator.Errors `json:"errors"`
	}{"Validation Error(s)", errs}
	renderJSON(w, status, res)
}

func parseJSON(w http.ResponseWriter, body io.ReadCloser, model interface{}) bool {
	defer body.Close()

	b, _ := ioutil.ReadAll(body)
	err := json.Unmarshal(b, model)

	if err != nil {
		e := &models.ErrorData{}
		e.Message = "Error in parsing json"
		e.Err = err
		renderERRORWithIn(w, e)
		return false
	}

	return true
}

func renderERRORWithIn(w http.ResponseWriter, err *models.ErrorData) {
	err.Set()
	renderJSON(w, err.Code, err)
}

func GetIDFromParams(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	params, _ := r.Context().Value("params").(httprouter.Params)
	idStr := params.ByName(key)
	id, err := strconv.ParseInt(idStr, 10, 64)
	isErr := true

	if err != nil {
		isErr = false
		e := &models.ErrorData{}
		e.Message = "Invalid ID"
		e.Code = http.StatusBadRequest
		e.Err = err
		renderERRORWithIn(w, e)
	}

	return id, isErr
}
