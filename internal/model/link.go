package model

import "time"

type Link struct {
	ID          int       `db:"id"           redis:"id"           json:"id"`
	ShortCode   string    `db:"short_code"   redis:"short_code"   json:"shortCode"`
	FullURL     string    `db:"full_url"     redis:"full_url"     json:"fullUrl"`
	AuthorEmail *string   `db:"author_email" redis:"author_email" json:"authorEmail"`
	ExpiresAt   time.Time `db:"expires_at"   redis:"expires_at"   json:"expiresAt"`
	CreatedAt   time.Time `db:"created_at"   redis:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   redis:"updated_at"   json:"updatedAt"`
}

func (l *Link) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}
