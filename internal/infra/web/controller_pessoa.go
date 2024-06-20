package web

import (
	"encoding/json"
	"github.com/ramoura/clean-arch-go/internal/application/usecase"
	"github.com/ramoura/clean-arch-go/internal/domain"
	http2 "github.com/ramoura/clean-arch-go/pkg/infra/http"
	"log"
	"net/http"
)

type PessoaRequest struct {
	ID         string   `json:"id"`
	Apelido    string   `json:"apelido"`
	Nome       string   `json:"nome"`
	Nascimento string   `json:"nascimento"`
	Stack      []string `json:"stack"`
}

type ControllerPessoas struct {
	CreatePessoa usecase.CreatePessoa
	GetPessoa    usecase.GetPessoa
	SearchPessoa usecase.SearchPessoa
}

func (controller *ControllerPessoas) RegisterRoutes(server http2.Server) {
	server.HandlerFunc("POST", "/pessoas", createPessoa(controller.CreatePessoa))
	server.HandlerFunc("GET", "/pessoas/{id}", getPessoa(controller.GetPessoa))
	server.HandlerFunc("GET", "/pessoas", searchPessoa(controller.SearchPessoa))
	server.HandlerFunc("GET", "/contagem-pessoas", contagemPessoas())
}

func searchPessoa(pessoa usecase.SearchPessoa) func(http2.Request) (error, *http2.Response) {
	return func(request http2.Request) (error, *http2.Response) {
		term := request.QueryValue("t")
		if term == "" {
			return nil, &http2.Response{
				Status: http.StatusBadRequest,
			}
		}

		pessoas, err := pessoa.Execute(term)
		if err != nil {
			log.Print(err.Error())
			return err, nil
		}

		resp := http2.Response{
			Status: http.StatusOK,
			Body:   pessoas,
		}
		return nil, &resp
	}
}

func contagemPessoas() func(http2.Request) (error, *http2.Response) {
	return func(request http2.Request) (error, *http2.Response) {
		resp := http2.Response{
			Status: http.StatusOK,
			Body:   "1000",
		}
		return nil, &resp
	}
}

func getPessoa(getPessoa usecase.GetPessoa) func(http2.Request) (error, *http2.Response) {
	return func(request http2.Request) (error, *http2.Response) {
		id := request.PathValue("id")

		pessoa, err := getPessoa.Execute(id)
		if err != nil {
			log.Print(err.Error())
			return err, nil
		}

		resp := http2.Response{
			Status: http.StatusOK,
			Body:   pessoa,
		}
		return nil, &resp

	}

}

func createPessoa(createPessoa usecase.CreatePessoa) func(http2.Request) (error, *http2.Response) {
	return func(request http2.Request) (error, *http2.Response) {
		var p PessoaRequest
		if err := json.NewDecoder(request.Body()).Decode(&p); err != nil {
			log.Print(err.Error())
			return err, nil
		}
		pessoa := domain.Pessoa{
			Apelido:    p.Apelido,
			Nome:       p.Nome,
			Nascimento: p.Nascimento,
			Stack:      p.Stack,
		}

		err := createPessoa.Execute(&pessoa)
		if err != nil {
			log.Print(err.Error())
			return err, nil
		}
		resp := http2.Response{
			Status: http.StatusCreated,
			Headers: map[string]string{
				"Location": "/pessoas/" + pessoa.ID,
			},
		}
		return nil, &resp
	}
}

//func getPessoa(db *sql.DB) web.HandlerFunc {
//	return func(w web.ResponseWriter, r *web.Request) {
//		id := r.PathValue("id")
//
//		var p PessoaRequest
//		var stackString string
//		row := db.QueryRow("SELECT ID, APELIDO, NOME, NASCIMENTO, STACK FROM pessoas WHERE ID = $1", id)
//		if err := row.Scan(&p.ID, &p.Apelido, &p.Nome, &p.Nascimento, &stackString); err != nil {
//			web.Error(w, err.Error(), web.StatusInternalServerError)
//			return
//		}
//
//		p.Stack = strings.Split(stackString, ",")
//		if err := json.NewEncoder(w).Encode(p); err != nil {
//			web.Error(w, err.Error(), web.StatusInternalServerError)
//			return
//		}
//	}
//}
//
//func searchPessoa(db *sql.DB) web.HandlerFunc {
//	return func(w web.ResponseWriter, r *web.Request) {
//
//		term := r.URL.Query().Get("t")
//		if term == "" {
//			w.WriteHeader(web.StatusBadRequest)
//			return
//		}
//
//		rows, err := db.Query("SELECT ID, APELIDO, NOME, NASCIMENTO, STACK FROM pessoas WHERE BUSCA_TRGM = $1 limit 50", "%"+term+"%")
//		if err != nil {
//			web.Error(w, err.Error(), web.StatusInternalServerError)
//			return
//		}
//		defer rows.Close()
//
//		var ps = []PessoaRequest{}
//		for rows.Next() {
//			var p PessoaRequest
//			var stackString string
//			if err := rows.Scan(&p.ID, &p.Apelido, &p.Nome, &p.Nascimento, &stackString); err != nil {
//				web.Error(w, err.Error(), web.StatusInternalServerError)
//				return
//			}
//			p.Stack = strings.Split(stackString, ",")
//			ps = append(ps, p)
//		}
//
//		if err := rows.Err(); err != nil {
//			web.Error(w, err.Error(), web.StatusInternalServerError)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(web.StatusOK)
//		if err := json.NewEncoder(w).Encode(ps); err != nil {
//			fmt.Println("Erro no encode de pessoa:" + err.Error())
//			web.Error(w, err.Error(), web.StatusInternalServerError)
//			return
//		}
//	}
//}
//func contagemPessoas(db *sql.DB) web.HandlerFunc {
//	return func(w web.ResponseWriter, r *web.Request) {
//		_, err := w.Write([]byte("1000"))
//		if err != nil {
//			web.Error(w, err.Error(), web.StatusInternalServerError)
//			return
//		}
//	}
//}
