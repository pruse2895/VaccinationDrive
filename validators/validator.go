package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

//Errors give type to holds the validation error messages
type Errors map[string]error

//MarshalJSON implements text marshaller
func (errs Errors) MarshalJSON() ([]byte, error) {
	err := make(map[string]string)
	for k, e := range errs {
		err[k] = e.Error()
	}

	return json.Marshal(err)
}

// ValidationFunc is a function that receives the value of a
// field and a parameter used for the respective validation tag.
type ValidationFunc func(v interface{}, param string) error

// Validator implements a validator
type Validator struct {
	// Tag name being used to get validation keys.
	tagName string
	// validationFuncs is a map of validation funcs specified by their name.
	validationFuncs map[string]ValidationFunc
	//errors holds the errror messages of the module
	errors Errors
	//module has the module name
	module string
	//for internationalization translate flag is set to true and to specify lang-code
	translate bool
	langCode  string
	//message holds the module speific messages
	Messages Messages
}

// New creates a new Validator with default validation funcs
func New(module string) *Validator {
	// if module == "" {
	// return nil, errors.New("$err:module cannot be empty")
	// }

	return &Validator{
		tagName: "validate",
		validationFuncs: map[string]ValidationFunc{
			"required": Required,
			"len":      Length,
			"min":      Min,
			"max":      Max,
			"email":    Email,
			"regexp":   Regex,
			"inList":   InList,
		},
		errors: Errors{},
		module: module,
	}
}

// Copy a validator
func (mv *Validator) copy(module string) *Validator {
	newFuncs := map[string]ValidationFunc{}
	for k, f := range mv.validationFuncs {
		newFuncs[k] = f
	}
	v := &Validator{
		tagName:         mv.tagName,
		validationFuncs: newFuncs,
		errors:          Errors{},
		module:          module,
	}

	if mv.translate {
		v.Translate(mv.langCode)
	}
	return v
}

// SetTag allows you to change the tag name used in structs
func (mv *Validator) SetTag(tag string) {
	mv.tagName = tag
}

// SetValidationFunc sets the function to be used for a given
// validation constraint. Calling this function with nil vf
// is the same as removing the constraint function from the list.
func (mv *Validator) SetValidationFunc(name string, vf ValidationFunc) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if vf == nil {
		delete(mv.validationFuncs, name)
		return nil
	}
	mv.validationFuncs[name] = vf
	return nil
}

//Translate enable translation
//And also fetch & set the multi language message
func (mv *Validator) Translate(langCode string) error {
	if langCode == "" {
		return errors.New("$err:lnaguage code not defined")
	}
	mv.translate = true
	mv.langCode = langCode
	mv.Messages = translator(langCode, mv.module)
	return nil
}

//GetAndAddError get error form message and add the error to the error map if not present
func (mv *Validator) GetAndAddError(field, errKey string, args ...string) {
	if field == "" {
		return
	}

	if _, ok := mv.errors[field]; !ok {
		err := mv.Messages.GetError(mv.module, field, errKey, args)
		if err == errNotFound {
			err = ErrValidation
		}

		mv.AddError(field, err)
	}
}

//AddError add the error to the error map if not present
func (mv *Validator) AddError(field string, err error) {
	if field == "" {
		return
	}

	if mv.errors == nil {
		mv.errors = Errors{}
	}

	if _, ok := mv.errors[field]; !ok {
		mv.errors[field] = err
	}
}

//RemoveError delete the error from errors map
func (mv *Validator) RemoveError(key string) {
	delete(mv.errors, key)
}

//AppendErrors append multiple errors
func (mv *Validator) AppendErrors(fname string, errs Errors) {
	for k, e := range errs {
		mv.AddError(fname+"."+k, e)
	}
}

//Result return validation errors & hasErr
func (mv *Validator) Result() (Errors, error) {
	if len(mv.errors) == 0 {
		return mv.errors, nil
	}

	return mv.errors, ErrValidation
}

// Validate validates the fields of a struct based
// on 'validator' tags and returns errors found indexed
// by the field name.
func (mv *Validator) Validate(v interface{}) (Errors, error) {
	sv := reflect.ValueOf(v)
	st := reflect.TypeOf(v)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return mv.Validate(sv.Elem().Interface())
	}
	if sv.Kind() != reflect.Struct && sv.Kind() != reflect.Interface {
		return mv.errors, errUnsupported
	}

	nfields := sv.NumField()
	for i := 0; i < nfields; i++ {
		fname := st.Field(i).Name
		if !unicode.IsUpper(rune(fname[0])) {
			continue
		}

		f := sv.Field(i)
		// deal with pointers
		for f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem()
		}
		tag := st.Field(i).Tag.Get(mv.tagName)
		if tag == "-" {
			continue
		}

		if tag != "" {
			mv.Valid(f.Interface(), fname, tag)
		}

		mv.deepValidateCollection(f, fname) // no-op if field is not a struct, interface, array, slice or map
	}

	return mv.Result()
}

func (mv *Validator) deepValidateCollection(f reflect.Value, fname string) {
	switch f.Kind() {
	case reflect.Struct, reflect.Interface, reflect.Ptr:
		v := mv.copy(fname)
		e, _ := v.Validate(f.Interface())
		mv.AppendErrors(fname, e)
	case reflect.Array, reflect.Slice:
		for i := 0; i < f.Len(); i++ {
			mv.deepValidateCollection(f.Index(i), fmt.Sprintf("%s[%d]", fname, i))
		}
	case reflect.Map:
		for _, key := range f.MapKeys() {
			mv.deepValidateCollection(key, fmt.Sprintf("%s[%+v](key)", fname, key.Interface())) // validate the map key
			value := f.MapIndex(key)
			mv.deepValidateCollection(value, fmt.Sprintf("%s[%+v](value)", fname, key.Interface()))
		}
	}
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func (mv *Validator) Valid(val interface{}, fname, tags string) error {
	if tags == "-" {
		return nil
	}
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return mv.Valid(v.Elem().Interface(), fname, tags)
	}
	switch v.Kind() {
	case reflect.Invalid:
		mv.validateVar(nil, fname, tags)
	default:
		mv.validateVar(val, fname, tags)
	}
	return nil
}

// validateVar validates one single variable
func (mv *Validator) validateVar(v interface{}, fname, tag string) {
	tags, err := mv.parseTags(tag)
	if err != nil {
		// unknown tag found, give up.
		mv.AddError(fname, err)
		return
	}

	mv.ValidateField(fname, v, tags)
}

//ValidateField execute the validation func
func (mv *Validator) ValidateField(fname string, v interface{}, tags []Tag) bool {
	for _, t := range tags {
		if err := t.Fn(v, t.Param); err != nil {
			if err == ErrValidation {
				mv.GetAndAddError(fname, t.Name, fmt.Sprintf("%v", v), t.Param)
				return false
			}
			log.Printf("$err: error validating, err: %s", err.Error())
		}
	}

	return true
}

// Tag represents one of the tag items
type Tag struct {
	Name  string         // name of the tag
	Fn    ValidationFunc // validation function to call
	Param string         // parameter to send to the validation function
}

// separate by no escaped commas
var sepPattern = regexp.MustCompile(`((?:^|[^\\])(?:\\\\)*),`)

func splitUnescapedComma(str string) []string {
	ret := []string{}
	indexes := sepPattern.FindAllStringIndex(str, -1)
	last := 0
	for _, is := range indexes {
		ret = append(ret, str[last:is[1]-1])
		last = is[1]
	}
	ret = append(ret, str[last:])
	return ret
}

// parseTags parses all individual tags found within a struct tag.
func (mv *Validator) parseTags(t string) ([]Tag, error) {
	tl := splitUnescapedComma(t)
	tags := make([]Tag, 0, len(tl))
	for _, i := range tl {
		i = strings.Replace(i, `\,`, ",", -1)
		tg := Tag{}
		v := strings.SplitN(i, "=", 2)
		tg.Name = strings.Trim(v[0], " ")
		if tg.Name == "" {
			return []Tag{}, errUnknownTag
		}
		if len(v) > 1 {
			tg.Param = strings.Trim(v[1], " ")
		}
		var found bool
		if tg.Fn, found = mv.validationFuncs[tg.Name]; !found {
			return []Tag{}, errUnknownTag
		}
		tags = append(tags, tg)
	}
	return tags, nil
}
