package link

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link/shortcode"
)

const linksEventsTopic = "links.events"

type LinkRepository interface {
	Create(ctx context.Context, shortCode, fullURL string, expiresAt time.Time) (*model.Link, error)
	GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error)
	GetByID(ctx context.Context, id int) (*model.Link, error)
}

type service struct {
	linkRepository                 LinkRepository
	producer                       sarama.SyncProducer
	shortCodeLength                int
	attemptsGenerateShortLinkLimit int
}

func NewService(
	linkRepository LinkRepository,
	producer sarama.SyncProducer,
	shortCodeLength int,
	attemptsGenerateShortLinkLimit int,
) *service {
	return &service{
		linkRepository:                 linkRepository,
		producer:                       producer,
		shortCodeLength:                shortCodeLength,
		attemptsGenerateShortLinkLimit: attemptsGenerateShortLinkLimit,
	}
}

// TODO: внедрить transaction outbox
// Сейчас есть проблема - плодим лишние сущности в бд, если брокер упадет. Обеспечиваем целостность конечного результата, но допускаем "фантомы" в БД, которые никто не будет юзать. Либо могут начать юзать - потенциальный баг: POST вернул error, в БД сущность есть, чувак брутфорсом подбирает шорткод, вызывает Visit метод, он отсылает в кафку сообщение, а консьюмер не понимает, как произошел переход по ссылке, которой для него нет (ранее при создании событие же не отправилось)
func (s *service) Create(ctx context.Context, fullURL string, ttlDays int) (*model.Link, error) {
	var attempts int
	var link *model.Link
	var repoErr error

	shortCode, err := shortcode.GenerateShortCode(s.shortCodeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, ttlDays)

	for attempts = 0; attempts < s.attemptsGenerateShortLinkLimit && link == nil; attempts++ {
		link, repoErr = s.linkRepository.Create(ctx, shortCode, fullURL, expiresAt)
	}

	if link == nil || repoErr != nil {
		return nil, fmt.Errorf("failed to create link after %d attempts: %w", attempts, repoErr)
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

	return link, nil
}

func (s *service) GetByShortCode(ctx context.Context, shortCode string) (*model.Link, error) {
	return s.linkRepository.GetByShortCode(ctx, shortCode)
}

// TODO: внедрить transaction outbox
// Сейчас есть проблема - не даем перейти по ссылке, если брокер упадет. Страдает перфоманс. Операция редиректа должна быть быстрая
func (s *service) Visit(ctx context.Context, shortCode string) (*model.Link, error) {
	link, err := s.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	msg, err := json.Marshal(newLinkVisitedMessage(link.ID, time.Now().UnixMilli()))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal link: %w", err)
	}

	// Всю сущность отдавать? В статистике же она уже есть, поэтому мб только ID + timestamp? для статы
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
