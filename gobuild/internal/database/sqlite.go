package database

import (
	"database/sql" // pacotes padrão para manipulação de banco de dados
	"fmt"          //para formtação de strings e erros
	"time"         // para manipulação de tempo

	_ "github.com/mattn/go-sqlite3" // driver SQLite (import por side effect)
)

var DB *sql.DB // Variável global para o banco de dados

// conecta (abre/gera) o arquivo do banco de dados SQLite
func Connect() error {


	//Abre ou cria o arquivo do banco de dados SQLite
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}

	// Verifica a conexão
	if err := db.Ping(); err != nil {
		// fecha o handle antes de retornar o erro (boa prática)
		_ = db.Close()
		return fmt.Errorf("erro ao verificar a conexão com o banco de dados: %v", err)
	}

	DB = db // Atribui a conexão à variável global

	//executa  migrações iniciais (criação de tabelas se não existirem)
	if err := migrate(); err != nil {
	_ = DB.Close() // em caso de erro, fecha a conexão
	DB = nil      // zera a variável global
	return fmt.Errorf("erro ao executar migrações: %v", err)
	}

	// mensagem de sucesso (log/feedback)
	fmt.Println("SQLite conectado com sucesso!")

	return nil
}

// migrate executa SQL de criação de tabelas iniciais
// inclui a tabela products (se ainda não existir) e a nova tabela stock_movements
func migrate() error {
	//schemaProducts para tabela products (matenmos compatibilidade com o que já usado)
	schemaProducts := ` 
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL DEFAULT 0,
		stock REAL NOT NULL DEFAULT 0,
		unit TEXT NOT NULL,
		category TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP	
	);
	`
	// Schema para a tebela stock_movements (nova tabela para histórico de movimentações de estoque)
	// - product_id: referencia ao produto (não há FK forçada aqui para simplicidade)
	// - tipo: "Entrada", "Saida", "ajuste"
	// - quantidade: número (pode ser decimal)
	// - created_at: timestamp automático
	schemaStock := `
	CREATE TABLE IF NOT EXISTS stock_movements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		tipo TEXT NOT NULL,
		quantidade REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	// execução da query de criação da tabela no DB.
	if _, err := DB.Exec(schemaProducts); err != nil {
		return fmt.Errorf("erro ao criar tabela products: %v", err)
	}
	if _, err := DB.Exec(schemaStock); err != nil {			
		return fmt.Errorf("erro ao criar tabela stock_movements: %v", err)
	}

	// exemplo opcional: podemos inserir um registro inicial se quisermos (comentei).
	_ = time.Now() // usado se quisermos logs de timestamp; mantido para referencia futura.

	return nil
}
