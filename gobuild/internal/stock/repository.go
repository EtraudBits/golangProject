package stock

import (
	"context"      // Para passar contexto em operações de banco de dados
	"database/sql" // pacote sql para manipulação de rows/ results
	"fmt"          // para formatação de strings e erros
)

// Repository gerencia operações de banco de dados para movimentações de estoque (stock_movements)
type Repository struct {
	DB *sql.DB // Conexão com o banco (injetada na criação do repositório)
}

// NewRepository cria uma nova instância do repositório de movimentações de estoque (conexão já pronta)
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
// Insert registra um movimento de estoque e retorna o ID inserido.
// Usamos ExecContext para passar contexto e facilitar cancelmento/timeouts.
func (r *Repository) Insert(ctx context.Context, m *Movement) (int64, error) {
	result, err := r.DB.ExecContext(ctx,
		`INSERT INTO stock_movements (product_id, tipo, quantidade) VALUES (?, ?, ?)`,	
		m.ProductID, m.Type, m.Quantity,
	)
	if err != nil {
		return 0, fmt.Errorf("erro ao inserir movimentação de estoque: %v", err)
	}

	id, err := result.LastInsertId() // obtém o ID do novo registro
	if err != nil {
		return 0, fmt.Errorf("erro ao obter ID da movimentação inserida: %v", err)
	}

	return id, nil
}
// GetByProduct retorna historico de movimentos de um produto (ordenado desc por data)
func (r *Repository) GetByProduct(ctx context.Context, productID int) ([]Movement, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, product_id, tipo, quantidade, created_at 
		FROM stock_movements
		WHERE product_id = ?
		ORDER BY created_at DESC`, productID,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar movimentações de estoque: %v", err)
	}
	defer rows.Close() // garante fechamento das rows após uso

	var list []Movement
	for rows.Next() {
		var m Movement
		if err := rows.Scan(&m.ID, &m.ProductID, &m.Type, &m.Quantity, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("erro ao escanear movimentação de estoque: %v", err)
		}
		list = append(list, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das movimentações: %v", err)
	}

	return list, nil
}