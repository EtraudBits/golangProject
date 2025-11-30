package main

import (
	"github.com/rs/zerolog/log"

	"github.com/EtraudBits/golangProject/ApiStudents/api"
)

func main() {
  //Chama todas funções do diretorio API
  server := api.NewServer ()
  
  server.ConfigureRoutes()

  if err:= server.Start(); err != nil {
    log.Fatal().Err(err).Msg("Failed to start server: %s") //usando log.Fatal para sinalizar o erro
  }

}