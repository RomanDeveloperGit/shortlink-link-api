package link

import (
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	"github.com/google/uuid"
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
	Meta      meta  `json:"meta"`
	LinkId    int   `json:"linkId"`
	Timestamp int64 `json:"timestamp"`
}

func newLinkVisitedMessage(linkId int, timestamp int64) *linkVisitedMessage {
	return &linkVisitedMessage{
		Meta:      meta{Event: linkVisitedEvent, EventId: uuid.New().String()},
		LinkId:    linkId,
		Timestamp: timestamp,
	}
}
