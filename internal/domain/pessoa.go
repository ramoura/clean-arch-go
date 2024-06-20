package domain

import (
	"time"
)

type Pessoa struct {
	ID         string
	Apelido    string
	Nome       string
	Nascimento string
	Stack      []string
}

func (pessoa Pessoa) IsValid() bool {
	return !(pessoa.Apelido == "" || pessoa.Nome == "" || pessoa.Nascimento == "" || len(pessoa.Stack) == 0)
}

func (pessoa Pessoa) IsValidDate() bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, pessoa.Nascimento)
	return err == nil
}
