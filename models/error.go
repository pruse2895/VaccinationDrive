package models

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-pg/pg"
)

// ErrorData struct
type ErrorData struct {
	Message string `json:"message"`
	Err     error  `json:"error"`
	IsDbErr bool   `json:"-"`
	Code    int    `json:"code"`
}

func (e ErrorData) Error() string {
	e.Set()
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

//Set func set error in model
func (e *ErrorData) Set() {

	if e.IsDbErr {
		e.SetDbError()
		return
	}

	if e.Code == 0 {
		e.Code = http.StatusBadRequest
	}

	if e.Message != "" {
		return
	}

	switch e.Code {
	case 404:
		e.Message = fmt.Sprintf("%s not found", e)
	default:
		e.Message = fmt.Sprintf("%d: Error", e.Code)
	}
}

//SetDbError set db error in model
func (e *ErrorData) SetDbError() {

	switch e.Err {
	case pg.ErrNoRows:
		e.Code = http.StatusNotFound
		e.Message = fmt.Sprintf("Data not found")
		return
	}

	dbErr, ok := (e.Err).(pg.Error)

	if !ok {
		e.Code = http.StatusBadRequest
		if e.Message == "" {
			e.Message = fmt.Sprintf("DB Error: %s", e.Stack())
		}
		return
	}

	switch dbErr.Field('C') {
	case "23505":
		e.Code = http.StatusBadRequest
		e.Message = dbErr.Field('D')

		//get unqiue field
		key := strings.Split(e.Message, "=")
		keyIndex := strings.Index(key[0], " ")
		str := key[0]
		// log.Printf("error", str[keyIndex: len(str)-1])
		uniqueField := str[keyIndex : len(str)-1]
		uniqueFieldName := uniqueField + ") already exists"
		if uniqueFieldName != "" {
			e.Message = uniqueFieldName
		}

	default:
		e.Code = http.StatusBadRequest
		if e.Message == "" {
			e.Message = fmt.Sprintf("DB Error: %s", e.Stack())
		}
	}
}

// Stack func
func (e ErrorData) Stack() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

//MarshalJSON implements json marshaller
// func (e ErrorData) MarshalJSON() ([]byte, error) {
// 	type Alias ErrorData
// 	return json.Marshal(&struct {
// 		Err string `json:"error"`
// 		Alias
// 	}{
// 		Err:   e.Stack(),
// 		Alias: (Alias)(e),
// 	})
// }
