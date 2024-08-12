// This file contains the repository implementation layer.
package repository

import (
	"database/sql"
	"errors"
)

type Repository struct {
	exec Executor
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(dbUrl string) *Repository {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
	return &Repository{
		exec: db,
	}
}

func (r *Repository) GetSqlDb() (*sql.DB, error) {
	var (
		sqlDb *sql.DB
		ok    bool
	)

	if sqlDb, ok = r.exec.(*sql.DB); !ok {
		return nil, errors.New("exec is not a sqlDb")
	}

	return sqlDb, nil
}

func (r *Repository) SetExecutor(exec Executor) {
	r.exec = exec
}
