package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type PgAdapter struct {
	db *sql.DB
}

func NewPgAdapter(dbUser string, dbPass string, dbAddr string, dbName string) (*PgAdapter, error) {
	db, err := setupDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbAddr, dbName))
	return &PgAdapter{db: db}, err
}

func setupDB(strConn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", strConn)
	if err != nil {
		return nil, err
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)                  // Máximo de conexões abertas
	db.SetMaxIdleConns(25)                  // Máximo de conexões inativas
	db.SetConnMaxLifetime(30 * time.Second) // Tempo máximo de vida útil de uma conexão

	return db, nil
}

func (pga *PgAdapter) Query(statement string, params ...any) (*sql.Rows, error) {
	return pga.db.Query(statement, params...)
}

func (pga *PgAdapter) QueryRow(statement string, params ...any) *sql.Row {
	return pga.db.QueryRow(statement, params...)
}

func (pga *PgAdapter) Exec(statement string, params ...any) (sql.Result, error) {
	return pga.db.Exec(statement, params...)
}

func (pga *PgAdapter) Close() error {
	return pga.db.Close()
}
