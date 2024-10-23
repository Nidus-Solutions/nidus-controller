package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID       string `gorm:"primaryKey; not null; unique_index;" json:"id"`
	Name     string `gorm:"type=varchar(255); not null;" json:"name"`
	Phone    string `gorm:"type=varchar(25); not null; unique_index;" json:"phone"`
	Mail     string `gorm:"type=varchar(35); not null; unique_index;" json:"mail"`
	CPF      string `gorm:"type=varchar(11); not null; unique_index;" json:"cpf"`
	UserName string `gorm:"type=varchar(25); not null; unique_index;" json:"username"`
	Password string `gorm:"type=varchar(25); not null;" json:"password"`

	// Aqui é feito o relacionamento entre o usuário e os projetos, do tipo 1 para muitos
	Projects []Project `gorm:"foreignKey:UserID" json:"projects"`
}

func NewUser() *User {
	User := User{
		ID: uuid.NewString(),
	}

	return &User
}
