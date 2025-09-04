package common_test

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestResponse_NewResponse(t *testing.T) {
	r := common.Response()
	assert.NotNil(t, r, "Response() should return a new response instance")
	assert.Equal(t, 0, r.HttpStatus, "New response should have zero HttpStatus")
	assert.Empty(t, r.Status, "New response should have empty Status")
	assert.Empty(t, r.Message, "New response should have empty Message")
	assert.Nil(t, r.Data, "New response should have nil Data")
	assert.Nil(t, r.Errors, "New response should have nil Errors")
}

type testData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestResponse_SetData(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		status     int
		expectData interface{}
		expectCode int
	}{
		{
			name:       "Success with data and status",
			data:       testData{Name: "John", Age: 30},
			status:     http.StatusCreated,
			expectData: testData{Name: "John", Age: 30},
			expectCode: http.StatusCreated,
		},
		{
			name:       "Success with data and default status",
			data:       "test data",
			expectData: "test data",
			expectCode: http.StatusOK,
		},
		{
			name:       "Success with data and invalid status",
			status:     http.StatusInternalServerError,
			data:       "test data",
			expectData: "test data",
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := common.Response()
			if tt.status != 0 {
				r.SetData(tt.data, tt.status)
			} else {
				r.SetData(tt.data)
			}

			assert.Equal(t, "success", r.Status, "Status should be 'success'")
			assert.Equal(t, tt.expectData, r.Data, "Data should match expected")
			if tt.status != 0 {
				assert.Equal(t, tt.expectCode, r.HttpStatus, fmt.Sprintf("Http status should be %d", r.HttpStatus))
			} else {
				assert.Equal(t, tt.expectCode, http.StatusOK, "Http status should be 200")
			}
		})
	}
}

type mockValidationError struct{}

func (m mockValidationError) Error() string     { return "validation error" }
func (m mockValidationError) All() []error      { return []error{errors.New("field is required")} }
func (m mockValidationError) ErrorOrNil() error { return m }

func TestResponse_SetError(t *testing.T) {
	// Setup validator
	v := validator.New()
	type TestStruct struct {
		Name string `validate:"required"`
	}

	// Create a validation error
	testData := TestStruct{}
	err := v.Struct(testData)
	require.Error(t, err)

	tests := []struct {
		name           string
		err            error
		expectStatus   string
		expectMessage  string
		expectCode     int
		expectHttpCode int
		hasErrors      bool
	}{
		{
			name:           "Validation error",
			err:            err,
			expectStatus:   "failed",
			expectMessage:  "Validation failed",
			expectCode:     4002000,
			hasErrors:      true,
			expectHttpCode: 400,
		},
		{
			name:           "Custom error",
			err:            common.NewError(http.StatusNotFound, 4041001, "Not found"),
			expectStatus:   "failed",
			expectMessage:  "Not found",
			expectCode:     4041001,
			expectHttpCode: 404,
		},
		{
			name:           "Custom error with HttpStatus zero",
			err:            common.NewError(0, 5000001, "Internal server error"),
			expectStatus:   "failed",
			expectMessage:  "Internal server error",
			expectCode:     5000001,
			expectHttpCode: 500,
		},
		{
			name:           "Fiber error",
			err:            fiber.NewError(http.StatusForbidden, "Access denied"),
			expectStatus:   "failed",
			expectMessage:  "Access denied",
			expectCode:     http.StatusForbidden * 10000,
			expectHttpCode: 0,
		},
		{
			name:           "Standard error",
			err:            errors.New("something went wrong"),
			expectStatus:   "failed",
			expectMessage:  "An unexpected server error occurred. Please try again later.",
			expectCode:     5000001,
			expectHttpCode: 0,
		},
		{
			name:           "Nil error",
			err:            nil,
			expectStatus:   "failed",
			expectMessage:  "An unexpected server error occurred. Please try again later.",
			expectCode:     5000001,
			expectHttpCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := common.Response().SetError(tt.err)

			assert.Equal(t, tt.expectStatus, r.Status, "Status should be 'failed'")
			assert.Equal(t, tt.expectMessage, r.Message, "Error message should match")
			assert.Equal(t, tt.expectCode, r.Code, "Error code should match")

			if tt.hasErrors {
				assert.NotNil(t, r.Errors, "Validation errors should be set")
			}

			if tt.expectHttpCode != 0 {
				assert.Equal(t, tt.expectHttpCode, r.HttpStatus, "HTTP status code should be set")
			}
		})
	}
}

func TestResponse_SetPagination(t *testing.T) {
	r := common.Response()

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	query := ctx.Request().URI().QueryArgs()
	query.Set("page", strconv.Itoa(1))
	query.Set("pageSize", strconv.Itoa(10))

	p := common.NewPaginationParams(ctx)

	pagination := p.GetPaginationResponse(ctx.Request(), 50)

	r.SetPagination(pagination)

	assert.Equal(t, pagination, r.Pagination, "Pagination should be set correctly")
}

func TestResponse_Send(t *testing.T) {
	app := fiber.New()

	t.Run("Success response with data", func(t *testing.T) {
		app.Get("/test", func(c *fiber.Ctx) error {
			r := common.Response().SetData("test data")
			return r.Send(c)
		})

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err, "Should not return error")

		resp, err := app.Test(req, -1)
		require.NoError(t, err, "Should not return error")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be 200")
	})

	t.Run("Success response with custom status", func(t *testing.T) {
		app.Get("/test-created", func(c *fiber.Ctx) error {
			r := common.Response().SetData("created", http.StatusCreated)
			return r.Send(c)
		})

		req, err := http.NewRequest(http.MethodGet, "/test-created", nil)
		require.NoError(t, err, "Should not return error")

		resp, err := app.Test(req, -1)
		require.NoError(t, err, "Should not return error")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status code should be 201")
	})
}

func TestResponse_SendWithError(t *testing.T) {
	app := fiber.New()

	t.Run("Custom error response", func(t *testing.T) {
		app.Get("/not-found", func(c *fiber.Ctx) error {
			err := common.NewError(http.StatusNotFound, 4041001, "Not found")
			return common.Response().SetError(err).Send(c)
		})

		req, err := http.NewRequest(http.MethodGet, "/not-found", nil)
		require.NoError(t, err, "Should not return error")

		resp, err := app.Test(req, -1)
		require.NoError(t, err, "Should not return error")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status code should be 404")
	})

	t.Run("Validation error response", func(t *testing.T) {
		app.Get("/validation-error", func(c *fiber.Ctx) error {
			type TestStruct struct {
				Name string `validate:"required"`
			}

			testData := TestStruct{}
			err := validator.New().Struct(testData)
			return common.Response().SetError(err).Send(c)
		})

		req, err := http.NewRequest(http.MethodGet, "/validation-error", nil)
		require.NoError(t, err, "Should not return error")

		resp, err := app.Test(req, -1)
		require.NoError(t, err, "Should not return error")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code should be 400")
	})

	t.Run("Fiber error response", func(t *testing.T) {
		app.Get("/fiber-error", func(c *fiber.Ctx) error {
			return common.Response().SetError(fiber.ErrForbidden).Send(c)
		})

		req, err := http.NewRequest(http.MethodGet, "/fiber-error", nil)
		require.NoError(t, err, "Should not return error")

		resp, err := app.Test(req, -1)
		require.NoError(t, err, "Should not return error")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode, "Status code should be 403")
	})
}

func TestResponse_SetMessage(t *testing.T) {
	r := common.Response().SetMessage("Custom message")
	assert.Equal(t, "Custom message", r.Message, "Message should be set")
}
