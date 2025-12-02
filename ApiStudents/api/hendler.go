package api

import (
	"net/http"
	"strconv"

	"errors"

	"github.com/EtraudBits/golangProject/ApiStudents/schemas" //pacote do repositório schemas
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log" //formatar os logs (deixar mais organizado)
	"gorm.io/gorm"
)

// Handler //Funcões executada quando a rota é chamada
func (api *API) getStudents(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  students, err := api.DB.GetStudents() //funções que chama do repositório db a Função GetStudents.
 
  if err != nil {
    return c.String(http.StatusNotFound, "Failed to get students") //podemos usar o StatusNotFound, caso não tenha nenhum usuario (estudante) cadastrarado ele mostra a msg.
  }

  //cria uma função para chamar a lista de estudande com a nova formatação criada no schemas.go
  listOfStudents := map[string] []schemas.StudentResponse{"students:": schemas.NewResponse(students)}
  
  return c.JSON(http.StatusOK, listOfStudents) // retorna uma resposta HTTP com status 200 (ok) -> para aplicações mais robustas chama o c.JSON
}

func (api *API) createStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder. Função que recebe o POST
  studentReq := StudentRequest{}
  // cria um bind para organizar as informações dinamica com o JSON -> usa os dados do request.go para validar a criação
  if err := c.Bind(&studentReq); err != nil { //Adequa a estrutura do struct 
    return err // já que na funcão createStuendet retorna error, caso não consiga executar a função retorna erro.
  }

  //chama o metodo validate do request.go
  if err := studentReq.Validate(); err != nil {
    log.Error().Err(err).Msg("[api] error validating struct")
    return c.String(http.StatusBadRequest, "Erro validating student")
  }
  //converter a estrutura para o student
  student := schemas.Student {
  //faz um De - Para 
  Name: studentReq.Name, //Name recebe sturdentReq.Name
  Email: studentReq.Email, //Email recebe studentReq.Email
  CPF: studentReq.CPF, //CPF recebe studentReq.CPF
  Age: studentReq.Age, //Age recebe studentReq.Age
  Active: *studentReq.Active, //Active recebe studentReq.Active -> pegando o ponteiro (*)
}

  if err := api.DB.AddStudent(student); err != nil {
     return c.String(http.StatusInternalServerError, "Error to Create student") // Caso ocorra o erro -> retorna uma resposta HTTP com erro interno do Servidor
  } //funçao publica na pacote db -> Chama a função student dinamica acima, tratando o erro, se houver!
  return c.JSON(http.StatusOK, student) // canso não ocorra erro -> retorna uma resposta HTTP com status 200 (ok)
  
}

func (api *API) getStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  id, err := strconv.Atoi(c.Param("id")) //Obtém o parâmetro de rota chamado "id" da URL. ->ex.: em /studante/10 -> id será "10" -> transforma a string em Inteiro usando strconv.Atoi
  if err != nil { //trata o erro
	return c.String(http.StatusInternalServerError, "Failed to Get student ID") //msg do erro quando vem do servidor (por ex.: digitou um ID que não exista no BD)
  }
  
  student, err := api.DB.GetStudent(id) //Pode não encontrar um student com esse id -> STATUS NOT FOUND (404) ou pode ter algum problema para encontrar o student (temos que tratar esses erros)
  if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found") 
  }

  if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to Get student") //msg do erro quando vem do servidor (por ex.: digitou um ID que não exista no BD)
  }

  return c.JSON(http.StatusOK, student) // se não ocorrer nenhum dos dois erros acima -> retorna uma resposta HTTP com status 200 (ok)
}

func (api *API) updateStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
   id, err := strconv.Atoi(c.Param("id")) //Obtém o parâmetro de rota chamado "id" da URL. ->ex.: em /studante/10 -> id será "10" -> transforma a string em Inteiro usando strconv.Atoi
  if err != nil { //trata o erro
	return c.String(http.StatusInternalServerError, "Failed to Get student ID") //msg do erro quando vem do servidor (por ex.: digitou um ID que não exista no BD)
  }
  //precisamos procurar o que esta vindo do PUT para fazer a comparação com o dados atual.
  receivedStudent := schemas.Student{} //variavel do estudante recebido
  if err := c.Bind(&receivedStudent); err != nil {
	return err
  }
//função para atualização dos dados
	updatingStudent, err := api.DB.GetStudent(id) //Pode não encontrar um student com esse id -> STATUS NOT FOUND (404) ou pode ter algum problema para encontrar o student (temos que tratar esses erros)
  if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found") 
  }

  if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to Get student") //msg do erro quando vem do servidor (por ex.: digitou um ID que não exista no BD)
  }

  student := upDateStudentInfo(receivedStudent, updatingStudent)  //retorna studant (atualização) chamando a função upDateStudentInfo  -> que esta salvando no bd
	
if err := api.DB.UpdateStudent(student); err != nil {
	return c.String(http.StatusInternalServerError, "Failed to save student")
}
  
  return c.JSON(http.StatusOK, student) // retorna uma resposta HTTP com status 200 (ok)
}

func (api *API) deleteStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  id, err := strconv.Atoi(c.Param("id")) //Obtém o parâmetro de rota chamado "id" da URL. ->ex.: em /studante/10 -> id será "10" -> transforma a string em Inteiro usando strconv.Atoi
    if err != nil { //trata o erro
	return c.String(http.StatusInternalServerError, "Failed to Get student ID") //msg do erro quando vem do servidor (por ex.: digitou um ID que não exista no BD)
}
	student, err := api.DB.GetStudent(id) //Pode não encontrar um student com esse id -> STATUS NOT FOUND (404) ou pode ter algum problema para encontrar o student (temos que tratar esses erros)
  if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found") 
  }

  if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to Get student") //msg do erro quando vem do servidor (por ex.: digitou um ID que não exista no BD)
  }

  if err := api.DB.DeleteStudent(student); err != nil {
	return c.String(http.StatusInternalServerError, "Failed to delete student")
  }
  
  return c.JSON(http.StatusOK, student) // retorna uma resposta HTTP com status 200 (ok)
}

//busca a função de deletar do db

//função para fazer a comparação do receivedStudent e student (updatingStudant) - do tipo db.Student
func upDateStudentInfo(receivedStudent, student schemas.Student) schemas.Student {
	if receivedStudent.Name != "" { // se o campo Name do receivedStudent for diferente de vazio -> Name é uma string
		student.Name = receivedStudent.Name //retorna Name do student.Name (atualizando o campo)
	}
	if receivedStudent.CPF > 0 { // se o campo CPF for maio que 0 do receivedStudent -> CPF é um Int por isso do maior que
		student.CPF = receivedStudent.CPF //retorna CPF do student.CPF (atualizando o campo) 
	}
	if receivedStudent.Email != "" { // se o campo Email do receivedStudent for diferente de vazio -> Email é uma string
		student.Email = receivedStudent.Email //retorna Email do student.Email (atualizando o campo)
	}
	if receivedStudent.Age > 0 { // se o campo Name do receivedStudent for diferente de vazio ->  Age é um Int por isso do maior que
		student.Age = receivedStudent.Age //retorna Name do student.Name (atualizando o campo)
	}
	if receivedStudent.Active != student.Active{ //se o campo Active do receivedStudent for diferente do campo Active do student -> Active é um bool (comparação)
		student.Active = receivedStudent.Active//retorna Active do student.Active (atualizando o campo)
	}
	return student //retorna a função dos dados atualizado

}