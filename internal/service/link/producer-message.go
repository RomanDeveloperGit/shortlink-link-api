package link

import (
	"time"

	"github.com/google/uuid"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
)

const (
	linkCreatedEvent string = "link.created"
	linkVisitedEvent string = "link.visited"
)

type meta struct {
	Event   string `json:"event"`
	EventId string `json:"eventId"`
}

type linkCreatedMessage struct {
	Meta meta       `json:"meta"`
	Link model.Link `json:"link"`
}

func newLinkCreatedMessage(link *model.Link) *linkCreatedMessage {
	return &linkCreatedMessage{
		Meta: meta{
			Event:   linkCreatedEvent,
			EventId: uuid.New().String(),
		},
		Link: *link,
	}
}

type linkVisitedMessage struct {
	Meta      meta      `json:"meta"`
	LinkID    int       `json:"linkId"`
	Timestamp time.Time `json:"timestamp"`
}

func newLinkVisitedMessage(linkID int) *linkVisitedMessage {
	return &linkVisitedMessage{
		Meta: meta{
			Event:   linkVisitedEvent,
			EventId: uuid.New().String(),
		},
		LinkID:    linkID,
		Timestamp: time.Now(),
	}
}
