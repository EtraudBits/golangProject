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

// StockService define o que o budget precisa saber sobre estoque
type StockService interface { // interface para checar estoque -> para o budget não depender diretamente do módulo de estoque
	// Saida reduz o estoque de um produto
	Saida(ctx context.Context, productID int, quantity float64) error
	// Entrada aumenta o estoque de um produto
	Entrada(ctx context.Context, productID int, quantity float64) error
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
	stock StockService // checa estoque (via interface)
}

// Construtor do Service (falicita testes e facilita manutenção) -> injeção de dependência
func NewService (repo *Repository, product ProductReader, stock StockService,) *Service {
	return &Service{
		repo: repo,
		product: product,
		stock: stock,
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
	
	// Dar saída no estoque para cada item do orçamento
	for _, item := range budgetItems {
		err := s.stock.Saida(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return nil, err
		}
	}
	
	budget.ID = int64(id)
	return budget, nil
}

// GetByID retorna um orçamento completo (cabeçalho + itens)
func (s *Service) GetByID(ctx context.Context, id int64) (*Budget, error) {
	// 1 -> Buscar o orçamento (cabeçalho)
	budget, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if budget == nil {
		return nil, errors.New("orçamento não encontrado")
	}

	// 2 -> Buscar os itens do orçamento
	items, err := s.repo.ListItemsByBudget(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3 -> Associar os itens ao orçamento
	budget.Items = items

	// 4 -> Retornar orçamento completo
	return budget, nil
}

	
// List retorna todos os orçamentos com seus itens
func (s *Service) List(ctx context.Context) ([]Budget, error) {
	
	// delega a busca pra o repository
	budgets, err := s.repo.ListBudgets(ctx)
	if err != nil {
		return nil, err
	}

	// regras futuras podem ser aplicadas aqui
	// - filtros
	// - paginação
	// - permissões

	// retorno final	
	return budgets, nil
	}

	// cria metodo cancelar orçamento
	func (s *Service) Cancel(ctx context.Context, id int64) error {
		// buscar item do orçamento
		items, err := s.repo.ListItemsByBudget(ctx, id)
		if err != nil {
			return err		
		}

		if len(items) == 0 {
			return errors.New("orçamento não encontrado ou sem itens")
		}

		// Devolve o estoque para cada item
		for _, item := range items {
			err := s.stock.Entrada(ctx, item.ProductID, item.Quantity) //metodo entrada no stock/service.go
			if err != nil {
				return err
			}
		}

		// Atualiza o status do orçamento para "Cancelado"
		return s.repo.Cancel(ctx, id)
	}