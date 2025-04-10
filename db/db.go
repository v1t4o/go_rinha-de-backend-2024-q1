package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Conectar abre a conexão com o banco de dados
func Connect() (*sql.DB, error) {
	configConnect := "user=golang password=golang dbname=rinha sslmode=disable"

	db, erro := sql.Open("postgres", configConnect)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
