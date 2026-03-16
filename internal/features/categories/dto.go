package category

type CategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type CategoryResponse struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
