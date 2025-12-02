package api

import (
	"github.com/EtraudBits/golangProject/ApiStudents/db"

	_ "github.com/EtraudBits/golangProject/ApiStudents/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// cria uma struct TIPO api com dois campos -> echo e o db
type API struct {
	Echo *echo.Echo // (*) PONTEIRO tras tudo que o echo irá usar
	DB *db.StudentHandler //(*) PONTEIRO tras tudo que o db irá usar 
}

// @title Sdudent API
// @version 1.0
// @description This is a sample server Student API
// @host localhost:8080
// @BasePath /
// @schemes http

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
  api.Echo.GET("/swagger/*", echoSwagger.WrapHandler) //rota para configurar o swagger
  }

  //função para startar o server 
  func (api *API) Start () error { //aponta para API
 // Start server // inica o servidor na porta "8080"
 	return api.Echo.Start(":8080")
  }


