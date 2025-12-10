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

// migrate executa as migrações iniciais do banco de dados
// aqui criamos as tabelas necessárias se elas não existirem
func migrate() error {
	// Definição do shema SQL para a tabela products.
	//id: PK autoincrement, name: texto, price e stock: Real (números com ponto),
	//unit: unidade (m2, kg, un...), category: texto, created_at: timestamp.
	schema := `
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
	// execução da query de criação da tabela no DB.
	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("erro ao criar tabela products: %v", err)
	}

	// exemplo opcional: podemos inserir um registro inicial se quisermos (comentei).
	_ = time.Now() // usado se quisermos logs de timestamp; mantido para referencia futura.

	return nil
}
