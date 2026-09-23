package model

import "time"

type Link struct {
	ID        int       `db:"id"         json:"id"`
	ShortCode string    `db:"short_code" json:"shortCode"`
	FullURL   string    `db:"full_url"   json:"fullUrl"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func (l *Link) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}
