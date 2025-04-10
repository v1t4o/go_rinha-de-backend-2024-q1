package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go_rinha-de-backend-2024-q1/db"

	"github.com/gorilla/mux"
)

type Transaction struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}

type Balance struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type BankStatement struct {
	Saldo             int           `json:"saldo"`
	UltimasTransacoes []Transaction `json:"ultimas_transacoes"`
}

type Customer struct {
	ID     int
	Limite int
	Saldo  int
}

type TransactionDB struct {
	ID          int
	ClienteID   int
	Valor       int
	Tipo        string
	Descricao   string
	RealizadaEm time.Time
}

// CriarUsuario insere um usuário no banco de dados
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	clienteID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
		return
	}

	var transacao Transaction
	err = json.NewDecoder(r.Body).Decode(&transacao)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	if transacao.Valor <= 0 {
		http.Error(w, "Valor da transação inválido", http.StatusUnprocessableEntity)
		return
	}

	if transacao.Tipo != "c" && transacao.Tipo != "d" {
		http.Error(w, "Tipo de transação inválido", http.StatusUnprocessableEntity)
		return
	}

	if len(transacao.Descricao) < 1 || len(transacao.Descricao) > 10 {
		http.Error(w, "Descrição da transação inválida", http.StatusUnprocessableEntity)
		return
	}

	tx, err := db.Connect()
	if err != nil {
		http.Error(w, "Erro ao iniciar transação no banco de dados", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var cliente Customer
	err = tx.QueryRow("SELECT id, limite, saldo FROM clientes WHERE id = $1", clienteID).Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Erro ao buscar cliente no banco de dados", http.StatusInternalServerError)
		return
	}

	novoSaldo := cliente.Saldo
	if transacao.Tipo == "c" {
		novoSaldo += transacao.Valor
	} else if transacao.Tipo == "d" {
		novoSaldo -= transacao.Valor
		if novoSaldo < -cliente.Limite {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}

	_, err = tx.Exec("UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, clienteID)
	if err != nil {
		http.Error(w, "Erro ao atualizar saldo do cliente", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO transacoes (cliente_id, valor, tipo, descricao)
		VALUES ($1, $2, $3, $4)
	`, clienteID, transacao.Valor, transacao.Tipo, transacao.Descricao)
	if err != nil {
		http.Error(w, "Erro ao registrar transação", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Erro ao commitar transação no banco de dados", http.StatusInternalServerError)
		return
	}

	resposta := struct {
		Limite int `json:"limite"`
		Saldo  int `json:"saldo"`
	}{
		Limite: cliente.Limite,
		Saldo:  novoSaldo,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resposta)
}

// BuscarUsuarios traz todos os usuários salvos no banco de dados
func Extract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	clienteID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID do cliente inválido", http.StatusBadRequest)
		return
	}

	var cliente Cliente
	err = db.QueryRow("SELECT id, limite, saldo FROM clientes WHERE id = $1", clienteID).Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Erro ao buscar cliente no banco de dados", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT valor, tipo, descricao, realizada_em
		FROM transacoes
		WHERE cliente_id = $1
		ORDER BY realizada_em DESC
		LIMIT 10
	`, clienteID)
	if err != nil {
		http.Error(w, "Erro ao buscar transações do cliente", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ultimasTransacoes []TransacaoDetalhada
	for rows.Next() {
		var transacao TransacaoDetalhada
		err = rows.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			http.Error(w, "Erro ao ler transação do banco de dados", http.StatusInternalServerError)
			return
		}
		ultimasTransacoes = append(ultimasTransacoes, transacao)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Erro ao iterar sobre as transações", http.StatusInternalServerError)
		return
	}

	extrato := Extrato{
		Saldo: Saldo{
			Total:       cliente.Saldo,
			DataExtrato: time.Now().UTC(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: ultimasTransacoes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(extrato)
}
