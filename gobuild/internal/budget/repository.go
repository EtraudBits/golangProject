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

// ListBudgets retorna todos os orçamentos com seus itens
func (r *Repository) ListBudgets(ctx context.Context) ([]Budget, error) {

	// 1-> Busca todos os orçamentos (cabeçalho)

	rows, err := r.DB.QueryContext(ctx,
	`SELECT id, customer, total, created_at
	FROM budgets
	ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao lista orçamentos: %w", err)}
		defer rows.Close()
		
		var budgets []Budget

		//2-> Itera sobre os orçamentos
		for rows.Next() {
			var b Budget

			// Lê os dados do orçamento
			if err := rows.Scan(&b.ID, &b.Customer, &b.Total, &b.CreatedAt); err != nil {
				return nil, fmt.Errorf("erro ao ler orçamento: %w, err")
			}

			// 3-> Buca os itens deste orçamento
			items, err := r.getItemsByBudget(ctx, b.ID)
			if err != nil {
				return nil, err
			}

			b.Items = items
			budgets = append(budgets, b)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
		return budgets, nil
	}
	
	// getItemsBybudget busca os itens de um orçamento específico
	func (r *Repository) getItemsByBudget(ctx context.Context, budgetID int64) ([]BudgetItem, error) {

		rows, err := r.DB.QueryContext(ctx,
		`SELECT product_id, product, quantity, unit_price, subtotal
		FROM budget_items
		WHERE budget_id = ?`,
		budgetID,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar intens do orçamento: %w", err)
		}
		defer rows.Close()

		var items []BudgetItem

		for rows.Next() {
			var it BudgetItem

			if err := rows.Scan(
				&it.ProductID,
				&it.Product,
				&it.Quantity,
				&it.UnitPrice,
				&it.Subtotal,
			); err != nil {
				return nil, fmt.Errorf("erro ao ler item do orçamento: %w", err)
			}

			items = append(items, it)
		}
		return items, nil
	}

	// GetByID busca um orçamento pelo ID junto com seus itens
	func (r *Repository) GetByID(ctx context.Context, id int64) (*Budget, error) {

		// 1-> Busca o orçamento (cabeçalho)
		row := r.DB.QueryRowContext(ctx,
		`SELECT id, customer, total, created_at
		FROM budgets
		WHERE id = ?`,
		id,
		)

		var b Budget
		if err := row.Scan(&b.ID, &b.Customer, &b.Total, &b.CreatedAt); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil // Orçamento não encontrado
			}
			return nil, err
		}

		// 2-> Busca os itens do orçamento
		rows, err := r.DB.QueryContext(ctx,
		`SELECT id, budget_id, product_id, product, quantity, unit_price, subtotal
		FROM budget_items
		WHERE budget_id = ?`, 
		b.ID,
		)
		
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var item BudgetItem
			if err := rows.Scan(
				&item.ID,
				&item.BudgetID,
				&item.ProductID,
				&item.Product,
				&item.Quantity,
				&item.UnitPrice,
				&item.Subtotal,
			); err != nil {
				return nil, err
			}
			b.Items = append(b.Items, item)
		}
		return &b, nil
	}