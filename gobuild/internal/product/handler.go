package product

import (
	"net/http" // para constantes de status HTTP
	"strconv"  // para conversão de string para int

	"github.com/labstack/echo/v4" // framework web Echo
)

// handler contém as dependências para os manipuladores de produto (service - injeção simples)
type Handler struct {
	svc *Service // serviço que contém regras de negócio
}

// NewHandler cria um novo handler com a dependência do serviço injetada
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// registerRoutes registra as rotas de produto no grupo Echo fornecido
// passaremos um grupo Echo para organizar as rotas de produto
func (h *Handler) RegisterRoutes(g *echo.Group) {
	//POST/api/products - > criar produto
	g.POST("", h.Create)
	// GET/api/products - > listar todos produtos
	g.GET("", h.list)
	//GET /api/products/:id -> obter produto por ID
	g.GET("/:id", h.Get)
	// PUT /api/products/:id -> atualizar produto por ID
	g.PUT("/:id", h.Update)
	// DELETE /api/products/:id -> deletar produto por ID
	g.DELETE("/:id", h.Delete)
	// Rota para consultar o estoque de um produto
	g.GET("/:id/stock", h.GetStock)

}
// Create lida com a criação de um produto via JSON no corpo da requisição.
// Retorna 201 Created com o id do novo produto ou 400/500 em caso de erro.
func (h *Handler) Create(c echo.Context) error {
	var req Produto 
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos: " + err.Error()})
	}
	//chama o serviço para criar o produto
	id, err := h.svc.Create(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "falha ao criar produto: " + err.Error()})
	}
	//retorna 201 Created com o ID do novo produto
	return c.JSON(http.StatusCreated, map[string]int64{"id": id})
}

// List retorna todos os produtos cadastrados no banco de dados.
// Retorna 200 OK com a lista de produtos ou 500 em caso de erro.
func (h *Handler) list(c echo.Context) error {
	//chama o serviço para obter todos os produtos
	produtos, err := h.svc.List(c.Request().Context())
	if err != nil {
		// retorna 500 em caso de erro
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "falha ao listar produtos: " + err.Error()})
	}
	//retorna 200 OK com a lista de produtos
	return c.JSON(http.StatusOK, produtos)
}
// Get retorna um produto por id. Lida com id inválido e 404 se não encontrado.
func (h *Handler) Get(c echo.Context) error {
	// Lê parâmetro :id da URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}
	// busca via serviçe
	p, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		// se o erro indicar que o produto não foi encontrado, retorna 404
		if err.Error() == "produto com ID "+strconv.Itoa(id)+" não encontrado" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "falha ao obter produto: " + err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}
// Update atualiza um produto por id.
func (h *Handler) Update(c echo.Context) error {
	// Lê parâmetro :id da URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	// bind do JSON para struct Produto
	var req Produto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dados inválidos: " + err.Error()})
	}

	//garante que o id do payload seja o mesmo do URL
	req.ID = id

	//chama o serviço para atualizar o produto
	if err := h.svc.Update(c.Request().Context(), &req); err != nil {
		// se o erro indicar que o produto não foi encontrado, retorna 404
		if err.Error() == "produto com ID "+strconv.Itoa(id)+" não encontrado" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "falha ao atualizar produto: " + err.Error()})
	}
	//retorna 200 OK com mensagem de sucesso
	return c.JSON(http.StatusOK, map[string]string{"message": "produto atualizado com sucesso"})
}

// Delete remove um produto por id.
// Retorna 204 No Content em caso de sucesso (sem corpo) ou 400/404/500 em caso de erro.
func (h *Handler) Delete(c echo.Context) error {
	// Lê parâmetro :id da URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}

	// chama service para deletar
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		// se o erro indicar que o produto não foi encontrado, retorna 404
		if err.Error() == "produto com ID "+strconv.Itoa(id)+" não encontrado" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "falha ao deletar produto: " + err.Error()})
	}
	//retorna 204 No Content em caso de sucesso
	return c.NoContent(http.StatusNoContent)
}

// GetStock retorna o estoque atual de um produto
func (h *Handler) GetStock(c echo.Context) error {
	// Lê o ID da URL
	idStr := c.Param("id")

	//Converte string em Int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error" : "id inválido",
		})
	}
	// Chama o Service para buscar o Estoque
	stock, err := h.svc.GetStock(c.Request().Context(), id)
	if err != nil {
		// se o produto não existir, retorna 404
		return c.JSON(http.StatusNotFound, map[string]string{
			"error" : err.Error(),
		})
	}
	// Retorna apenas o que o cliente precisa
	return c.JSON(http.StatusOK, map[string]interface{}{
		"product_id": id,
		"estoque": stock,
	})
}