package data

import (
	"database/sql"

	"github.com/paemuri/gorduchinha/app/constant"
	"github.com/pkg/errors"
)

type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func parseError(err error, entity string) error {
	if err == nil {
		return nil
	}

	switch errors.Cause(err) {
	case sql.ErrNoRows:
		return errors.WithStack(constant.NewErrorEntityNotFound(entity))
	}

	return errors.WithStack(err)
}
