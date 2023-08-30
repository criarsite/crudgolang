package main

import (
	"database/sql"
	"fmt" // Formatar

	//"log"      // Monitoramento
	"net/http" // Web
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexaoBD() (conexao *sql.DB) {

	stringConexao := "adm:123456@/alunos?charset=utf8&parseTime=True&loc=Local"

	conexao, err := sql.Open("mysql", stringConexao)
	if err != nil {
		panic(err.Error())
	}

	return conexao
}

var layout = template.Must(template.ParseGlob("layout/*"))

func main() {

	//ROTAS

	http.HandleFunc("/", Inicio)
	http.HandleFunc("/show", Show)

	http.HandleFunc("/criar", Cadastrar)
	http.HandleFunc("/inserir", Inserir)

	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/atualizar", Atualizar)

	http.HandleFunc("/deletar", Deletar)

	fmt.Println("Servidor Conectado")
	http.ListenAndServe(":5000", nil)
}

func Show(w http.ResponseWriter, r *http.Request) {
	contexto := conexaoBD()
	nId := r.URL.Query().Get("id")
	usuario, err := contexto.Query("SELECT * FROM alunos WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	aluno := Aluno{}

	for usuario.Next() {
		var id int
		var nome, email string
		err = usuario.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}

		aluno.ID = id
		aluno.Nome = nome
		aluno.Email = email

	}

	layout.ExecuteTemplate(w, "show", aluno)
	//defer contexto.Close()
}

func Atualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nome := r.FormValue("nome")
		email := r.FormValue("email")

		contexto := conexaoBD()

		atualizar, err := contexto.Prepare(" UPDATE alunos SET nome=?,email=? WHERE id=? ")
		if err != nil {
			panic(err.Error())
		}

		atualizar.Exec(nome, email, id)
		http.Redirect(w, r, "/", 301)

	}
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idResult := r.URL.Query().Get("id")
	//fmt.Println(idResult)

	contexto := conexaoBD()
	registro, err := contexto.Query("SELECT * FROM alunos WHERE id=?", idResult)

	aluno := Aluno{}
	for registro.Next() {
		var id int
		var nome, email string
		err = registro.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}
		aluno.ID = id
		aluno.Nome = nome
		aluno.Email = email
	}

	//	fmt.Println(aluno)
	layout.ExecuteTemplate(w, "editar", aluno)

}

func Deletar(w http.ResponseWriter, r *http.Request) {
	idResult := r.URL.Query().Get("id")
	//fmt.Println(idResult)

	contexto := conexaoBD()
	remover, err := contexto.Prepare("DELETE FROM alunos WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	remover.Exec(idResult)
	http.Redirect(w, r, "/", 301)
}

func Inserir(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nome := r.FormValue("nome")
		email := r.FormValue("email")

		contexto := conexaoBD()
		cadastrar, err := contexto.Prepare("INSERT INTO alunos(Nome, Email) VALUES(?,?) ")
		if err != nil {
			panic(err.Error())
		}

		cadastrar.Exec(nome, email)

		http.Redirect(w, r, "/", 301)

	}
}

func Cadastrar(w http.ResponseWriter, r *http.Request) {
	layout.ExecuteTemplate(w, "criar", nil)
}

type Aluno struct {
	ID    int
	Nome  string
	Email string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	contexto := conexaoBD()
	listar, err := contexto.Query("SELECT * FROM alunos")
	if err != nil {
		panic(err.Error())
	}

	aluno := Aluno{}
	listaAlunos := []Aluno{}

	for listar.Next() {
		var id int
		var nome, email string
		err = listar.Scan(&id, &nome, &email)
		if err != nil {
			panic(err.Error())
		}

		aluno.ID = id
		aluno.Nome = nome
		aluno.Email = email

		listaAlunos = append(listaAlunos, aluno)
	}
	//fmt.Println(listaAlunos)

	//layout.ExecuteTemplate(w, "inicio", nil)
	layout.ExecuteTemplate(w, "inicio", listaAlunos)
}
