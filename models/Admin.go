package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model

	ID      string `gorm:"primaryKey; not null; unique_index;"`
	Name    string `gorm:"not null; unique_index;"`
	IsAdmin bool   `gorm:"not null; default:true;"`
}

func NewAdmin() *Admin {
	Admin := Admin{
		ID: uuid.NewString(),
	}

	return &Admin
}
