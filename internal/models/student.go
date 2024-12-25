package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Course    string    `json:"course"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
