package link

import "github.com/RomanDeveloperGit/shortlink-link-api/internal/model"

type createLinkRequest struct {
	FullURL     string  `json:"fullUrl"     validate:"required,url"`
	TTLDays     int     `json:"ttlDays"     validate:"gt=0,lte=30"`
	AuthorEmail *string `json:"authorEmail" validate:"omitempty,email"`
}

type linkResponse struct {
	model.Link

	IsExpired bool `json:"isExpired"`
}

func newLinkResponse(link *model.Link) *linkResponse {
	return &linkResponse{
		Link:      *link,
		IsExpired: link.IsExpired(),
	}
}
