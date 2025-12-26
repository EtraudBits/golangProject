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

// Função construtora
// NewRepository cria um novo repository de budget
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// CreateBudget cria um orçamento com seus itens dentro de uma transação
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
	// Inserindo o orçamento
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

	// Pegando o ID gerado -> será usado nos itens
	budgetID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("erro ao obter ID do ormaçamento: %w", err)
	}

	// Inserindo os Itens
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

	// commit final -> aqui o banco confirma tudo
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("erro ao commitar transação %w", err)
	}

	// retorno final
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
		return nil, fmt.Errorf("erro ao lista orçamentos: %w", err)
	}
	defer rows.Close()

	var budgets []Budget

	//2-> Itera sobre os orçamentos
	for rows.Next() {
		var b Budget

		//Mapeia colunas -> struct
		if err := rows.Scan(
			&b.ID,
			&b.Customer,
			&b.Total,
			&b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler orçamento: %w", err)
		}
		budgets = append(budgets, b)
	}
	// Verifica erros na iteração
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return budgets, nil
}

// ListItemsByBudget retorna os itens de um orçamento específico
func (r *Repository) ListItemsByBudget(
	ctx context.Context,
	budgetID int64,
) ([]BudgetItem, error) {

	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, budget_id, product_id, product, quantity, unit_price, subtotal
		 FROM budget_items
		 WHERE budget_id = ?
		 ORDER BY id`,
		budgetID,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar itens do orçamento: %w", err)
	}
	defer rows.Close()

	var items []BudgetItem

	for rows.Next() {
		var it BudgetItem

		if err := rows.Scan(
			&it.ID,
			&it.BudgetID,
			&it.ProductID,
			&it.Product,
			&it.Quantity,
			&it.UnitPrice,
			&it.Subtotal,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear item do orçamento: %w", err)
		}

		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração dos itens do orçamento: %w", err)
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

func (r *Repository) Cancel(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE budgets SET status = 'CANCELADO' WHERE id = ?`,
		id,
	)
	return err
}

// UpdateBudget atualiza um orçamento e seus itens dentro de uma transação
func (r *Repository) UpdateBudget(
	ctx context.Context,
	budget *Budget,
	items []BudgetItem,
) error {

	// 1-> Inicia a transação
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 2-> Atualiza o cabeçalho do orçamento
	_, err = tx.ExecContext(ctx,
		`UPDATE budgets SET customer = ?, total = ? WHERE id = ?`,
		budget.Customer,
		budget.Total,
		budget.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3-> Remove os itens antigos
	_, err = tx.ExecContext(ctx,
		`DELETE FROM budget_items WHERE budget_id = ?`,
		budget.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4-> Insere os novos itens
	for _, item := range items {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO budget_items (budget_id, product_id, product, quantity, unit_price, subtotal)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			budget.ID,
			item.ProductID,
			item.Product,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 5-> Commit da transação
	return tx.Commit()
}

// DeleteBudget remove um orçamento e seus itens dentro de uma transação
func (r *Repository) DeleteBudget(ctx context.Context, id int64) error {
	// 1-> Inicia a transação
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	// 2-> Deleta os itens do orçamento primeiro (evita lixo no banco)
	_, err = tx.ExecContext(ctx,
		`DELETE FROM budget_items WHERE budget_id = ?`,
		id,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao deletar itens do orçamento: %w", err)
	}

	// 3-> Deleta o orçamento
	result, err := tx.ExecContext(ctx,
		`DELETE FROM budgets WHERE id = ?`,
		id,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao deletar orçamento: %w", err)
	}

	// 4-> Verifica se o orçamento existia (RowsAffected garante que o ID existe)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return sql.ErrNoRows // Orçamento não encontrado
	}

	// 5-> Commit da transação
	return tx.Commit()
}

