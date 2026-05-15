package entity

type WorkOrder struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Images      string `json:"images"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

const (
	StatusPending   = 1
	StatusProcessing = 2
	StatusCompleted = 3
)