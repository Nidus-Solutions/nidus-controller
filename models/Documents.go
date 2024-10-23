// Model para o documento
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model

	ID        string `gorm:"primaryKey; not null; unique_index;" json:"id"`
	ProjectID string `gorm:"not null" json:"projectId"`
	Name      string `gorm:"type=varchar(255); not null;" json:"name"`
	Link      string `gorm:"type=varchar(255); not null;" json:"link"`
}

// Gerando um novo ID para o documento
func NewDocument() *Document {
	Document := Document{
		ID: uuid.NewString(),
	}

	return &Document
}
