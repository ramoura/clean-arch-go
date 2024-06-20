package integration

import (
	"bytes"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/ramoura/clean-arch-go/internal/application/usecase"
	"github.com/ramoura/clean-arch-go/internal/domain"
	"github.com/ramoura/clean-arch-go/internal/infra/repository"
	"github.com/ramoura/clean-arch-go/internal/infra/web"
	"github.com/ramoura/clean-arch-go/pkg/infra/database"
	http2 "github.com/ramoura/clean-arch-go/pkg/infra/http"
	"github.com/ramoura/clean-arch-go/test"
	"net/http"
	"net/http/httptest"

	"os"
	"testing"
)

var db database.Connection
var httpServer http2.Server

func TestMain(m *testing.M) {
	testDatabase := test.SetupTestDatabase()
	db = testDatabase.Connection

	pessoaRepo := repository.PessoaRepositoryDatabase{
		Connection: testDatabase.Connection,
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

	httpServer = http2.NewHttpServerGin()
	controller.RegisterRoutes(httpServer)

	// Run tests
	code := m.Run()

	// Cleanup
	testDatabase.TearDown()

	os.Exit(code)
}

func TestCreatePessoa(t *testing.T) {
	t.Cleanup(func() {
		db.Exec("DELETE FROM pessoas")
	})
	pessoa := web.PessoaRequest{
		Apelido:    "Test",
		Nome:       "Pessoa Teste",
		Nascimento: "1985-01-01",
		Stack:      []string{"Go", "Java"},
	}

	body, _ := json.Marshal(pessoa)
	req := httptest.NewRequest("POST", "/pessoas", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	httpServer.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}

	location := w.Header().Get("Location")
	if location == "" {
		t.Error("Expected Location header to be set")
	}

}

func TestGetPessoa(t *testing.T) {
	t.Cleanup(func() {
		db.Exec("DELETE FROM pessoas")
	})

	pessoa := web.PessoaRequest{
		Apelido:    "Test",
		Nome:       "Pessoa Teste",
		Nascimento: "1985-01-01",
		Stack:      []string{"Go", "Java"},
	}

	req, w := createPessoa(t, pessoa)

	location := w.Header().Get("Location")

	req = httptest.NewRequest("GET", location, nil)
	w = httptest.NewRecorder()
	httpServer.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var retrieved domain.Pessoa
	if err := json.NewDecoder(w.Body).Decode(&retrieved); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if retrieved.Nome != pessoa.Nome {
		t.Errorf("Expected Nome %s, got %s", pessoa.Nome, retrieved.Nome)
	}

	if retrieved.Apelido != pessoa.Apelido {
		t.Errorf("Expected Apelido %s, got %s", pessoa.Apelido, retrieved.Apelido)
	}

}

func createPessoa(t *testing.T, pessoa web.PessoaRequest) (*http.Request, *httptest.ResponseRecorder) {
	// Primeiro, crie uma pessoa para poder recuperá-la

	body, _ := json.Marshal(pessoa)
	req := httptest.NewRequest("POST", "/pessoas", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	httpServer.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}
	return req, w
}

func TestSearchPessoa(t *testing.T) {
	t.Cleanup(func() {
		db.Exec("DELETE FROM pessoas")
	})
	createPessoa(t, web.PessoaRequest{
		Apelido:    "Test 1",
		Nome:       "Pessoa Teste 1",
		Nascimento: "1985-01-01",
		Stack:      []string{"Go", "Java"},
	})
	createPessoa(t, web.PessoaRequest{
		Apelido:    "Test 2",
		Nome:       "Pessoa Teste 2",
		Nascimento: "1985-01-01",
		Stack:      []string{"Python,Java,Ruby"},
	})
	createPessoa(t, web.PessoaRequest{
		Apelido:    "Test 3",
		Nome:       "Pessoa Teste 3",
		Nascimento: "1985-01-01",
		Stack:      []string{"Python,Ruby"},
	})

	req := httptest.NewRequest("GET", "/pessoas?t=java", nil)
	w := httptest.NewRecorder()
	httpServer.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var results []domain.Pessoa
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}
