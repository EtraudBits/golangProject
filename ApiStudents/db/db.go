package db

import (
	"github.com/rs/zerolog/log" //formatar os logs (deixar mais organizado)

	//"gorm.io/driver/sqlite" // troque por: "github.com/glebarez/sqlite" abaixo
	"github.com/EtraudBits/golangProject/ApiStudents/schemas" //pacote do repositório schemas
	"github.com/glebarez/sqlite"                              //importa o GORM
	"gorm.io/gorm"
)

//cria uma struct tipo Studenthandler
type StudentHandler struct { //criação para poder atrelar as funções em uma estrutura especifica
	DB *gorm.DB //estrutura DB apontando para o gorm.DB
}


func Init() *gorm.DB { //função Publica
// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("./student.db"), &gorm.Config{}) //a forma que se cria o GORM usando o Banco de Dados SQLite
	if err != nil { //trata o erro, caso não consiga executar
		log.Fatal().Err(err).Msgf("Failed to initialize SQLite: %s", err.Error()) //usando log.Fatal para sinalizar o erro
	}

	db.AutoMigrate(&schemas.Student{}) //aponta para a struct Student (Gerenciar a estrutura Student)

	return db
}

func NewStudentHandler (db *gorm.DB) *StudentHandler {
	return &StudentHandler{DB: db}
}

func (s *StudentHandler) AddStudent (student schemas.Student) error { //função publica - pode usar fora do pacote db (inicia com a primeira letra MAIUCULA) passando o strudent como pararamento e retornar um erro.
	
	
	if result := s.DB.Create(&student); result.Error != nil { //temos que o usar o & (é comercial) para chamar a variavel
		log.Error().Msg("Failed to Create Student!")
		return result.Error
	}
	log.Info().Msg("Create Student!")
	return nil  //caso não tenha ocorrido nenhum erro, passa a exibir a mensagem

}

func (s *StudentHandler) GetStudents () ([]schemas.Student, error) { //função publica -> listar usuario , retornando a lista de usuarios (estudantes) + erro.
students := []schemas.Student{} // retorna uma lista de usuarios (estudantes)
err := s.DB.Find(&students).Error//consultar essa tabela usando o metodo find do GORM
return students, err
}

//metodo para filtrar o campo active
func (s *StudentHandler) GetFilteredStudent (active bool) ([]schemas.Student, error) { //função publica -> listar usuario , retornando a lista de usuarios (estudantes) + erro.
filteredStudents := []schemas.Student{} // retorna uma lista de usuarios (estudantes)
err := s.DB.Where("active = ?", active).Find(&filteredStudents)
return filteredStudents, err.Error
}


func (s *StudentHandler) GetStudent (id int) (schemas.Student, error) { //função publica -> Busca um unico usuario , retornando um usuarios (estudante) + erro. buscando pelo o ID
var student schemas.Student // variavel student do tipo Stundet (guarda um usuarios (estudante))
err := s.DB.First(&student, id).Error//para este caso usamos o First para procurar apenas um dado, diferente do Find. (metodo do GORM) usando o parametro ID.
return student, err
}

//metodo para salvar os update (atualização) do hendler.go
func (s *StudentHandler) UpdateStudent (updateStudent schemas.Student) error { 
	return s.DB.Save(&updateStudent).Error //Salva os dados atualizado
}

//metodo para deletar 
func (s *StudentHandler) DeleteStudent (student schemas.Student) error { 
	return s.DB.Delete(&student).Error //Deleta o dado informado
}