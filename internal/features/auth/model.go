package auth

import (
	"time"
)

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type UserProvider struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	ProviderName   string    `db:"provider_name"`
	ProviderUserID string    `db:"provider_user_id"`
	CreatedAt      time.Time `db:"created_at"`
}
