package budget

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"strconv"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes (g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id/cancel", h.Cancel)
}

// CreateItemRequest represeta um item enviado pelo o cliente
type CreateItemRequest struct {
	ProductID int `json:"product_ID"`
	Quantity float64 `json:"quantity"`
}
type CreateBudgetRequest struct {
	Customer string `json:"customer"`
	Items []CreateItemRequest `json:"items"`
}

func (h *Handler) Create(c echo.Context) error {

	var req CreateBudgetRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "JSON invalido",
		})
	}
	budget, err := h.svc.Create(
		c.Request().Context(),
		req.Customer,
		req.Items,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, budget)
}

func (h *Handler) List(c echo.Context) error {

	budgets, err := h.svc.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error" : err.Error(),
		})
	}
	return c.JSON(http.StatusOK, budgets)
}

// GetByID retorna um orçamento especifico
func (h *Handler) GetByID(c echo.Context) error {

	// 1 -> ler ID da URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{ // trata o erro com status 400
			"error": "id inválido",
		})
	}

	// 2 -> chamar service
	budget, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "orçamento não encontrado" {
		return c.JSON(http.StatusNotFound, map[string]string{ // trata o erro com status 404
			"error": err.Error(),
		})
	}
		return c.JSON(http.StatusInternalServerError, map[string]string{ // trata o erro com status 500
			"error": err.Error(),
		})
		
	}
	// 3 -> retornar resposta
	return c.JSON(http.StatusOK, budget) // trata o sucesso com status 200
}

// Cancel cancela um orçamento
func (h *Handler) Cancel (c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string {
			"error": "id inválido",
		})
	}
	if err := h.svc.Cancel(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string {
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string {
		"message": "orçamento cancelado com sucesso",
	})
}