package db

import (
	"github.com/rs/zerolog/log" //formatar os logs (deixar mais organizado)

	//"gorm.io/driver/sqlite" // troque por: "github.com/glebarez/sqlite" abaixo
	"github.com/glebarez/sqlite" //importa o GORM
	"gorm.io/gorm"
)

//cria uma struct tipo Studenthandler
type StudentHandler struct { //criação para poder atrelar as funções em uma estrutura especifica
	DB *gorm.DB //estrutura DB apontando para o gorm.DB
}
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
		log.Fatal().Err(err).Msgf("Failed to initialize SQLite: %s", err.Error()) //usando log.Fatal para sinalizar o erro
	}

	db.AutoMigrate(&Student{}) //aponta para a struct Student (Gerenciar a estrutura Student)

	return db
}

func NewStudentHandler (db *gorm.DB) *StudentHandler {
	return &StudentHandler{DB: db}
}

func (s *StudentHandler) AddStudent (student Student) error { //função publica - pode usar fora do pacote db (inicia com a primeira letra MAIUCULA) passando o strudent como pararamento e retornar um erro.
	
	
	if result := s.DB.Create(&student); result.Error != nil { //temos que o usar o & (é comercial) para chamar a variavel
		log.Error().Msg("Failed to Create Student!")
		return result.Error
	}
	log.Info().Msg("Create Student!")
	return nil  //caso não tenha ocorrido nenhum erro, passa a exibir a mensagem

}

func (s *StudentHandler) GetStudents () ([]Student, error) { //função publica -> listar usuario , retornando a lista de usuarios (estudantes) + erro.
students := []Student{} // retorna uma lista de usuarios (estudantes)



err := s.DB.Find(&students).Error//consultar essa tabela usando o metodo find do GORM
return students, err
}

