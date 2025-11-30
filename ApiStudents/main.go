package main

import (
	"log"

	"github.com/EtraudBits/golangProject/ApiStudents/api"
)

func main() {
  //Chama todas funções do diretorio API
  server := api.NewServer ()
  
  server.ConfigureRoutes()

  if err:= server.Start(); err != nil {
    log.Fatal(err)
  }

}