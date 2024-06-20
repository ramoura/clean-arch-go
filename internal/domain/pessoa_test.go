package domain

import (
	"testing"
)

func TestPessoaIsValid(t *testing.T) {
	tests := []struct {
		pessoa   Pessoa
		expected bool
	}{
		{Pessoa{Apelido: "apelido", Nome: "nome", Nascimento: "2000-01-01", Stack: []string{"Go"}}, true},
		{Pessoa{Apelido: "", Nome: "nome", Nascimento: "2000-01-01", Stack: []string{"Go"}}, false},
		{Pessoa{Apelido: "apelido", Nome: "", Nascimento: "2000-01-01", Stack: []string{"Go"}}, false},
		{Pessoa{Apelido: "apelido", Nome: "nome", Nascimento: "", Stack: []string{"Go"}}, false},
		{Pessoa{Apelido: "apelido", Nome: "nome", Nascimento: "2000-01-01", Stack: []string{}}, false},
	}

	for _, test := range tests {
		result := test.pessoa.IsValid()
		if result != test.expected {
			t.Errorf("For Pessoa %+v, expected %v, got %v", test.pessoa, test.expected, result)
		}
	}
}

func TestPessoaIsValidDate(t *testing.T) {
	tests := []struct {
		pessoa   Pessoa
		expected bool
	}{
		{Pessoa{Nascimento: "2000-01-01"}, true},
		{Pessoa{Nascimento: "2000-13-01"}, false},
		{Pessoa{Nascimento: "2000-01-32"}, false},
		{Pessoa{Nascimento: "invalid-date"}, false},
		{Pessoa{Nascimento: ""}, false},
	}

	for _, test := range tests {
		result := test.pessoa.IsValidDate()
		if result != test.expected {
			t.Errorf("For Pessoa %+v, expected %v, got %v", test.pessoa, test.expected, result)
		}
	}
}
