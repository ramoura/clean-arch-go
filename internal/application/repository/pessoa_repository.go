package repository

import "github.com/ramoura/clean-arch-go/internal/domain"

type PessoaRepository interface {
	CreatePessoa(Pessoa *domain.Pessoa) error
	GetPessoa(id string) (*domain.Pessoa, error)
	SearchPessoaWithTerm(term string) ([]*domain.Pessoa, error)
}
