package api

import (
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


