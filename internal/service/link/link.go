package link

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/apperror"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/pkg/observability"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link/shortcode"
)

const linksEventsTopic = "links.events"

type LinkRepository interface {
	Create(
		ctx context.Context,
		shortCode, fullURL string,
		authorEmail *string,
		expiresAt time.Time,
	) (*model.Link, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error)
	GetByID(ctx context.Context, id int) (*model.Link, error)
}

type service struct {
	linkRepository                 LinkRepository
	producer                       sarama.SyncProducer
	cache                          *redis.Client
	logger                         *slog.Logger
	shortCodeLength                int
	attemptsGenerateShortLinkLimit int
}

func NewService(
	linkRepository LinkRepository,
	producer sarama.SyncProducer,
	cache *redis.Client,
	logger *slog.Logger,
	shortCodeLength int,
	attemptsGenerateShortLinkLimit int,
) *service {
	return &service{
		linkRepository:                 linkRepository,
		producer:                       producer,
		cache:                          cache,
		logger:                         logger,
		shortCodeLength:                shortCodeLength,
		attemptsGenerateShortLinkLimit: attemptsGenerateShortLinkLimit,
	}
}

func (s *service) buildCacheKey(shortCode string) string {
	// ключом делаем именно шорткод, потому что важно обеспечить быструю работу метода Visit
	return fmt.Sprintf("link:%s", shortCode)
}

func (s *service) saveLinkInCache(ctx context.Context, link *model.Link) error {
	cacheKey := s.buildCacheKey(link.ShortCode)

	pipe := s.cache.Pipeline()
	pipe.HSet(
		ctx,
		cacheKey,
		link,
	)
	pipe.Expire(ctx, cacheKey, time.Until(link.ExpiresAt))

	_, err := pipe.Exec(ctx)
	if err != nil {
		observability.NewObservedLogger(ctx, s.logger).
			Error("failed to save link in redis cache", slog.Any("error", err))

		// намеренно не выходим из функции, потому что это не должно ломать процесс.
		// пусть лучше серис работает с задержками, чем не работает совсем
	}

	return err
}

// TODO: внедрить transaction outbox
// Сейчас есть проблема - плодим лишние сущности в бд, если брокер упадет. Обеспечиваем целостность конечного результата, но допускаем "фантомы" в БД, которые никто не будет юзать. Либо могут начать юзать - потенциальный баг: POST вернул error, в БД сущность есть, чувак брутфорсом подбирает шорткод, вызывает Visit метод, он отсылает в кафку сообщение, а консьюмер не понимает, как произошел переход по ссылке, которой для него нет (ранее при создании событие же не отправилось)
func (s *service) Create(
	ctx context.Context,
	fullURL string,
	ttlDays int,
	authorEmail *string,
) (*model.Link, error) {
	var attemptsCount int

	var link *model.Link
	var repoErr error

	shortCode, err := shortcode.GenerateShortCode(s.shortCodeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, ttlDays)

	for attemptsCount = 0; attemptsCount < s.attemptsGenerateShortLinkLimit && link == nil; attemptsCount++ {
		link, repoErr = s.linkRepository.Create(ctx, shortCode, fullURL, authorEmail, expiresAt)
	}

	if link == nil || repoErr != nil {
		return nil, fmt.Errorf(
			"failed to create link after %d attempts: %w",
			attemptsCount,
			repoErr,
		)
	}

	msg, err := json.Marshal(newLinkCreatedMessage(link))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal link: %w", err)
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: linksEventsTopic,
		Key: sarama.StringEncoder(
			shortCode, // чтобы событие Visit не было прочитано раньше, чем Create
		),
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message to kafka: %w", err)
	}

	_ = s.saveLinkInCache(ctx, link)

	return link, nil
}

func (s *service) GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	cachedLink := &model.Link{}
	err := s.cache.HGetAll(ctx, s.buildCacheKey(shortCode)).Scan(cachedLink)

	if err == nil && cachedLink.ShortCode != "" {
		return cachedLink, nil
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		observability.NewObservedLogger(ctx, s.logger).
			Error("failed to get link from redis cache", slog.Any("error", err))
	}

	link, err := s.linkRepository.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	_ = s.saveLinkInCache(ctx, link)

	return link, nil
}

// TODO: внедрить transaction outbox
// Сейчас есть проблема - не даем перейти по ссылке, если брокер упадет. Страдает перфоманс. Операция редиректа должна быть быстрая
func (s *service) Visit(ctx context.Context, shortCode string) (*model.Link, error) {
	link, err := s.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	isExpired := link.IsExpired()

	if isExpired {
		return nil, fmt.Errorf("link is expired: %w", apperror.ErrLinkExpired)
	}

	msg, err := json.Marshal(newLinkVisitedMessage(link.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal link: %w", err)
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: linksEventsTopic,
		Key: sarama.StringEncoder(
			shortCode, // чтобы событие Visit не было прочитано раньше, чем Create
		),
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return link, nil
}

func (s *service) GetByID(ctx context.Context, id int) (*model.Link, error) {
	return s.linkRepository.GetByID(ctx, id)
}
