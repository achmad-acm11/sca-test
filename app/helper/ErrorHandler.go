package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"sca-integrator/app/exception"
	"strings"
)

func ErrorHandler(err error) {
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}
}

func customMessage(field string, typeMessage string, param string) string {
	var msg string
	switch typeMessage {
	case "required":
		msg = fmt.Sprintf("%s %s", field, "field is required")
	case "oneof":
		msg = fmt.Sprintf("%s %s", field, "field is oneof "+param)
	}
	return msg
}

func ErrorHandlerValidator(err error) {
	if err != nil {
		errors := make(map[string]string)

		for _, e := range err.(validator.ValidationErrors) {
			msg := customMessage(strings.ToLower(e.Field()), e.Tag(), e.Param())
			errors[strings.ToLower(e.Field())] = msg
		}
		panic(exception.NewValidationError(errors))
	}
}
