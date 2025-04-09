package main

import (
	"go_rinha-de-backend-2024-q1/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/clientes/{id}/transacoes", server.CreateTransaction).Methods(http.MethodPost)
	router.HandleFunc("/clientes/{id}/extrato", server.Extract).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":5000", router))
}
