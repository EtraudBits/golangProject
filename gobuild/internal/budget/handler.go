package budget

import (
	"net/http"

	"github.com/labstack/echo/v4"
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