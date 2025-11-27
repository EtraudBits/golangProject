package db

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite" //importa o GORM
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model //inclui as estruturas da struct model do pacote gorm
	Name string
	CPF int
	Email string
	Age int
	Active bool

}


func Init() *gorm.DB { //função Publica
// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("student.db"), &gorm.Config{}) //a forma que se cria o GORM usando o Banco de Dados SQLite
	if err != nil { //trata o erro, caso não consiga executar
		log.Fatalln(err)
	}

	db.AutoMigrate(&Student{}) //aponta para a struct Student (Gerenciar a estrutura Student)

	return db
}

func AddStudent () { //função publica
	db := Init()

	student := Student {
		Name: "Duarte",
		CPF: 12345678911,
		Email: "du84arte@gamil.com",
		Age: 41,
		Active: true,
	}

	if result := db.Create(&student); result.Error != nil {
		fmt.Println ("Erro to create student")
	}

	fmt.Println("Create student!")

}