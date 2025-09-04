package common

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type response struct {
	HttpStatus int
	Code       int                `json:"code,omitempty"`
	Status     string             `json:"status,omitempty"`
	Message    string             `json:"message,omitempty"`
	Data       any                `json:"data,omitempty"`
	Errors     interface{}        `json:"errors,omitempty"`
	Pagination paginationResponse `json:"pagination,omitempty"`
}

// Response create new response instance
func Response() *response {
	return &response{}
}

// SetError sets the error response
func (r *response) SetError(err error) *response {
	r.Status = "failed"

	// Validation errors (go-playground/validator)
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		if errors.As(err, &vErrs) {
			var errs []string
			trans, _ := translator.GetTranslator("en")
			for _, vErr := range vErrs {
				errs = append(errs, vErr.Translate(trans))
			}

			r.Errors = errs
			r.Code = 4002000
			r.Message = "Validation failed"
			r.HttpStatus = 400

			return r
		}
	}

	// Custom Error
	var cuserr Error
	if errors.As(err, &cuserr) {
		status := cuserr.HTTPStatus
		if status == 0 {
			status = http.StatusInternalServerError
		}

		r.Code = cuserr.Code
		r.Message = cuserr.Message
		r.HttpStatus = status
		return r
	}

	// Fiber Error
	var fErr *fiber.Error
	if errors.As(err, &fErr) {
		r.Code = fErr.Code * 10000
		r.Message = fErr.Message
		r.HttpStatus = fErr.Code

		return r
	}

	// Fallback
	r.Code = 5000001
	r.Message = "An unexpected server error occurred. Please try again later."
	r.HttpStatus = 500

	return r
}

// SetData sets the data response, status code is optional, default is 200
func (r *response) SetData(data any, status ...int) *response {
	r.Status = "success"
	r.Data = data

	if len(status) > 0 {
		r.HttpStatus = status[0]
	}

	if r.HttpStatus >= http.StatusMultipleChoices || r.HttpStatus < http.StatusOK {
		r.HttpStatus = http.StatusOK
	}

	return r
}

// SetMessage sets the message response
func (r *response) SetMessage(message string) *response {
	r.Message = message
	return r
}

// SetPagination sets the pagination response
func (r *response) SetPagination(pagination paginationResponse) *response {
	r.Pagination = pagination
	return r
}

// Send sends the response, if HttpStatus is less than 200, it will be set to 200
func (r *response) Send(ctx *fiber.Ctx) error {
	if r.HttpStatus >= http.StatusOK {
		ctx = ctx.Status(r.HttpStatus)
	}

	return ctx.JSON(r)
}
