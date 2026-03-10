package category

import "time"

func ToCategoryResponse(c Category) CategoryResponse {
	var updatedAt string
	if c.UpdatedAt != nil {
		updatedAt = c.UpdatedAt.Format(time.RFC3339)
	}

	var createdAt string

	if c.CreatedAt != nil {
		createdAt = c.CreatedAt.Format(time.RFC3339)
	}

	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Type:      c.Type,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
