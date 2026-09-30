package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/apperror"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(
	ctx context.Context,
	shortCode string,
	fullURL string,
	expiresAt time.Time,
) (*model.Link, error) {
	var link model.Link

	const q = `
		INSERT INTO links (short_code, full_url, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, short_code, full_url, expires_at, created_at, updated_at
	`

	err := r.db.GetContext(
		ctx,
		&link,
		q,
		shortCode,
		fullURL,
		expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create link: %w", err)
	}

	return &link, nil
}

func (r *repository) GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	var link model.Link

	q := `
		SELECT * FROM links
		WHERE short_code = $1
	`

	err := r.db.GetContext(ctx, &link, q, shortCode)
	if err != nil {
		baseError := fmt.Errorf("failed to get link by short code: %w", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(
				baseError,
				fmt.Errorf("typed app error: %w", apperror.ErrLinkNotFound),
			)
		}

		return nil, baseError
	}

	return &link, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.Link, error) {
	var link model.Link

	q := `
		SELECT * FROM links
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &link, q, id)
	if err != nil {
		baseError := fmt.Errorf("failed to get link by id: %w", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Join(
				baseError,
				fmt.Errorf("typed app error: %w", apperror.ErrLinkNotFound),
			)
		}

		return nil, baseError
	}

	return &link, nil
}
