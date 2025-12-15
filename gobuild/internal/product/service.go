package product

import (
	"context" // Para passar contexto em operações de banco de dados
	"errors"  // para manipulação de erros
	"fmt"     // para formatação de strings e erros
)

type Service struct {
	repo *Repository // dependencia do repositorio para persistencia
}
// NewService cria uma nova instância do serviço de produtos
func NewService(r *Repository) *Service {
	return &Service{
		repo: r,
	}		
}
// validateProduto realiza validações basicas no produto antes de salvar/atualizar.
func (s *Service) ValidateProduto(p *Produto) error {
	// nome não pode ser vazio
	if p.Name == "" {
		return errors.New("o nome do produto não pode ser vazio")
	}
	// preço não pode ser negativo
	if p.Preco < 0 {
		return errors.New("o preço do produto não pode ser negativo")
	}
	// estoque não pode ser negativo
	if p.Estoque < 0 {
		return errors.New("o estoque do produto não pode ser negativo")
	}
	// unidade não pode ser vazia
	if p.Unidade == "" {
		return errors.New("a unidade do produto não pode ser vazia")
	}
	// categoria não pode ser vazia
	if p.Categoria == "" {
		return errors.New("a categoria do produto não pode ser vazia")
	}
	return nil // todas as validações passaram

}
// criar método lista todos os produtos cadastrados
func (s *Service) List(ctx context.Context) ([]Produto, error) {
	//chama o repositório para obter todos os produtos
	produtos, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %v", err)
	}
	return produtos, nil
}
// create cria um novo produto após validação dos dados

func (s *Service) Create(ctx context.Context, p *Produto) (int64, error) {
	// valida o produto antes de criar
	if err := s.ValidateProduto(p); err != nil {
		return 0, fmt.Errorf("validação do produto falhou: %v", err)
	}
	//cria via repo.
	return s.repo.Create(ctx, p)
}

// Get retorna produto por ID, ou erro se não encontrado
func (s *Service) Get(ctx context.Context, id int) (*Produto, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produto: %v", err)
	}
	if p == nil {
		// retorna erro para o handler decidir status 404
		return nil, fmt.Errorf("produto com ID %d não encontrado", id)
	}	
	return p, nil
}
// update atualiza os dados de um produto existente
func (s *Service) Update(ctx context.Context, p *Produto) error {

	if err := s.ValidateProduto(p); err != nil {
		return fmt.Errorf("validação do produto falhou: %v", err)

	}

	//verifica se o produto existe
existing, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		return fmt.Errorf("erro ao obter produto existente: %v", err)
	}

	if existing == nil {
		return fmt.Errorf("produto com ID %d não encontrado", p.ID)
	}

	//atualiza apenas se o produto existir
	return s.repo.Update(ctx, p)
}

// delete remove um produto pelo ID
func (s *Service) Delete(ctx context.Context, id int) error {
	//verifica se o produto existe
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao obter produto existente: %v", err)
	}
	// se não existir, retorna erro
	if existing == nil {
		return fmt.Errorf("produto com ID %d não encontrado", id)
	}
	//deleta o produto
	return s.repo.Delete(ctx, id)
}
// --- Função para pesquisar o estoque de cada produto --
// GetStock retorna apenas o estoque atual de um produto
func (s *Service) GetStock(ctx context.Context, id int) (float64, error) {
	// Busca o produto pelo ID usando o repositório
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, nil
	}
	// Se o produto não existir, retorna erro
	if p == nil {
		return 0, errors.New("Produto não encontrado")
	}
	// Retorna apenas o estoque
	return p.Estoque, nil
}