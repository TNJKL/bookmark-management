package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Message represents the structure of a response message
type Message struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

//Common response messages

var (
	InternalErrResponse = Message{
		Message: "Processing error",
		Details: nil,
	}
	InputErrResponse = Message{
		Message: "Input error",
		Details: nil,
	}
)

// InputFieldError builds a Message with per-field errors, or InputErrResponse if err isn't a validation error.
func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}
	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}
	return Message{
		Message: "Input error",
		Details: errs,
	}
}
