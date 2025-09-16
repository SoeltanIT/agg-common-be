package common

import (
	"fmt"
	"net/http"
)

// Error is a custom error type for API responses
type Error struct {
	HTTPStatus int
	Code       int
	Message    string
}

// Error implements the error interface, returns the error message
func (e Error) Error() string {
	return e.Message
}

// NewError creates a new Error instance
func NewError(httpStatus int, code int, message string) Error {
	return Error{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// Constants for common errors
var (
	// Error 400
	ErrInsufficientBalance        = Error{HTTPStatus: http.StatusBadRequest, Code: 4001001, Message: "The player does not have sufficient balance to complete this transaction"}
	ErrInvalidBonus               = Error{HTTPStatus: http.StatusBadRequest, Code: 4001002, Message: "The bonus provided is invalid or no longer available"}
	ErrEmptyClientId              = Error{HTTPStatus: http.StatusBadRequest, Code: 4001003, Message: "Client ID is missing. Please provide a valid Client ID"}
	ErrInvalidClientSecret        = Error{HTTPStatus: http.StatusBadRequest, Code: 4001004, Message: "The provided client secret is invalid"}
	ErrGameInActive               = Error{HTTPStatus: http.StatusBadRequest, Code: 4001005, Message: "The selected game is currently inactive"}
	ErrInvalidSignature           = Error{HTTPStatus: http.StatusBadRequest, Code: 4001006, Message: "The request signature is invalid. Please check your credentials"}
	ErrMissingAggregatorSignature = Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       4002001,
		Message:    "Missing X-Aggregator-Signature header. Please provide a valid signature",
	}

	// Error 401
	ErrUnauthorized         = Error{HTTPStatus: http.StatusUnauthorized, Code: 4010001, Message: "You are not authorized to access this resource"}
	ErrInvalidToken         = Error{HTTPStatus: http.StatusUnauthorized, Code: 4010002, Message: "The provided access token is invalid"}
	ErrMissingAuthorization = Error{HTTPStatus: http.StatusUnauthorized, Code: 4010003, Message: "Authorization header is missing"}
	ErrExpiredToken         = Error{HTTPStatus: http.StatusUnauthorized, Code: 4010004, Message: "The access token has expired. Please login again"}

	// Error 403
	ErrForbidden      = Error{HTTPStatus: http.StatusForbidden, Code: 4030001, Message: "You do not have permission to access this resource"}
	ErrSessionExpired = Error{HTTPStatus: http.StatusForbidden, Code: 4031001, Message: "Your session has expired"}

	// Error 404
	ErrProviderNotFound     = Error{HTTPStatus: http.StatusNotFound, Code: 4041001, Message: "The specified game provider could not be found"}
	ErrSessionNotFound      = Error{HTTPStatus: http.StatusNotFound, Code: 4041002, Message: "The session you are trying to access does not exist or is invalid"}
	ErrPlayerNotFound       = Error{HTTPStatus: http.StatusNotFound, Code: 4041003, Message: "The requested player could not be found."}
	ErrGameNotFound         = Error{HTTPStatus: http.StatusNotFound, Code: 4041004, Message: "The requested game could not be found"}
	ErrDuplicateTransaction = Error{HTTPStatus: http.StatusConflict, Code: 4091001, Message: "This transaction has already been processed"}
	ErrRecordNotFound       = func(entity, id string) Error {
		return Error{
			HTTPStatus: http.StatusNotFound,
			Code:       4042999,
			Message:    fmt.Sprintf("The specified %s with ID '%s' could not be found", entity, id),
		}
	}

	// Error 5xx
	ErrServerError = Error{HTTPStatus: http.StatusInternalServerError, Code: 5000001, Message: "An unexpected server error occurred. Please try again later."}
)

func ValidationError(message string) Error {
	return NewError(http.StatusBadRequest, 4002999, message)
}
