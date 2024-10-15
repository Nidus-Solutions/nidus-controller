package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Mail      string    `json:"mail"`
	CPF       string    `json:"cpf"`
	CreatedAt time.Time `json:"createAt"`
}

func NewUser() *User {
	User := User{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
	}

	return &User
}
