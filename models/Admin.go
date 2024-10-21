package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model

	ID       string `gorm:"primaryKey; not null; unique_index;" json:"id"`
	Name     string `gorm:"not null; unique_index;" json:"name"`
	Username string `gorm:"not null; unique_index;" json:"username"`
	Password string `gorm:"not null;" json:"password"`
	IsAdmin  bool   `gorm:"not null; default:true;" json:"is_admin"` // talvez eu nem use
}

func NewAdmin() *Admin {
	Admin := Admin{
		ID: uuid.NewString(),
	}

	return &Admin
}
