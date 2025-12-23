package stock

import (
	"context"      // Para passar contexto em operações de banco de dados
	"database/sql" // pacote sql para manipulação de rows/ results
	"errors"       // para erros específicos
	"fmt"          // para formatação de strings e erros
)

// Service coordena regras de negócio para movimentações de estoque
// - verifica se o produto existe (pode usar repositório de produtos)
// - realiza a operação em transação (atualiza product.stock e insere movement)
// - previne estoque negativo (configuravel)
type Service struct {
	db *sql.DB			   // Conexão com o banco (injetada na criação do serviço)
	repo *Repository	   // Repositório de movimentações de estoque
	getProduct func(ctx context.Context, id int) (float64, error) // função para obter estoque do produto
	updateStock func(ctx context.Context, id int, newStock float64) error // função para atualizar estoque do produto
	//	allowNegative bool	   // se falso, previne estoque negativo (configurável é só deixar true por enquanto)
}
// ProductLite é uma visão reduzida do produto usada pelo serviço de estoque
// não precisamos de todos os campos, só do estoque atual
type ProductLite struct {
	ID int
	Stock float64
}

// NewService cria uma o serviço de estoque
// passamos também funções ulitárias para ler/atualizar o estoque do produto (injeção para simplicidade)
func NewService(db *sql.DB, repo *Repository,
	getProduct func(ctx context.Context, id int) (float64, error),
	updateStock func(ctx context.Context, id int, newStock float64) error,
	) *Service {
	return &Service{
		db: db,
		repo: repo,
		getProduct: getProduct,
		updateStock: updateStock,
	}
}

// helper: valida o tipo de movimento
func validType(t string) bool {
	return t == "Entrada" || t == "Saida" || t == "Ajuste"
}

//CreateMovement executa o fluxo completo de um movimento:
// 1. valida o tipo e quantidade
// 2. inicia transação
// 3. obtém estoque atual do produto
// 4. insere registro em stock_movements
// 5. commita a transação (ou rollback em caso de erro)
func (s *Service) CreateMovement(ctx context.Context, m *Movement) (int64, error) {
	//validações básicas
	if !validType(m.Type) {
		return 0, errors.New("tipo de movimentação inválido")
	}
	if m.Quantity <= 0 {
		return 0, errors.New("quantidade deve ser maior que zero")
	}

	// lê produto atual (via função injetada)
	currentStock, err := s.getProduct(ctx, m.ProductID)	
	if err != nil {
		return 0, fmt.Errorf("erro ao obter produto: %v", err)
	}

	// dependendo do tipo, calcula novo estoque
	newStock := currentStock // começa com estoque atual
	switch m.Type {
	case "Entrada":
		newStock += m.Quantity 
	case "Saida":
		newStock -= m.Quantity
	case "Ajuste":
		// ajuste significa que o estoque passa a ser exatamente m.Quantity
		newStock = m.Quantity	
}

// evita estoque negativo em saídas (pode ser configurável futuramente)
//if !s.allowNegative && newStock < 0 {
//	return 0, errors.New("estoque insuficiente para saída")
//}

// agora fazemos a operação dentro de uma transação para garantir atomicidade
tx, err := s.db.BeginTx(ctx, nil)
if err != nil {
	return 0, fmt.Errorf("erro ao iniciar transação: %v", err)
}

// 1) atualizar o estoque na tabela products usando a função injetada (que pode usar tx)
if err := s.updateStock(ctx, m.ProductID, newStock); err != nil {
	_ = tx.Rollback() // tenta rollback em caso de erro
	return 0, fmt.Errorf("erro ao atualizar estoque do produto: %v", err)
}

// 2) inserir o registro de movimento (usando repo, que usa a conexao principal - OK)
id, err := s.repo.Insert(ctx, m)
if err != nil {
	_ = tx.Rollback() // tenta rollback em caso de erro
	return 0, fmt.Errorf("erro ao inserir movimentação de estoque: %v", err)
}

// 3) commit da transação
if err := tx.Commit(); err != nil {
	return 0, fmt.Errorf("erro ao commitar transação: %v", err)
}

return id, nil
}

// GetHistory retorna o historico de movimentações para um produto
func (s *Service) GetHistory(ctx context.Context, productID int) ([]Movement, error) {
	return s.repo.GetByProduct(ctx, productID)
}

// Saida reduz o estoque de um produto
func (s *Service) Saida(ctx context.Context, productID int, quantity float64) error {
	// cria movimento de saída
	m := &Movement{
		ProductID: productID,
		Type: "Saida",
		Quantity: quantity,
	}
	_, err := s.CreateMovement(ctx, m)
	return err
}

func (s *Service) Entrada(ctx context.Context, productID int, quantity float64) error {
	// cria movimento de entrada
	m := &Movement{
		ProductID: productID,
		Type: "Entrada",
		Quantity: quantity,
	}
	_, err := s.CreateMovement(ctx, m)
	return err
	}



