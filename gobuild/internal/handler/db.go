package handler

import (
	"net/http"

	"github.com/EtraudBits/golangProject/gobuild/internal/database"

	"github.com/labstack/echo/v4"
)

//TestDBHandler testa a conexão com o banco de dados SQLite
//retorna 200 OK se o Ping no banco for bem-sucedido
func TestDBHandler (c echo.Context) error {

	// testa a conexão com o banco de dados
	err := database.DB.Ping()
	if err != nil {
		// se falhar, retorna erro 500
		return c.String(http.StatusInternalServerError, "Erro ao conectar ao banco de dados")
	}
	// se sucesso, retorna mensagem de sucesso
	return c.String(http.StatusOK, "Conexão com o banco de dados bem-sucedida")
}