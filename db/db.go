package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Conectar abre a conexão com o banco de dados
func Connect() (*sql.DB, error) {
	configConnect := "golang:golang@/rinha?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("pq", configConnect)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
