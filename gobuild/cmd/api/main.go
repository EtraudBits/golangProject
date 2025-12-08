package main

import (
	"log"

	"github.com/EtraudBits/golangProject/gobuild/internal/database"
	"github.com/EtraudBits/golangProject/gobuild/internal/server"
)

func main() {
// Conecta ao banco de dados
	err := database.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// criar instancia do servidor
	s := server.New()

	// registrar rotas
	s.RegisterRoutes()

	// iniciar servidor http
	s.Start()

}	