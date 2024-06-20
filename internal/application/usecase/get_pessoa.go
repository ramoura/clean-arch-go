package usecase

import (
	"github.com/ramoura/clean-arch-go/internal/application/repository"
	"github.com/ramoura/clean-arch-go/internal/domain"
)

type GetPessoa struct {
	Repo repository.PessoaRepository
}

func (getPessoa *GetPessoa) Execute(id string) (*domain.Pessoa, error) {
	pessoa, err := getPessoa.Repo.GetPessoa(id)
	if err != nil {
		return nil, err
	}
	return pessoa, nil
}
