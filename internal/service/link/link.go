package link

import (
	"context"
	"log/slog"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/randomstring"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/requestid"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/traceid"
)

type LinkRepository interface {
	Create(ctx context.Context, shortCode, fullURL string, expiresAt time.Time) (*model.Link, error)
}

type service struct {
	shortLinkLength                int
	attemptsGenerateShortLinkLimit int
	linkRepository                 LinkRepository
	logger                         *slog.Logger
}

func NewService(linkRepository LinkRepository, logger *slog.Logger, shortLinkLength int, attemptsGenerateShortLinkLimit int) *service {
	return &service{
		linkRepository:                 linkRepository,
		logger:                         logger,
		shortLinkLength:                shortLinkLength,
		attemptsGenerateShortLinkLimit: attemptsGenerateShortLinkLimit,
	}
}

func (s *service) Create(ctx context.Context, fullURL string, ttlDays int) (*model.Link, error) {
	var attempts int
	var link *model.Link
	var err error

	shortCode := randomstring.NewRandomString(s.shortLinkLength)
	expiresAt := time.Now().UTC().AddDate(0, 0, int(ttlDays))

	for attempts = 0; attempts < s.attemptsGenerateShortLinkLimit && link == nil; attempts++ {
		link, err = s.linkRepository.Create(ctx, shortCode, fullURL, expiresAt)
	}

	if link == nil {
		s.logger.Error("failed to create link",
			slog.String("short_code", shortCode),
			slog.String("request_id", requestid.RequestIDFromContext(ctx)),
			slog.String("trace_id", traceid.TraceIDFromContext(ctx)),
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	s.logger.Info("link created",
		slog.String("short_code", shortCode),
		slog.Int("attempts", attempts),
		slog.String("request_id", requestid.RequestIDFromContext(ctx)),
		slog.String("trace_id", traceid.TraceIDFromContext(ctx)),
	)

	return link, nil
}
