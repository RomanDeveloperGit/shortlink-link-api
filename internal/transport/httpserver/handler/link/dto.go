package link

import "github.com/RomanDeveloperGit/shortlink-link-api/internal/model"

type CreateLinkRequest struct {
	FullURL string `json:"fullUrl" validate:"required,url"`
	TTLDays int    `json:"ttlDays" validate:"gt=0,lte=30"`
}

type linkResponse struct {
	model.Link

	IsExpired bool `json:"isExpired"`
}

func NewLinkResponse(link *model.Link) *linkResponse {
	return &linkResponse{
		Link:      *link,
		IsExpired: link.IsExpired(),
	}
}
