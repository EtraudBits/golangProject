package api

import (
	"fmt"
	"net/http"
	"strconv"

	"errors"

	"github.com/EtraudBits/golangProject/ApiStudents/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Handler //Funcões executada quando a rota é chamada
func (api *API) getStudents(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  students, err := api.DB.GetStudents() //funções que chama do repositório db a Função GetStudents.
 
  if err != nil {
    return c.String(http.StatusNotFound, "Failed to get students") //podemos usar o StatusNotFound, caso não tenha nenhum usuario (estudante) cadastrarado ele mostra a msg.
  }
  return c.JSON(http.StatusOK, students) // retorna uma resposta HTTP com status 200 (ok) -> para aplicações mais robustas chama o c.JSON
}

func (api *API) createStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder. Função que recebe o POST
  student := db.Student {} // cria um bind para organizar as informações dinamica com o JSON

  if err := c.Bind(&student); err != nil { //Adequa a estrutura do struct 
    return err // já que na funcão createStuendet retorna error, caso não consiga executar a função retorna erro.
  }
  if err := api.DB.AddStudent(student); err != nil {
     return c.String(http.StatusInternalServerError, "Error to Create student") // Caso ocorra o erro -> retorna uma resposta HTTP com erro interno do Servidor
  } //funçao publica na pacote db -> Chama a função student dinamica acima, tratando o erro, se houver!
  return c.String(http.StatusOK, "Create student") // canso não ocorra erro -> retorna uma resposta HTTP com status 200 (ok)
  
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
  id := c.Param("id")// captura o "id" enviado na URL para identificar qual recurso deve ser atualizado
  updateStud := fmt.Sprintf("Update %s student", id)
  return c.String(http.StatusOK, updateStud) // retorna uma resposta HTTP com status 200 (ok)
}

func (api *API) deleteStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  id := c.Param("id") // Recebe o parâmetro "id" que indica qual estudante será deletado
  deleteStud := fmt.Sprintf("Delete %s student", id)
  return c.String(http.StatusOK, deleteStud) // retorna uma resposta HTTP com status 200 (ok)
}