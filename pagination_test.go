package common_test

import (
	"strconv"
	"testing"

	common "github.com/SoeltanIT/agg-common-be"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestNewPaginationParams(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  map[string]string
		expectedPage int
		expectedSize int
	}{
		{
			name:         "Default values",
			queryParams:  map[string]string{},
			expectedPage: 1,
			expectedSize: common.DefaultPageSize,
		},
		{
			name: "Custom page and size",
			queryParams: map[string]string{
				"page":     "2",
				"pageSize": "25",
			},
			expectedPage: 2,
			expectedSize: 25,
		},
		{
			name: "Invalid page number",
			queryParams: map[string]string{
				"page":     "-1",
				"pageSize": "10",
			},
			expectedPage: 1, // Should default to 1 for invalid values
			expectedSize: 10,
		},
		{
			name: "Non-numeric values",
			queryParams: map[string]string{
				"page":     "abc",
				"pageSize": "xyz",
			},
			expectedPage: 1, // Default values for non-numeric
			expectedSize: common.DefaultPageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			query := ctx.Request().URI().QueryArgs()
			for key, value := range tt.queryParams {
				query.Set(key, value)
			}

			params := common.NewPaginationParams(ctx)

			assert.Equal(t, tt.expectedPage, params.GetPage(), "Page number should match")
			assert.Equal(t, tt.expectedSize, params.GetPageSize(), "Page size should match")
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "First page",
			page:     1,
			pageSize: 10,
			expected: 0,
		},
		{
			name:     "Second page",
			page:     2,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "Third page with custom size",
			page:     3,
			pageSize: 25,
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			query := ctx.Request().URI().QueryArgs()
			query.Set("page", strconv.Itoa(tt.page))
			query.Set("pageSize", strconv.Itoa(tt.pageSize))

			p := common.NewPaginationParams(ctx)
			offset := p.CalculateOffset()
			assert.Equal(t, tt.expected, offset, "Offset should be calculated correctly")
		})
	}
}

func TestGetNext(t *testing.T) {
	tests := []struct {
		name          string
		currentPage   int
		pageSize      int
		totalItems    int64
		expectHasNext bool
	}{
		{
			name:          "Has next page",
			currentPage:   1,
			pageSize:      10,
			totalItems:    20,
			expectHasNext: true,
		},
		{
			name:          "No next page",
			currentPage:   2,
			pageSize:      10,
			totalItems:    15,
			expectHasNext: false,
		},
		{
			name:          "Exact page boundary",
			currentPage:   2,
			pageSize:      10,
			totalItems:    20,
			expectHasNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			query := ctx.Request().URI().QueryArgs()
			query.Set("page", strconv.Itoa(tt.currentPage))
			query.Set("pageSize", strconv.Itoa(tt.pageSize))

			p := common.NewPaginationParams(ctx)

			nextURL := p.GetNext(ctx.Request().URI(), tt.totalItems)

			if tt.expectHasNext {
				assert.NotEmpty(t, nextURL, "Next URL should not be empty when there are more items")
				assert.Contains(t, nextURL, "page="+strconv.Itoa(tt.currentPage+1), "Next URL should point to the next page")
			} else {
				assert.Empty(t, nextURL, "Next URL should be empty when there are no more items")
			}
		})
	}
}

func TestGetPrev(t *testing.T) {
	tests := []struct {
		name          string
		currentPage   int
		pageSize      int
		expectHasPrev bool
	}{
		{
			name:          "Has previous page",
			currentPage:   2,
			pageSize:      10,
			expectHasPrev: true,
		},
		{
			name:          "No previous page on first page",
			currentPage:   1,
			pageSize:      10,
			expectHasPrev: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			query := ctx.Request().URI().QueryArgs()
			query.Set("page", strconv.Itoa(tt.currentPage))
			query.Set("pageSize", strconv.Itoa(tt.pageSize))

			p := common.NewPaginationParams(ctx)

			prevURL := p.GetPrev(ctx.Request().URI())

			if tt.expectHasPrev {
				assert.NotEmpty(t, prevURL, "Previous URL should not be empty when not on first page")
				assert.Contains(t, prevURL, "page="+strconv.Itoa(tt.currentPage-1), "Previous URL should point to the previous page")
			} else {
				assert.Empty(t, prevURL, "Previous URL should be empty on first page")
			}
		})
	}
}

func TestGetPaginationResponse(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	query := ctx.Request().URI().QueryArgs()
	query.Set("page", strconv.Itoa(2))
	query.Set("pageSize", strconv.Itoa(10))

	p := common.NewPaginationParams(ctx)

	resp := p.GetPaginationResponse(ctx.Request(), 25)

	assert.Equal(t, int64(25), resp.Total, "Total count should match")
	assert.Equal(t, 2, resp.Page, "Current page should match")
	assert.Contains(t, resp.Next, "page=3", "Next page URL should point to page 3")
	assert.Contains(t, resp.Prev, "page=1", "Previous page URL should point to page 1")
}
