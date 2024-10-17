package entities

import (
	"github.com/google/uuid"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Mail  string `json:"mail"`
	CPF   string `json:"cpf"`
}

func NewUser() *User {
	User := User{
		ID: uuid.NewString(),
	}

	return &User
}
