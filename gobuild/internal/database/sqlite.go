package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
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
		return fmt.Errorf("erro ao verificar a conexão com o banco de dados: %v", err)
	}

	DB = db // Atribui a conexão à variável global

	fmt.Println("SQLite conectado com sucesso!")

	return nil
}