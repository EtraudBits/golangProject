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

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id/cancel", h.Cancel)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

// CreateItemRequest representa um item enviado pelo cliente
type CreateItemRequest struct {
	ProductID int     `json:"product_ID"`
	Quantity  float64 `json:"quantity"`
}

// CreateBudgetRequest representa os dados para criar um orçamento
type CreateBudgetRequest struct {
	Customer string              `json:"customer"`
	Items    []CreateItemRequest `json:"items"`
}

// UpdateBudgetRequest representa os dados para atualizar um orçamento
type UpdateBudgetRequest struct {
	Customer string              `json:"customer"`
	Items    []CreateItemRequest `json:"items"`
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
			"error": err.Error(),
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
func (h *Handler) Cancel(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
	}
	err = h.svc.Cancel(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "orçamento não encontrado ou sem itens" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "orçamento cancelado com sucesso",
	})
}

// Update atualiza um orçamento existente
func (h *Handler) Update(c echo.Context) error {

	// 1 -> ler ID da URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
	}

	// 2 -> ler o corpo da requisição
	var req UpdateBudgetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
	}

	// 3 -> chamar o service
	budget, err := h.svc.Update(
		c.Request().Context(),
		id,
		req.Customer,
		req.Items,
	)
	if err != nil {
		if err.Error() == "orçamento não encontrado" || err.Error() == "produto não encontrado" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// 4 -> retornar resposta com mensagem e budget atualizado
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "orçamento atualizado com sucesso",
		"budget":  budget,
	})
}

// Delete remove um orçamento e seus itens
func (h *Handler) Delete(c echo.Context) error {
	// 1 -> ler ID da URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
	}

	// 2 -> chamar o service para deletar
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "orçamento não encontrado" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// 3 -> retornar 204 No Content (padrão REST para DELETE bem-sucedido)
	return c.NoContent(http.StatusNoContent)
}
