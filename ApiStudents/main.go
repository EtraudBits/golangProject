package main

import (
  "github.com/labstack/echo/v4" // framework web usado para criar rotas e servidor HTTP
  "github.com/labstack/echo/v4/middleware" //middleware já prontos do Echo
  "log/slog" // logger estruturado da biblioteca Padrão
  "net/http" // códigos e funções HTTP da standard lib
  "errors" // para manipulação de erros
)

func main() {
  // Echo instance
  e := echo.New() //cria uma nova instância da estrutura principal do Echo, que representa o servidor WEB.

  // Middleware
  e.Use(middleware.Logger()) //Adiciona um middleware que registra logs de cada requisição HTTP feita no servidor.
  e.Use(middleware.Recover()) // adicona um middleware que recupera o servidor caso aconteça um panic, evitando que a aplicação caia.

  // Routes // Define rotas HTTP usando métodos
  e.GET("/students", getStudent)

  // Start server // inica o servidor na porta "8080"
  if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
    slog.Error("failed to start server", "error", err) // caso ocorra um erro que não seja o erro padrão de servidor fechado (http.ErrServerClosed), registra o erro usando slog.Error.
  }
}

// Handler //Funcões executada quando a rota é chamada
func getStudent(c echo.Context) error { //recebe um echo.context, que contém informações da requisição e métodos para responder.
  return c.String(http.StatusOK, "List of all students") // retorna uma resposta HTTP com status 200 (ok)
}