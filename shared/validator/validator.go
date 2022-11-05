package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func contains[V string | int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | any](s []V, e V) bool {
	flag := false
	for _, v := range s {
		flag = flag || reflect.DeepEqual(v, e)
	}
	return flag
}

func oneOfField(fl validator.FieldLevel) bool {
	currentValue, currentKind, _ := fl.ExtractType(fl.Field())
	reflValue, reflKind, _, found := fl.GetStructFieldOK2()
	if !found || reflKind != reflect.Slice || reflValue.Type().Elem().Kind() != currentKind {
		return false
	}

	switch v := currentValue.Interface().(type) {
	case string:
		return contains(reflValue.Interface().([]string), v)
	case int8:
		return contains(reflValue.Interface().([]int8), v)
	case int16:
		return contains(reflValue.Interface().([]int16), v)
	case int32:
		return contains(reflValue.Interface().([]int32), v)
	case int64:
		return contains(reflValue.Interface().([]int64), v)
	case int:
		return contains(reflValue.Interface().([]int), v)
	case uint8:
		return contains(reflValue.Interface().([]uint8), v)
	case uint16:
		return contains(reflValue.Interface().([]uint16), v)
	case uint32:
		return contains(reflValue.Interface().([]uint32), v)
	case uint64:
		return contains(reflValue.Interface().([]uint64), v)
	case uintptr:
		return contains(reflValue.Interface().([]uintptr), v)
	case float32:
		return contains(reflValue.Interface().([]float32), v)
	case float64:
		return contains(reflValue.Interface().([]float64), v)
	default:
		var s []any
		var ok bool
		if s, ok = reflValue.Interface().([]any); !ok {
			return false
		}
		return contains(s, v)
	}
}

// schemaTagName for register new name field by tag
func schemaTagName(fld reflect.StructField) string {
	name := fld.Tag.Get("schema")
	if name != "" {
		return name
	}
	return fld.Name
}

func init() {
	v = validator.New()
	// register custom validator
	v.RegisterTagNameFunc(schemaTagName)
	v.RegisterValidation("oneoffield", oneOfField, false)
}

type ValidationErrors = validator.ValidationErrors
type FieldError = validator.FieldError

func Get() *validator.Validate {
	return v
}
