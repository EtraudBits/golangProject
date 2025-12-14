package stock

import (
	"net/http" // para constantes de status HTTP
	"strconv"  // para conversão de strings

	"github.com/labstack/echo/v4" // framework web Echo
	"golang.org/x/net/context"    // para contexto em handlers
)

// Handler expoe endpoints HTTP para movimentoações de estoque.

type Handler struct {
	svc *Service // serviço de estoque
}

// NwHandler cria um handler com o serviço injetado
func NewHandler(svc *Service) *Handler { 
	return &Handler{svc: svc} // inicializa o handler com o serviço
}

//RegistreRoutes registra as rotas de estoque num grupo Echo,
//ex.: g := e.Group("/api/stock"); h.RegisterRoutes(g).
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/entrada", h.Entrada)
	g.POST("/saida", h.Saida)
	g.POST("/ajuste", h.Ajuste)
	g.GET("/historico/:product_id", h.Historico)
}

//Entrada esperam JSON: {"product_id": 1, "quantity": 10, "description": "Compra fornecedor"}
type movimentRequest struct {
	ProductID int `json:"product_id"`
	Quantity float64 `json:"quantity"`
}

// Entrada cria um movimento de tipo ENTRADA
func (h *Handler) Entrada(c echo.Context) error {
	var req movimentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "requisição inválida"})
	}
	m := &Movement {
		ProductID: req.ProductID,
		Type: "Entrada",
		Quantity: req.Quantity,
	}
	
	id, err:= h.svc.CreateMovement(c.Request().Context(), m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int64{"movement_id": id})
	}
	
	// saida cria um movimento de tipo SAIDA
	func (h *Handler) Saida(c echo.Context) error {
	var req movimentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "requisição inválida"})
	}
	
	m := &Movement{
		ProductID: req.ProductID,
		Type: "Saida",
		Quantity: req.Quantity,
	}

	id, err := h.svc.CreateMovement(c.Request().Context(), m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "requisição inválida"})
	}
	return c.JSON(http.StatusCreated, map[string]int64{"movement_id": id})
}

	//Ajuste define o estoque diretamente (tipo AJUSTE) - quantity é o novo estoque
	func (h *Handler) Ajuste(c echo.Context) error {
		var req movimentRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

	m := &Movement{
		ProductID: req.ProductID,
		Type: "Ajuste",
		Quantity: req.Quantity,
	}

	id, err := h.svc.CreateMovement(c.Request().Context(), m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]int64{"movement_id": id})
	}

	//Historico retorna lista de movimentos de um produto
	func (h *Handler) Historico(c echo.Context) error {
		idStr := c.Param("product_id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Product_id inválido"})
		}

		list, err := h.svc.GetHistory(context.Background(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro ao obter historico do produto"})
		}
		return c.JSON(http.StatusOK, list)
	}



