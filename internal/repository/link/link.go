package link

import (
	"context"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, shortCode, fullURL string, expiresAt time.Time) (*model.Link, error) {
	var link model.Link

	const q = `
		INSERT INTO links (short_code, full_url, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, short_code, full_url, expires_at, created_at, updated_at
	`

	err := r.db.GetContext(
		ctx,
		link,
		q,
		shortCode,
		fullURL,
		expiresAt,
	)

	return &link, err
}
