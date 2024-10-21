package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	ID        string  `gorm:"primaryKey; not null; unique_index;" json:"id"`
	UserId    string  `gorm:"not null;foreignKey=ID" json:"userId"`
	Name      string  `gorm:"type=varchar(255); not null;" json:"name"`
	Company   string  `gorm:"type=varchar(255); not null;" json:"company"`
	Documents string  `gorm:"type=string; not null;" json:"documents"`
	CNPJ      string  `gorm:"type=varchar(18); not null; unique_index;" json:"cnpj"`
	Value     float64 `gorm:"not null;" json:"value"`
	Deadline  string  `gorm:"type=string; not null;"`
}

func NewProject() *Project {
	Project := Project{
		ID: uuid.NewString(),
	}

	return &Project
}
