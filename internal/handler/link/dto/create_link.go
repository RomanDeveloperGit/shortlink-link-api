package dto

import "github.com/RomanDeveloperGit/shortlink-link-api/internal/model"

type CreateLinkRequest struct {
	FullURL string `json:"fullUrl" validate:"required,url"`
	TTLDays int    `json:"ttlDays" validate:"gt=0,lte=30"`
}

type createLinkResponse struct {
	model.Link

	IsExpired bool `json:"isExpired"`
}

func NewCreateLinkResponse(link *model.Link) createLinkResponse {
	return createLinkResponse{
		Link:      *link,
		IsExpired: link.IsExpired(),
	}
}
