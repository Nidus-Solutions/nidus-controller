package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID       string `gorm:"primaryKey; not null; unique_index;"`
	Name     string `gorm:"type=varchar(255); not null;"`
	Phone    string `gorm:"type=varchar(25); not null; unique_index;"`
	Mail     string `gorm:"type=varchar(35); not null; unique_index;"`
	CPF      string `gorm:"type=varchar(11); not null; unique_index;"`
	UserName string `gorm:"type=varchar(25); not null; unique_index;"`
	Password string `gorm:"type=varchar(25); not null;"`
}

func NewUser() *User {
	User := User{
		ID: uuid.NewString(),
	}

	return &User
}
