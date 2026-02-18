package models

import (
	"net/http"
	"strconv"
)

const (
	defaultPage    = 1
	defaultPerPage = 20
	maxPerPage     = 100
)

// PaginationParams holds pagination settings parsed from query parameters.
type PaginationParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// Offset returns the SQL OFFSET value for the current page.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// PaginatedResponse wraps a list of items with pagination metadata.
type PaginatedResponse[T any] struct {
	Items   []T `json:"items"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// ParsePagination extracts page and per_page from an HTTP request's query params.
func ParsePagination(r *http.Request) PaginationParams {
	page := parseIntQuery(r, "page", defaultPage)
	perPage := parseIntQuery(r, "per_page", defaultPerPage)

	if page < 1 {
		page = defaultPage
	}

	if perPage < 1 {
		perPage = defaultPerPage
	}

	if perPage > maxPerPage {
		perPage = maxPerPage
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
	}
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return val
}
