package server

import (
	"net/http"
)

type transaction struct {
	Valor       uint32 `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}

type transactionResult struct {
	Limite uint32 `json:"limite"`
	Saldo  string `json:"saldo"`
}

type extract struct {
	Total       uint32 `json:"total"`
	DataExtrato string `json:"data_extrato"`
	Limite      string `json:"limite"`
}

// CriarUsuario insere um usuário no banco de dados
func CreateTransaction(w http.ResponseWriter, r *http.Request) {

}

// BuscarUsuarios traz todos os usuários salvos no banco de dados
func Extract(w http.ResponseWriter, r *http.Request) {

}
