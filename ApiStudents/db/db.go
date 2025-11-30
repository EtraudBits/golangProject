package db

import (
	"fmt"
	"log"

	//"gorm.io/driver/sqlite" // troque por: "github.com/glebarez/sqlite" abaixo
	"github.com/glebarez/sqlite" //importa o GORM
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


func Init() *gorm.DB { //função Publica
// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("./student.db"), &gorm.Config{}) //a forma que se cria o GORM usando o Banco de Dados SQLite
	if err != nil { //trata o erro, caso não consiga executar
		log.Fatalln(err) // se der algum erro ele para a aplicação (log.fatal)
	}

	db.AutoMigrate(&Student{}) //aponta para a struct Student (Gerenciar a estrutura Student)

	return db
}

func AddStudent (student Student) error { //função publica - pode usar fora do pacote db (inicia com a primeira letra MAIUCULA) passando o strudent como pararamento e retornar um erro.
	db := Init() //não chama "db.Init" por já estamos dentro do pacote db. só inclui caso estejamos fora do pacote.
	
	if result := db.Create(&student); result.Error != nil { //temos que o usar o & (é comercial) para chamar a variavel
		return result.Error
	}

	fmt.Println("Create student!")
	return nil  //caso não tenha ocorrido nenhum erro, passa a exibir a mensagem

}