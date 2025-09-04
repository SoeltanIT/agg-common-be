package common_test

import (
	"errors"
	"net/http"
	"testing"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name       string
		httpStatus int
		code       int
		message    string
		expected   common.Error
	}{
		{
			name:       "Create new error with all fields",
			httpStatus: http.StatusBadRequest,
			code:       4001001,
			message:    "test error",
			expected: common.Error{
				HTTPStatus: http.StatusBadRequest,
				Code:       4001001,
				Message:    "test error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.NewError(tt.httpStatus, tt.code, tt.message)
			assert.Equal(t, tt.expected, result, "NewError() should create correct Error struct")
		})
	}
}

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name     string
		err      common.Error
		expected common.Error
	}{
		{
			name:     "ErrInsufficientBalance",
			err:      common.ErrInsufficientBalance,
			expected: common.Error{HTTPStatus: http.StatusBadRequest, Code: 4001001, Message: "The player does not have sufficient balance to complete this transaction"},
		},
		{
			name:     "ErrUnauthorized",
			err:      common.ErrUnauthorized,
			expected: common.Error{HTTPStatus: http.StatusUnauthorized, Code: 4010001, Message: "You are not authorized to access this resource"},
		},
		{
			name:     "ErrForbidden",
			err:      common.ErrForbidden,
			expected: common.Error{HTTPStatus: http.StatusForbidden, Code: 4030001, Message: "You do not have permission to access this resource"},
		},
		{
			name:     "ErrServerError",
			err:      common.ErrServerError,
			expected: common.Error{HTTPStatus: http.StatusInternalServerError, Code: 5000001, Message: "An unexpected server error occurred. Please try again later."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.HTTPStatus, tt.err.HTTPStatus, "%s: HTTPStatus should match", tt.name)
			assert.Equal(t, tt.expected.Code, tt.err.Code, "%s: Code should match", tt.name)
			assert.Equal(t, tt.expected.Message, tt.err.Message, "%s: Message should match", tt.name)
		})
	}
}

func TestErrRecordNotFound(t *testing.T) {
	tests := []struct {
		name     string
		entity   string
		id       string
		expected common.Error
	}{
		{
			name:   "Record not found with entity and ID",
			entity: "user",
			id:     "123",
			expected: common.Error{
				HTTPStatus: http.StatusNotFound,
				Code:       4042999,
				Message:    "The specified user with ID '123' could not be found",
			},
		},
		{
			name:   "Record not found with empty entity and ID",
			entity: "",
			id:     "",
			expected: common.Error{
				HTTPStatus: http.StatusNotFound,
				Code:       4042999,
				Message:    "The specified  with ID '' could not be found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.ErrRecordNotFound(tt.entity, tt.id)
			assert.Equal(t, tt.expected, result, "ErrRecordNotFound should return correct error")
		})
	}
}

func TestErrorImplementsErrorInterface(t *testing.T) {
	// This test ensures that the Error type implements the error interface
	var _ error = (*common.Error)(nil)
	assert.Implements(t, (*error)(nil), common.Error{}, "Error should implement the error interface")

	// Test that we can assign an Error to an error variable
	var err error = common.Error{Message: "test"}
	assert.Equal(t, "test", err.Error(), "Error should work as an error interface")
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "Same error values",
			err:      common.ErrUnauthorized,
			target:   common.ErrUnauthorized,
			expected: true,
		},
		{
			name:     "Different error values",
			err:      common.ErrUnauthorized,
			target:   common.ErrForbidden,
			expected: false,
		},
		{
			name:     "With standard error",
			err:      errors.New("some error"),
			target:   common.ErrUnauthorized,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			assert.Equal(t, tt.expected, result, "errors.Is() should return expected result")
		})
	}
}
