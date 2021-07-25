package validator

import (
	"errors"
	"fmt"
	"strings"
)

func init() {
	setDefaultMessage()
}

//Messages implements the error messages
type Messages map[string]string

var errNotFound = errors.New("$err:message not found")

//Get returns the specified error message of the field
//if key not available it returns 'not foud error'
func (e Messages) Get(key string, args ...string) error {
	if v, ok := e[key]; ok {
		return errors.New(formatMessage(v, args))
	}

	return errNotFound
}

func formatMessage(message string, args []string) string {
	for i, v := range args {
		message = strings.Replace(message, fmt.Sprintf("${%d}", i), v, 1)
	}
	return message
}

//defaultMessages gives the default error messages
// ${index} - is used for formating the message with dynamic values
//Eg. `${0} ${1} cannot be blank` - `'Category' 'name' cannot be blank`
// ${0} - Category which module & ${1} - name which is field
// to omit module name form message like - `${1} cannot be blank` - `'name' cannot be blank`
// Based on the ${postioning} message get formatted
var defaultMessages Messages

func setDefaultMessage() {
	defaultMessages = Messages{
		"required":  "${0} ${1} cannot be blank",
		"email":     "${0} ${2} is not a valid e-mail address",
		"matches":   "${0} ${2} does not match the required pattern ${3}",
		"len":       "${0} ${2} should be in size ${3}",
		"min":       "${0} ${2} is less than minimum value ${3}",
		"max":       "${0} ${2} exceeds maximum value ${3}",
		"inList":    "${0} ${2} is not contained within the list ${3}",
		"unique":    "${0} ${2} already exists",
		"pinunique": "${0} ${2} already exists",
		"regexp":    "${0} ${2} is not a valid data",
	}
}

//SetDefaultErrors is used to change the Default error message
func SetDefaultErrors(msgs Messages) {
	defaultMessages = msgs
}

//GetTranslatedMessage returns the message from i8n
type GetTranslatedMessage func(langCode, module string) Messages

var translator = func(langCode, module string) Messages {
	return Messages{}
}

//SetTranslator set translator func
func SetTranslator(t GetTranslatedMessage) {
	translator = t
}

//GetModule returns the translated label value
//if translation is not available send the label
func (e Messages) GetModule(module string) error {
	return e.Get(fmt.Sprintf("%s.label", module))
}

//GetLabel returns the translated label value
//if translation is not available send the label
func (e Messages) GetLabel(module, field string) error {
	return e.Get(fmt.Sprintf("%s.%s.label", module, field))
}

//GetDefaultMessage returns the translated label value
// if translation is not available send the label
func (e Messages) GetDefaultMessage(module, field, key string, args []string) error {
	m := module
	f := field

	if err := e.GetModule(field); err != errNotFound {
		m = err.Error()
	}

	if err := e.GetLabel(module, field); err != errNotFound {
		f = err.Error()
	}

	args = append([]string{m, f}, args...)
	k := fmt.Sprintf("default.error.%s", key)
	return e.Get(k, args...)
}

//GetMessage return the message by walk through field & errCode in translation.Message
func (e Messages) GetMessage(module, field, key string, args []string) error {
	k := fmt.Sprintf("%s.%s.error.%s", module, field, key)
	return e.Get(k, args...)
}

//GetError return the validation Error message
//First it search usinng field & errCode in translation.Message
//If not it search in translation.DefaultMessage
//If not search from Default message of validation package
func (e Messages) GetError(module, field, key string, args []string) error {
	//find property specific message based on language from translation
	err := e.GetMessage(module, field, key, args)

	//If no property specific find default message from translation
	if err == errNotFound {
		if err = e.GetDefaultMessage(module, field, key, args); err == errNotFound {
			args = append([]string{module, field}, args...)
			return defaultMessages.Get(key, args...)
		}
	}

	return err
}
