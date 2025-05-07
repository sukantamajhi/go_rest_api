package utils

import (
	"log"

	"github.com/go-playground/validator/v10"
)

func GetErrorMsg(fe validator.FieldError) string {
	log.Printf("Field error: %+v", fe)
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return "Should be at least " + fe.Param() + " characters"
	case "email":
		return "Invalid email address"
	case "unique":
		return "This field must be unique"
	case "eqfield":
		return "This field must be equal to " + fe.Param()
	case "nefield":
		return "This field must not be equal to " + fe.Param()
	case "gt":
		return "This field must be greater than " + fe.Param()
	case "gte":
		return "This field must be greater than or equal to " + fe.Param()
	case "lte":
		return "This field must be less than or equal to " + fe.Param()
	case "lt":
		return "This field must be less than " + fe.Param()
	case "len":
		return "This field must be " + fe.Param() + " characters long"
	case "max":
		return "This field must be less than or equal to " + fe.Param() + " characters"
	}
	return "Unknown error"
}

func GetErrorMsgs(fe validator.ValidationErrors) []string {
	var errorMessages []string
	for _, err := range fe {
		errorMessages = append(errorMessages, GetErrorMsg(err))
	}
	return errorMessages
}
