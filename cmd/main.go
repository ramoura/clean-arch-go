package main

import (
	"github.com/ramoura/clean-arch-go/internal/application/usecase"
	"github.com/ramoura/clean-arch-go/internal/infra/repository"
	"github.com/ramoura/clean-arch-go/internal/infra/web"
	"github.com/ramoura/clean-arch-go/pkg/infra/database"
	"github.com/ramoura/clean-arch-go/pkg/infra/http"
	"log"
)

func main() {
	connection, err := database.NewPgAdapter("root", "1234", "localhost:5432", "rinhadb")
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	pessoaRepo := repository.PessoaRepositoryDatabase{
		Connection: connection,
	}

	createPessoa := usecase.CreatePessoa{
		Repo: pessoaRepo,
	}

	getPessoa := usecase.GetPessoa{
		Repo: pessoaRepo,
	}

	searchPessoa := usecase.SearchPessoa{
		Repo: pessoaRepo,
	}
	controller := web.ControllerPessoas{
		CreatePessoa: createPessoa,
		GetPessoa:    getPessoa,
		SearchPessoa: searchPessoa,
	}

	httpServer := http.NewHttpServerMux()

	controller.RegisterRoutes(httpServer)
	httpServer.Start()

}
