package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EtraudBits/golangProject/gobuild/internal/handler"
	"github.com/labstack/echo/v4"
)

// Server represtenta o servidor HTTP usando Echo framework
type Server struct {
	Echo *echo.Echo
}

// New cria uma nova instância do servidor
func New() *Server {
	e := echo.New()

	//e.Use(middlleware.Logger())
	//e.Use(middlleware.Recover())

	return &Server{
		Echo: e,
	}

}
// RegisterRoutes registra todas as rotas da API
func (s *Server) RegisterRoutes() {

	// Rota Raiz
	s.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API gobuild rodando com SQLite!")
	})

	// Rota de Testa a conexão do Banco de Dados
	s.Echo.GET("/db-test", handler.TestDBHandler)
}

func (s *Server) Start() {
	fmt.Println("Iniciando o servidor na porta 8080...")

	if err := s.Echo.Start(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}	