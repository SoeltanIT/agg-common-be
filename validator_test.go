package common_test

import (
	"testing"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	t.Run("Validator returns a non-nil validator instance", func(t *testing.T) {
		v := common.Validator()
		assert.NotNil(t, v, "Validator() should return a non-nil validator instance")
	})
}

func TestValidatorWithCustomTags(t *testing.T) {
	type TestStruct struct {
		Username string `validate:"required,min=3" json:"username"`
		Email    string `validate:"required,email" json:"email"`
	}

	tests := []struct {
		name     string
		input    TestStruct
		hasErr   bool
		errCount int
	}{
		{
			name:     "Valid input",
			input:    TestStruct{Username: "testuser", Email: "test@example.com"},
			hasErr:   false,
			errCount: 0,
		},
		{
			name:     "Missing required fields",
			input:    TestStruct{},
			hasErr:   true,
			errCount: 2,
		},
		{
			name:     "Invalid email format",
			input:    TestStruct{Username: "testuser", Email: "invalid-email"},
			hasErr:   true,
			errCount: 1,
		},
		{
			name:     "Username too short",
			input:    TestStruct{Username: "ab", Email: "test@example.com"},
			hasErr:   true,
			errCount: 1,
		},
	}

	v := common.Validator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if tt.hasErr {
				assert.Error(t, err, "Expected validation error")
				if err != nil {
					validationErrors := err.(validator.ValidationErrors)
					assert.Len(t, validationErrors, tt.errCount, "Unexpected number of validation errors")
				}
			} else {
				assert.NoError(t, err, "Expected no validation errors")
			}
		})
	}
}
