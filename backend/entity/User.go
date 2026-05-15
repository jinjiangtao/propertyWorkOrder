package entity

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

const (
	RoleUser      = 1
	RoleAdmin     = 2
)