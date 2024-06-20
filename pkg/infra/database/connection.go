package database

import "database/sql"

type Connection interface {
	Query(statement string, params ...any) (*sql.Rows, error)
	QueryRow(statement string, params ...any) *sql.Row
	Exec(statement string, params ...any) (sql.Result, error)
	Close() error
}
