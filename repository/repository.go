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

	sqlDb := &SqlDb{
		db: db,
	}

	return &Repository{
		exec: sqlDb,
	}
}

func (r *Repository) GetSqlDb() (SqlDbInterface, error) {
	var (
		sqlDb *SqlDb
		ok    bool
	)

	if sqlDb, ok = r.exec.(*SqlDb); !ok {
		return nil, errors.New("exec is not a sqlDb")
	}

	return sqlDb, nil
}

func (r *Repository) SetExecutor(exec Executor) {
	r.exec = exec
}
