package transaction

type CreateTransaction struct {
	Title     string    `json:"title" binding:"required"`
	Amount    float64   `json:"amount" binding:"required"`
	Type      Type      `json:"type" binding:"required"`
	Category  Category  `json:"category"`
	Frequency Frequency `json:"frequency"`
}

type TransactionResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	Type      Type      `json:"type"`
	Category  Category  `json:"category"`
	Frequency Frequency `json:"frequency"`
	CreatedAt string    `json:"created_at"`
}

type UpdateTransaction struct {
	Title string `json:"date"`
}
