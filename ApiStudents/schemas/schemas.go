package schemas

import (
	"time"

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

type StudentResponse struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UPdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
	Name string `json:"name"` //temos que dizer como é nome do campo vai querer atrelar no json que recebe no post
	CPF int `json:"cpf"`
	Email string `json:"email"`
	Age int `json:"age"`
	Active bool `json:"active"`
}

//Cria uma nova resposta para o usuario
func NewResponse(students []Student) []StudentResponse { //recebe os studants que retorna uma nova resposta (lista)
	//percorre a lista de estudantes, criando uma lista nova de acordo com a nova estrutura (StudentResponse)
	studentResponses := []StudentResponse{}

	for _, student := range students {

		studentResponse := StudentResponse {

			ID: int(student.ID),
			CreatedAt: student.CreatedAt,
			UPdatedAt: student.UpdatedAt,
			//DeletedAt: student.DeletedAt, -> não dquero que ele apareça então tiramos a informação de Delete
			Name: student.Name,
			Email: student.Email,
			CPF: student.CPF,
			Age: student.Age,
			Active: student.Active,
		}
		studentResponses = append(studentResponses, studentResponse)
	}
	return studentResponses
}