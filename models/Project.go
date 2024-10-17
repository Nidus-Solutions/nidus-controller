package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	ID       string    `gorm:"primaryKey; not null; unique_index"`
	UserId   string    `gorm:"not null;foreignKey=ID"`
	Name     string    `gorm:"type=varchar(255); not null;"`
	Company  string    `gorm:"type=varchar(255); not null;"`
	CNPJ     string    `gorm:"type=varchar(14); not null; unique_index;"`
	Value    float64   `gorm:"not null;"`
	Deadline time.Time `gorm:"type=timestamp; not null;"`
}

func NewProject() *Project {
	Project := Project{
		ID: uuid.NewString(),
	}

	return &Project
}
