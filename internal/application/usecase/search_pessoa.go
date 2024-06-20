package usecase

import (
	"github.com/ramoura/clean-arch-go/internal/application/repository"
	"github.com/ramoura/clean-arch-go/internal/domain"
)

type SearchPessoa struct {
	Repo repository.PessoaRepository
}

func (searchPessoa *SearchPessoa) Execute(term string) ([]*domain.Pessoa, error) {
	pessoa, err := searchPessoa.Repo.SearchPessoaWithTerm(term)
	if err != nil {
		return nil, err
	}
	return pessoa, nil
}
