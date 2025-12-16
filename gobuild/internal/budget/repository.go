package budget

// Importações necessárias para o funcionamento do repositório
import (
	"context"      // Controle de tempo/cancelamento
	"database/sql" // API padrão do Go para banco
	"fmt"
)

// Repository lida exclusivamente com SQL do modulo budget
type Repository struct {
	DB *sql.DB
}
//Função construtora
// NewRepository cria um novo repository de budget
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

//CreateBudget cria um orçamento com seus itens dentro de uma transação
func (r *Repository) CreateBudget(
	ctx context.Context, //Contexto da requisição
	budget *Budget, //cabeçalho
	items []BudgetItem, //itens do orçamento
) (int64, error) {
	//iniciando a transação
	tx, err := r.DB.BeginTx(ctx, nil) //a partir daqui tudo usa tx, se falhar rollback
	if err != nil {
		return 0, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
//Inserindo o orçamento
result, err := tx.ExecContext(ctx,
	`INSERT INTO budgets (customer, total)
	VALUES (?, ?)`,
	budget.Customer,
	budget.Total,
)
if err != nil {
	tx.Rollback()
	return 0, fmt.Errorf("erro ao inserir orçamento: %w", err)
}

//Pegando o ID gerado -> será usado nos itens
budgetID, err := result.LastInsertId()
if err != nil {
	tx.Rollback()
	return 0, fmt.Errorf("erro ao obter ID do ormaçamento: %w", err)
}

//Inserindo os Itens
for _, item := range items {
	_, err := tx.ExecContext(ctx,
	`INSERT INTO budget_items (budget_id, product_id, product, quantity, unit_price, subtotal)
	VALUES(?, ?, ?, ?, ?, ?)`,
	budgetID,
	item.ProductID, 
	item.Product, 
	item.Quantity,
	item.UnitPrice, 
	item.Subtotal,
)
if err != nil {
	tx.Rollback()
	return 0, fmt.Errorf("erro ao inserir item do orçamento: %w", err)
	}
}

//commit final -> aqui o banco confirma tudo
if err := tx.Commit(); err != nil {
	return 0, fmt.Errorf("erro ao commitar transação %w", err)
}

//retorno final
return budgetID, nil
}