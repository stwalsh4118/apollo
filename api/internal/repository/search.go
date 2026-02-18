package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/sean/apollo/api/internal/models"
)

// ErrInvalidQuery indicates a malformed FTS5 search query.
var ErrInvalidQuery = errors.New("invalid search query")

// SearchRepository defines search operations.
type SearchRepository interface {
	Search(ctx context.Context, query string, params models.PaginationParams) (*models.PaginatedResponse[models.SearchResult], error)
}

// SQLiteSearchRepository implements SearchRepository using FTS5.
type SQLiteSearchRepository struct {
	db *sql.DB
}

// NewSearchRepository creates a new SQLiteSearchRepository.
func NewSearchRepository(db *sql.DB) *SQLiteSearchRepository {
	return &SQLiteSearchRepository{db: db}
}

// searchSQL: snippet column index 3 = body column in search_index(entity_type, entity_id, title, body).
const searchSQL = `
SELECT entity_type, entity_id, title,
       snippet(search_index, 3, '<mark>', '</mark>', '...', 30)
FROM search_index
WHERE search_index MATCH ?
ORDER BY rank
LIMIT ? OFFSET ?
`

const countSearchSQL = `SELECT COUNT(*) FROM search_index WHERE search_index MATCH ?`

func (r *SQLiteSearchRepository) Search(ctx context.Context, query string, params models.PaginationParams) (*models.PaginatedResponse[models.SearchResult], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, countSearchSQL, query).Scan(&total); err != nil {
		return nil, classifySearchError(err)
	}

	rows, err := r.db.QueryContext(ctx, searchSQL, query, params.PerPage, params.Offset())
	if err != nil {
		return nil, classifySearchError(err)
	}
	defer rows.Close()

	var results []models.SearchResult

	for rows.Next() {
		var sr models.SearchResult
		if err := rows.Scan(&sr.EntityType, &sr.EntityID, &sr.Title, &sr.Snippet); err != nil {
			return nil, fmt.Errorf("scan search result: %w", err)
		}

		results = append(results, sr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate search results: %w", err)
	}

	if results == nil {
		results = []models.SearchResult{}
	}

	return &models.PaginatedResponse[models.SearchResult]{
		Items:   results,
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	}, nil
}

func classifySearchError(err error) error {
	if strings.Contains(err.Error(), "fts5: syntax error") {
		return fmt.Errorf("%w: %s", ErrInvalidQuery, err.Error())
	}

	return fmt.Errorf("search: %w", err)
}

// Verify interface compliance at compile time.
var _ SearchRepository = (*SQLiteSearchRepository)(nil)
