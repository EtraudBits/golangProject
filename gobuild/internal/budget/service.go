package budget

import (
	"context" // padrão GO para requests, banco, cancelamento
	"errors"  // criar erros claros de negócio
)

//cria uma interface que não depende diretamente do modulo product
// ProductReader define o que o budget precisa saber sobre produtos

type ProductReader interface {
	GetByID(ctx context.Context, id int) (*ProductLite, error) // o budget só quer: ID, Nome e Preço -> então usaremos o ProductLite (struct)
}

type ProductLite struct {
	ID int
	Name string
	Price float64
}

//Criação do Service

type Service struct {
	repo *Repository //fala com o banco (budget_repository)
	product ProductReader //lê produtos (via interface)
}

// Construtor do Service (falicita testes e facilita manutenção) -> injeção de dependência
func NewService (repo *Repository, product ProductReader) *Service {
	return &Service{
		repo: repo,
		product: product,
	}
}

//Regra principal (criar Orçamento)
func (s *Service) Create(ctx context.Context, customer string, items []CreateItemRequest) (*Budget, error) {
	if customer == "" {
		return nil, errors.New("cliente é obrigatório!")
	}

	if len(items) == 0 {
		return nil, errors.New("orçamento precisa de ao menos um item")
	}

	budget := &Budget{
		Customer: customer,
		Total: 0, //total começa zerado.
	}

	//processar itens
	var budgetItems []BudgetItem // lista de itens finais

	for _, item := range items {
		p, err := s.product.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, errors.New("produto não encontrado")
		}
		//calcula o subtotal
		subtotal := item.Quantity * p.Price

		//criação item do orçamento
		bi := BudgetItem {
			ProductID: p.ID,
			Product: p.Name,
			Quantity: item.Quantity,
			UnitPrice: p.Price,
			Subtotal: subtotal,
		}
		//somar total
		budget.Total += subtotal
		budgetItems = append(budgetItems, bi)
	}
	//Salva no banco (persistencia isolada no repository)
	id, err := s.repo.CreateBudget(ctx, budget, budgetItems)
	if err != nil {
		return nil, err
	}
	budget.ID = int64(id)

	return budget, nil
}
	

//DTO (request auxiliar)
// CreateItemRequest representa um item recebido via API
type CreateItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity float64 `json:"quantity"`
}

