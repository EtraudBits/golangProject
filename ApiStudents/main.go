package main

import (
	"errors" // para manipulação de erros
	"fmt"
	"log/slog" // logger estruturado da biblioteca Padrão
	"net/http" // códigos e funções HTTP da standard lib

	"github.com/EtraudBits/golangProject/ApiStudents/db" //importa o modulo da pasta go.mod - acrescentando no final /db
	"github.com/labstack/echo/v4"                        // framework web usado para criar rotas e servidor HTTP
	"github.com/labstack/echo/v4/middleware"             //middleware já prontos do Echo
)

func main() {
  // Echo instance
  e := echo.New() //cria uma nova instância da estrutura principal do Echo, que representa o servidor WEB.

  // Middleware
  e.Use(middleware.Logger()) //Adiciona um middleware que registra logs de cada requisição HTTP feita no servidor.
  e.Use(middleware.Recover()) // adicona um middleware que recupera o servidor caso aconteça um panic, evitando que a aplicação caia.

  // Routes // Define rotas HTTP usando métodos
  e.GET("/students", getStudents) // buscar lista de estudantes
  e.POST("/students", createStudent) //criar estudantes
  e.GET("/students/:id", getStudent) // No singular queremos pegar apenas um estudante
  e.PUT("/students/:id", updateStudent) //Atualizar 
  e.DELETE("/students/:id", deleteStudent) //Deletar

  // Start server // inica o servidor na porta "8080"
  if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
    slog.Error("failed to start server", "error", err) // caso ocorra um erro que não seja o erro padrão de servidor fechado (http.ErrServerClosed), registra o erro usando slog.Error.
  }
}

// Handler //Funcões executada quando a rota é chamada
func getStudents(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  return c.String(http.StatusOK, "List of all students") // retorna uma resposta HTTP com status 200 (ok)
}

func createStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder. Função que recebe o POST
  student := db.Student {} // cria um bind para organizar as informações dinamica com o JSON

  if err := c.Bind(&student); err != nil { //Adequa a estrutura do struct 
    return err // já que na funcão createStuendet retorna error, caso não consiga executar a função retorna erro.
  }
  if err := db.AddStudent(student); err != nil {
     return c.String(http.StatusInternalServerError, "Error to Create student") // Caso ocorra o erro -> retorna uma resposta HTTP com erro interno do Servidor
  } //funçao publica na pacote db -> Chama a função student dinamica acima, tratando o erro, se houver!
  return c.String(http.StatusOK, "Create student") // canso não ocorra erro -> retorna uma resposta HTTP com status 200 (ok)
  
}

func getStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  id := c.Param("id") //Obtém o parâmetro de rota chamado "id" da URL. ->ex.: em /studante/10 -> id será "10"
  getStud := fmt.Sprintf("Get %s student", id)
  return c.String(http.StatusOK, getStud) // retorna uma resposta HTTP com status 200 (ok)
}

func updateStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  id := c.Param("id")// captura o "id" enviado na URL para identificar qual recurso deve ser atualizado
  updateStud := fmt.Sprintf("Update %s student", id)
  return c.String(http.StatusOK, updateStud) // retorna uma resposta HTTP com status 200 (ok)
}

func deleteStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  id := c.Param("id") // Recebe o parâmetro "id" que indica qual estudante será deletado
  deleteStud := fmt.Sprintf("Delete %s student", id)
  return c.String(http.StatusOK, deleteStud) // retorna uma resposta HTTP com status 200 (ok)
}