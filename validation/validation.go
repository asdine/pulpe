package validation

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.CustomTypeTagMap.Set("pulpe-json", validateJSON)
}

func validateJSON(i interface{}, o interface{}) bool {
	switch t := i.(type) {
	case []byte:
		return govalidator.IsJSON(string(t))
	case string:
		return govalidator.IsJSON(t)
	case json.RawMessage:
		return govalidator.IsJSON(string(t))
	case *json.RawMessage:
		if t != nil {
			return validateJSON(*t, o)
		}
	}
	return true
}

// Validate validates and saves all the govalidator errors in a ValidatorError.
func Validate(s interface{}) error {
	ok, err := govalidator.ValidateStruct(s)
	if ok {
		return nil
	}

	errs, ok := err.(govalidator.Errors)
	if !ok || errs == nil {
		return nil
	}

	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var verr validationError

	for i := range errs {
		e, ok := errs[i].(govalidator.Error)
		if !ok {
			continue
		}

		f, ok := typ.FieldByName(e.Name)
		if !ok {
			// shouldn't happen.
			panic("unknown field")
		}

		var name = e.Name

		tag := f.Tag.Get("json")
		if tag != "" {
			if idx := strings.Index(tag, ","); idx != -1 {
				name = tag[:idx]
			} else {
				name = tag
			}
		}

		verr = AddError(verr, name, e.Err).(validationError)
	}

	return verr
}
