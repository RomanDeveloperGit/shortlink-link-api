package requestvalidation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func ParseErrors(err error) []FieldError {
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return []FieldError{{Field: "", Error: err.Error()}}
	}

	result := make([]FieldError, 0, len(validationErrs))
	for _, e := range validationErrs {
		result = append(result, FieldError{
			Field: e.Field(),
			Error: formatMessage(e),
		})
	}
	return result
}

func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	// Строковые
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "url":
		return "Must be a valid URL"
	case "uuid", "uuid4":
		return "Must be a valid UUID"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "alpha":
		return "Must contain only letters"
	case "numeric":
		return "Must contain only numbers"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", e.Param())
	case "contains":
		return fmt.Sprintf("Must contain %q", e.Param())
	case "startswith":
		return fmt.Sprintf("Must start with %q", e.Param())
	case "endswith":
		return fmt.Sprintf("Must end with %q", e.Param())

	// Длины строк/слайсов
	case "min":
		return fmt.Sprintf("Must be at least %s", e.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", e.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s", e.Param())

	// Числовые сравнения
	case "gt":
		return fmt.Sprintf("Must be greater than %s", e.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", e.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", e.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", e.Param())
	case "eq":
		return fmt.Sprintf("Must be equal to %s", e.Param())
	case "ne":
		return fmt.Sprintf("Must not be equal to %s", e.Param())
	case "gtfield":
		return fmt.Sprintf("Must be greater than field %q", e.Param())
	case "gtefield":
		return fmt.Sprintf("Must be greater than or equal to field %q", e.Param())
	case "ltfield":
		return fmt.Sprintf("Must be less than field %q", e.Param())
	case "ltefield":
		return fmt.Sprintf("Must be less than or equal to field %q", e.Param())

	// Разное
	case "dive":
		return "One of the items is invalid"
	case "required_if":
		return "This field is required"
	case "required_unless":
		return "This field is required"
	case "excluded_if":
		return "This field must not be set"

	default:
		return e.Error()
	}
}
