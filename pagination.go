package common

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/valyala/fasthttp"
)

// PaginationParams : contains the pagination parameters.
type PaginationParams struct {
	Page     int
	PageSize int
}

// GetPage : returns the page number.
func (p PaginationParams) GetPage() int {
	return p.Page
}

// GetPageSize : returns the page size.
func (p PaginationParams) GetPageSize() int {
	return p.PageSize
}

// DefaultPageSize is the default number of items per page.
const DefaultPageSize = 10

// NewPaginationParams : creates a new PaginationParams with default values.
func NewPaginationParams(c *fiber.Ctx) PaginationParams {
	page := c.QueryInt("page", 1)

	if page < 1 {
		page = 1
	}

	pageSize := c.QueryInt("pageSize", DefaultPageSize)

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// CalculateOffset : calculates the offset for the SQL query based on pagination parameters.
func (p PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetPaginationResponse : returns the pagination response.
func (p PaginationParams) GetPaginationResponse(req *fasthttp.Request, total int64) paginationResponse {
	return paginationResponse{
		Next:  p.GetNext(req.URI(), total),
		Prev:  p.GetPrev(req.URI()),
		Total: total,
		Page:  p.Page,
	}
}

// GetPrev : returns the previous page URL.
func (p PaginationParams) GetPrev(baseURL *fasthttp.URI) string {
	if p.Page <= 1 {
		return ""
	}

	var previousURL fasthttp.URI
	baseURL.CopyTo(&previousURL)

	previousURL.QueryArgs().Set("page", strconv.Itoa(p.Page-1))
	previousURL.QueryArgs().Set("pageSize", strconv.Itoa(p.PageSize))

	return previousURL.String()
}

// GetNext : returns the next page URL.
func (p PaginationParams) GetNext(baseURL *fasthttp.URI, total int64) string {
	if total <= int64(p.Page*p.PageSize) {
		return ""
	}

	var previousURL fasthttp.URI
	baseURL.CopyTo(&previousURL)

	previousURL.QueryArgs().Set("page", strconv.Itoa(p.Page+1))
	previousURL.QueryArgs().Set("pageSize", strconv.Itoa(p.PageSize))

	return previousURL.String()
}

// paginationResponse : contains the pagination response.
type paginationResponse struct {
	Next  string `json:"next"`
	Prev  string `json:"prev"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
}
