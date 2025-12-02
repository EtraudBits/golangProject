package schemas

import (
	"gorm.io/gorm"
)

// struct do tipo Student
type Student struct {
	gorm.Model //inclui as estruturas da struct model do pacote gorm
	Name string `json:"Name"` //temos que dizer como é nome do campo vai querer atrelar no json que recebe no post
	CPF int `json:"CPF"`
	Email string `json:"Email"`
	Age int `json:"Age"`
	Active bool `json:"Active"`

}
