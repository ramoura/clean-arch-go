package repository

import (
	"github.com/ramoura/clean-arch-go/internal/domain"
	"github.com/ramoura/clean-arch-go/pkg/infra/database"
	"strings"
)

type PessoaRepositoryDatabase struct {
	Connection database.Connection
}

func (repo PessoaRepositoryDatabase) CreatePessoa(pessoa *domain.Pessoa) error {
	stackString := strings.Join(pessoa.Stack, ",")
	if _, err := repo.Connection.Exec("INSERT INTO pessoas (ID, APELIDO, NOME, NASCIMENTO, STACK) VALUES ($1, $2, $3, $4, $5)", pessoa.ID, pessoa.Apelido, pessoa.Nome, pessoa.Nascimento, stackString); err != nil {
		return err
	}

	return nil
}

func (repo PessoaRepositoryDatabase) GetPessoa(id string) (*domain.Pessoa, error) {
	var pessoa domain.Pessoa
	var stackString string
	row := repo.Connection.QueryRow("SELECT ID, APELIDO, NOME, NASCIMENTO, STACK FROM pessoas WHERE ID = $1", id)
	if err := row.Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &pessoa.Nascimento, &stackString); err != nil {
		return nil, err
	}

	pessoa.Stack = strings.Split(stackString, ",")
	return &pessoa, nil
}

func (repo PessoaRepositoryDatabase) SearchPessoaWithTerm(term string) ([]*domain.Pessoa, error) {
	rows, err := repo.Connection.Query("SELECT ID, APELIDO, NOME, NASCIMENTO, STACK FROM pessoas WHERE BUSCA_TRGM ilike $1 limit 50", "%"+term+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pessoas = []*domain.Pessoa{}
	for rows.Next() {
		var pessoa domain.Pessoa
		var stackString string
		if err := rows.Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &pessoa.Nascimento, &stackString); err != nil {
			return nil, err
		}
		pessoa.Stack = strings.Split(stackString, ",")
		pessoas = append(pessoas, &pessoa)
	}
	return pessoas, nil
}
