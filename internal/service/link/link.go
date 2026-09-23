package link

import (
	"context"
	"fmt"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link/shortcode"
)

type LinkRepository interface {
	Create(ctx context.Context, shortCode string, fullURL string, expiresAt time.Time) (*model.Link, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error)
	GetByID(ctx context.Context, id int) (*model.Link, error)
}

type service struct {
	shortCodeLength                int
	attemptsGenerateShortLinkLimit int
	linkRepository                 LinkRepository
}

func NewService(
	linkRepository LinkRepository,
	shortCodeLength int,
	attemptsGenerateShortLinkLimit int,
) *service {
	return &service{
		linkRepository:                 linkRepository,
		shortCodeLength:                shortCodeLength,
		attemptsGenerateShortLinkLimit: attemptsGenerateShortLinkLimit,
	}
}

func (s *service) Create(ctx context.Context, fullURL string, ttlDays int) (*model.Link, error) {
	var attempts int
	var link *model.Link
	var err error

	shortCode := shortcode.GenerateShortCode(s.shortCodeLength)
	expiresAt := time.Now().UTC().AddDate(0, 0, ttlDays)

	for attempts = 0; attempts < s.attemptsGenerateShortLinkLimit && link == nil; attempts++ {
		link, err = s.linkRepository.Create(ctx, shortCode, fullURL, expiresAt)
	}

	if link == nil || err != nil {
		return nil, fmt.Errorf("failed to create link after %d attempts: %w", attempts, err)
	}

	return link, nil
}

func (s *service) GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	return s.linkRepository.GetByShortCode(ctx, shortCode)
}

func (s *service) GetByID(ctx context.Context, id int) (*model.Link, error) {
	return s.linkRepository.GetByID(ctx, id)
}
