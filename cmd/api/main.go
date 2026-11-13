package main

import (
	"log"

	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/database"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/env"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/handler"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/loteria"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/repository"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/server"
	"github.com/maiconDeSouza/lotofacil_lab_v2/internal/service"
)

func main() {
	env := env.InitENV()
	collections, errapp := database.NewConnectMongo(env)
	if errapp != nil {
		log.Fatalf("%s : %v", errapp.Message, errapp.Err)
	}

	loteria := loteria.NewLoteriaClient(env)
	repo := repository.NewRepository(collections.CollectionContest, collections.CollectionGame)
	service := service.NewService(repo, *loteria)
	handler := handler.NewHandler(service)
	server := server.NewServer(handler, env)

	server.Run()

	// env := env.InitENV()
	// loteria := loteria.NewLoteriaClient(env)

	// res, _ := loteria.GetFrequentNumbers()

	// fmt.Println(res)

}
