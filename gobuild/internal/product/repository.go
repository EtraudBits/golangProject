package product

import (
	"context"      // Para passar contexto em operações de banco de dados
	"database/sql" // pacote sql para manipulação de rows/ results
	"fmt"          // para formatação de strings e erros
)

type Repository struct {
	DB *sql.DB // Conexão com o banco (injetada na criação do repositório)
}

// NewRepository cria uma nova instância do repositório de produtos
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// create insere um novo produto no banco de dados e retorna o ID inserido.
func (r *Repository) Create(ctx context.Context, p *Produto) (int64, error) {
	// Query INSERT com Placeholders (compativel com SQLite)
	result, err := r.DB.ExecContext(ctx,
 `INSERT INTO products (name, price, stock, unit, category, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		p.Name, p.Preco, p.Estoque, p.Unidade, p.Categoria, &p.DataCriacao)
	if err != nil {
		return 0, fmt.Errorf("erro ao inserir produto: %v", err)
	}

	id, err := result.LastInsertId() // obtém o ID do novo registro
	if err != nil {
		return 0, fmt.Errorf("erro ao obter ID do produto inserido: %v", err)
	}

	return id, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]Produto, error) {
// executa a query SELECT para buscar todos os produtos
rows, err := r.DB.QueryContext(ctx,  `SELECT id, name, price, stock, unit, category, created_at FROM products`)
if err != nil {
	return nil, fmt.Errorf("erro ao buscar produtos: %v", err)
}

defer rows.Close() // garante que as rows serão fechadas após o uso

var produtos []Produto
// itera sobre os resultados e faz o scan em structs Produto
for rows.Next() {

	var p Produto
	if err := rows.Scan(&p.ID, &p.Name, &p.Preco, &p.Estoque, &p.Unidade, &p.Categoria, &p.DataCriacao); err != nil {
		return nil, fmt.Errorf("erro ao escanear produto: %v", err)
	}
	produtos = append(produtos, p)
}
// checa por erros na iteração
if err := rows.Err(); err != nil {
	return nil, fmt.Errorf("erro durante iteração dos produtos: %v", err)
}
return produtos, nil
}
// GetByID busca um produto pelo seu ID (chave primária)
func (r *Repository) GetByID(ctx context.Context, id int) (*Produto, error) {

	row := r.DB.QueryRowContext(ctx, `SELECT id, name, price, stock, unit, category, created_at FROM products WHERE id = ?`, id)

	var p Produto
	if err := row.Scan(&p.ID, &p.Name, &p.Preco, &p.Estoque, &p.Unidade, &p.Categoria, &p.DataCriacao); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // produto não encontrado
		}
		return nil, fmt.Errorf("erro ao escanear produto: %v", err)
	}
	return &p, nil
}

// Update atualiza os dados de um produto existente no banco de dados (atualiza pelo ID)
func (r *Repository) Update(ctx context.Context, p *Produto) error {

	_, err := r.DB.ExecContext(ctx,
		`UPDATE products SET name = ?, price = ?, stock = ?, unit = ?, category = ? WHERE id = ?`,
		p.Name, p.Preco, p.Estoque, p.Unidade, p.Categoria, p.ID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %v", err)
	}
	return nil
}

func (r *Repository) Delete (ctx context.Context, id int) error {
	// executa a query DELETE para remover o produto pelo ID
	_, err := r.DB.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id	)	
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %v", err)
	}
	return nil	
}

