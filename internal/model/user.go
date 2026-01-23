package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Dismiss   int       `json:"dismiss"`
}
