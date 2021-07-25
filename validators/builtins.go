// Package validator implements value validations
package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var (

	//ErrValidation represents the default validation message
	ErrValidation = errors.New("$err:validation error")

	// ErrUnsupported is the error error returned when a validation rule
	// is used with an unsupported variable type
	errUnsupported = errors.New("unsupported type")

	// ErrBadParameter is the error returned when an invalid parameter
	// is provided to a validation rule (e.g. a string where an int was
	// expected (max=foo,len=bar) or missing a parameter when one is required (len=))
	errBadParameter = errors.New("bad parameter")

	// ErrUnknownTag is the error returned when an unknown tag is found
	errUnknownTag = errors.New("unknown tag")

	// ErrInvalid is the error returned when variable is invalid
	// (normally a nil pointer)
	errInvalid = errors.New("invalid value")
)

type field struct {
	fieldName, code, val string
	obj                  interface{}
}

//Required validates whether a variable value not null
func Required(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		valid = len(st.String()) != 0
	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		valid = st.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() != 0
	case reflect.Float32, reflect.Float64:
		valid = st.Float() != 0
	case reflect.Bool:
		valid = st.Bool()
	case reflect.Invalid:
		valid = false // always invalid
	case reflect.Struct:
		valid = true // always valid since only nil pointers are empty
	default:
		if t, ok := v.(time.Time); ok {
			valid = !t.IsZero()
		}

		return errUnsupported
	}

	if !valid {
		return ErrValidation

	}

	return nil
}

//Length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func Length(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	valid := true
	if st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return nil
		}
		st = st.Elem()
	}
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		valid = int64(len(st.String())) == p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		valid = int64(st.Len()) == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		valid = st.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return errBadParameter
		}
		valid = st.Uint() == p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return errBadParameter
		}
		valid = st.Float() == p
	default:
		return errUnsupported
	}

	if !valid {
		return ErrValidation
	}

	return nil
}

// Min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func Min(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	invalid := false
	if st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return nil
		}
		st = st.Elem()
	}
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		invalid = int64(len(st.String())) < p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		invalid = int64(st.Len()) < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		invalid = st.Int() < p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return errBadParameter
		}
		invalid = st.Uint() < p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return errBadParameter
		}
		invalid = st.Float() < p
	default:
		return errUnsupported
	}

	if invalid {
		return ErrValidation
	}

	return nil
}

// Max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func Max(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	var invalid bool
	if st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return nil
		}
		st = st.Elem()
	}
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		invalid = int64(len(st.String())) > p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		invalid = int64(st.Len()) > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(param)
		if err != nil {
			return errBadParameter
		}
		invalid = st.Int() > p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(param)
		if err != nil {
			return errBadParameter
		}
		invalid = st.Uint() > p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(param)
		if err != nil {
			return errBadParameter
		}
		invalid = st.Float() > p
	default:
		return errUnsupported
	}

	if invalid {
		return ErrValidation
	}

	return nil
}

// Regex is the builtin validation function that checks
// whether the string variable matches a regular expression
func Regex(v interface{}, param string) error {
	s, ok := v.(string)
	if !ok {
		sptr, ok := v.(*string)
		if !ok {
			return errUnsupported
		}
		if sptr == nil {
			return nil
		}
		s = *sptr
	}

	re, err := regexp.Compile(param)
	if err != nil {
		return errBadParameter
	}

	if !re.MatchString(s) {
		return ErrValidation
	}

	return nil
}

//Email validation using regex
func Email(v interface{}, param string) error {
	emailRegexString := "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	return Regex(v, emailRegexString)
}

//InList execute the stringer and validate against the param
//Default in in-valid param is 'Unknown'
func InList(v interface{}, param string) error {
	if param == "" {
		param = "Unknown"
	}

	if d := fmt.Sprintf("%s", v); d == param {
		return ErrValidation
	}
	return nil
}

// asInt retuns the parameter as a int64
// or panics if it can't convert
func asInt(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return 0, errBadParameter
	}
	return i, nil
}

// asUint retuns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) (uint64, error) {
	i, err := strconv.ParseUint(param, 0, 64)
	if err != nil {
		return 0, errBadParameter
	}
	return i, nil
}

// asFloat retuns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) (float64, error) {
	i, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0.0, errBadParameter
	}
	return i, nil
}
