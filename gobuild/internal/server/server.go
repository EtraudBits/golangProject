package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/EtraudBits/golangProject/gobuild/internal/budget"
	"github.com/EtraudBits/golangProject/gobuild/internal/database"
	dbhandler "github.com/EtraudBits/golangProject/gobuild/internal/handler" // handler de /db-test
	"github.com/EtraudBits/golangProject/gobuild/internal/product"
	stockpkg "github.com/EtraudBits/golangProject/gobuild/internal/stock"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server é o wrapper do Echo usado para organizar o app
type Server struct {
	Echo *echo.Echo
}

// New cria o servidor com middlewares básicos
func New() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return &Server{Echo: e}
}

// RegisterRoutes registra todas as rotas da aplicação
func (s *Server) RegisterRoutes() {
	// rota raiz
	s.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API gobuild rodando com SQLite! 🚀")
	})

	// rota de teste do banco
	s.Echo.GET("/db-test", dbhandler.TestDBHandler)

	// --- produtos (já existentes) ---
	repo := product.NewRepository(database.DB)
	svc := product.NewService(repo)
	h := product.NewHandler(svc)
	gp := s.Echo.Group("/api/products")
	h.RegisterRoutes(gp)

	// --- estoque (criado antes do budget para injeção de dependência) ---
	stockRepo := stockpkg.NewRepository(database.DB)

	// Função injetada para ler produto (ProductLite) — usa o produto repo/serviço já existente.
	getProduct := func(ctx context.Context, id int) (float64, error) {
		// reutilizamos repository product para buscar apenas stock (podia ser otimizada)
		p, err := repo.GetByID(ctx, id)
		if err != nil {
			return 0, err
		}
		if p == nil {
			return 0, nil
		}
		return p.Estoque, nil
	}

	// Função injetada para atualizar estoque do produto
	updateStock := func(ctx context.Context, id int, newStock float64) error {
		// Faz update direto na tabela products
		_, err := database.DB.ExecContext(ctx, `UPDATE products SET stock = ? WHERE id = ?`, newStock, id)
		return err
	}

	stockSvc := stockpkg.NewService(database.DB, stockRepo, getProduct, updateStock)
	stockHandler := stockpkg.NewHandler(stockSvc)
	gs := s.Echo.Group("/api/stock")
	stockHandler.RegisterRoutes(gs)

	// -- Modulo budget (depois do stock, pois depende dele)
	// cria o repositório de budget -> fala com o banco
	budgetRepo := budget.NewRepository(database.DB)

	//cria o service de budget (injetando stockSvc)
	budgetSvc := budget.NewService(budgetRepo, svc, stockSvc)

	// cria o handler HTTP do budget
	budgetHandler := budget.NewHandler(budgetSvc)

	// cria um grupo de Rotas /api/budgets
	gb := s.Echo.Group("/api/budgets")
	budgetHandler.RegisterRoutes(gb)


}

// Start inicia o servidor
func (s *Server) Start() {
	fmt.Println("🔥 Servidor iniciado em http://localhost:8080")
	if err := s.Echo.Start(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
