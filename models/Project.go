package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model

	ID       string  `gorm:"primaryKey; not null; unique_index;" json:"id"`
	UserID   string  `gorm:"not null" json:"userId"`
	Name     string  `gorm:"type=varchar(255); not null;" json:"name"`
	Company  string  `gorm:"type=varchar(255); not null;" json:"company"`
	CNPJ     string  `gorm:"type=varchar(18); not null; unique_index;" json:"cnpj"`
	Value    float64 `gorm:"not null;" json:"value"`
	Deadline string  `gorm:"type=varchar(255); not null;" json:"deadline"`

	// Aqui é feito o relacionamento entre o projeto e os documentos, do tipo 1 para muitos
	Documents []Document `gorm:"foreignKey:ProjectID" json:"documents"`
}

// Gerando o ID do projeto
func NewProject() *Project {
	Project := Project{
		ID: uuid.NewString(),
	}

	return &Project
}
