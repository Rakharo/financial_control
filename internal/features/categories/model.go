package category

import "time"

type Category struct {
	ID        uint64
	UserID    *uint64
	Name      string
	Type      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
