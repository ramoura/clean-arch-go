package usecase

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ramoura/clean-arch-go/internal/application/repository"
	"github.com/ramoura/clean-arch-go/internal/domain"
	"log"
)

type CreatePessoa struct {
	Repo repository.PessoaRepository
}

func (createPessoa *CreatePessoa) Execute(pessoa *domain.Pessoa) error {

	if !pessoa.IsValidDate() {
		return errors.New("data de nascimento inválida. Use o formato AAAA-MM-DD")
	}

	if !pessoa.IsValid() {
		return errors.New("dados inválidos.")
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	pessoa.ID = newUUID.String()

	err = createPessoa.Repo.CreatePessoa(pessoa)
	if err != nil {
		return err
	}

	return nil

}
