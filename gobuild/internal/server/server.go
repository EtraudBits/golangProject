package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EtraudBits/golangProject/gobuild/internal/database" // conexão com o banco de dados (DB)
	"github.com/EtraudBits/golangProject/gobuild/internal/handler"  // pacote antigo de handlers (db test)
	"github.com/EtraudBits/golangProject/gobuild/internal/product"  // pacote product (novo)
	"github.com/labstack/echo/v4"
)

// Server represtenta o servidor HTTP usando Echo framework
type Server struct {
	Echo *echo.Echo // instância do Echo usada para gerenciar rotas e requisições
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
// além das rotas existentes, registra as rotas de produto /api/products com o CRUD completo
func (s *Server) RegisterRoutes() {

	// Rota Raiz
	s.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API gobuild rodando com SQLite!")
	})

	// Rota de Testa a conexão do Banco de Dados
	s.Echo.GET("/db-test", handler.TestDBHandler)


// --- Configuração do módulo de Produtos ---
// cria repositório, serviço e handler para produtos usando a conexão do DB.
// database.DB vem do pacote internal/database, previamente inicializado em main.go.
repo := product.NewRepository(database.DB) // cria repositório com a conexão do DB
svc := product.NewService(repo) // injeta repo no service
h := product.NewHandler(svc) // cria o handler http para products

//cria um grupo de rotas para organizar endpoints realcionados a produtos.
// As rotas ficarão sob o prefixo /api/products
g := s.Echo.Group("/api/products")
//registra as rotas de produto no grupo criado
h.RegisterRoutes(g)
}

// Start inicia o servidor HTTP na porta 8080
// caso de erro, encerra o programa com log fatal
func (s *Server) Start() {
	fmt.Println("Iniciando o servidor na porta 8080...")

	if err := s.Echo.Start(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}	