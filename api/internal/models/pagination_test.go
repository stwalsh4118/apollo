package models_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sean/apollo/api/internal/models"
)

func TestParsePaginationDefaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	p := models.ParsePagination(req)

	if p.Page != 1 {
		t.Fatalf("expected default page 1, got %d", p.Page)
	}

	if p.PerPage != 20 {
		t.Fatalf("expected default per_page 20, got %d", p.PerPage)
	}
}

func TestParsePaginationCustomValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/items?page=3&per_page=10", nil)
	p := models.ParsePagination(req)

	if p.Page != 3 {
		t.Fatalf("expected page 3, got %d", p.Page)
	}

	if p.PerPage != 10 {
		t.Fatalf("expected per_page 10, got %d", p.PerPage)
	}
}

func TestParsePaginationClampsMaxPerPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/items?per_page=500", nil)
	p := models.ParsePagination(req)

	if p.PerPage != 100 {
		t.Fatalf("expected clamped per_page 100, got %d", p.PerPage)
	}
}

func TestParsePaginationHandlesInvalidValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/items?page=abc&per_page=-1", nil)
	p := models.ParsePagination(req)

	if p.Page != 1 {
		t.Fatalf("expected fallback page 1 for invalid input, got %d", p.Page)
	}

	if p.PerPage != 20 {
		t.Fatalf("expected fallback per_page 20 for negative input, got %d", p.PerPage)
	}
}

func TestPaginationParamsOffset(t *testing.T) {
	tests := []struct {
		page    int
		perPage int
		want    int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 10, 20},
		{1, 50, 0},
	}

	for _, tc := range tests {
		p := models.PaginationParams{Page: tc.page, PerPage: tc.perPage}
		got := p.Offset()

		if got != tc.want {
			t.Errorf("Offset() for page=%d per_page=%d: got %d, want %d", tc.page, tc.perPage, got, tc.want)
		}
	}
}
