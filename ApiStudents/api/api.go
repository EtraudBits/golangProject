package api

import (
	// para manipulação de erros
	"fmt"      // logger estruturado da biblioteca Padrão
	"net/http" // códigos e funções HTTP da standard lib

	"github.com/EtraudBits/golangProject/ApiStudents/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// cria uma struct TIPO api com dois campos -> echo e o db
type API struct {
	Echo *echo.Echo // (*) PONTEIRO tras tudo que o echo irá usar
	DB *db.StudentHandler //(*) PONTEIRO tras tudo que o db irá usar 
}


func NewServer() *API {//Cria uma função que inicializa tudo -> função para ser o servidor -> Usando a Struct API

  // Echo instance
  e := echo.New() //cria uma nova instância da estrutura principal do Echo, que representa o servidor WEB.

  // Middleware
  e.Use(middleware.Logger()) //Adiciona um middleware que registra logs de cada requisição HTTP feita no servidor.
  e.Use(middleware.Recover()) // adicona um middleware que recupera o servidor caso aconteça um panic, evitando que a aplicação caia.

  database := db.Init() //puxa a inicialização do db

	studentDB := db.NewStudentHandler(database)

  return &API{
	Echo: e,
	DB: studentDB,
  }
}
 
  //função para organizar (configurar) as Rotas
  func (api *API) ConfigureRoutes() {
	// Routes // Define rotas HTTP usando métodos
  api.Echo.GET("/students", api.getStudents) // buscar lista de estudantes -> atrelar a api
  api.Echo.POST("/students", api.createStudent) //criar estudantes
  api.Echo.GET("/students/:id", api.getStudent) // No singular queremos pegar apenas um estudante
  api.Echo.PUT("/students/:id", api.updateStudent) //Atualizar 
  api.Echo.DELETE("/students/:id", api.deleteStudent) //Deletar
  }

  //função para startar o server 
  func (api *API) Start () error { //aponta para API
 // Start server // inica o servidor na porta "8080"
 	return api.Echo.Start(":8080")
  }


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
  id := c.Param("id") //Obtém o parâmetro de rota chamado "id" da URL. ->ex.: em /studante/10 -> id será "10"
  getStud := fmt.Sprintf("Get %s student", id)
  return c.String(http.StatusOK, getStud) // retorna uma resposta HTTP com status 200 (ok)
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