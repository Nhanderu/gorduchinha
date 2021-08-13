package data

import (
	"database/sql"

	"github.com/paemuri/gorduchinha/app/contract"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func Connect(url string) (contract.DataManager, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	p := new(pool)
	p.pool = db
	p.repo = repo{db}

	return p, nil
}

type pool struct {
	repo
	pool *sql.DB
}

func (p *pool) Begin() (contract.TransactionManager, error) {

	tx, err := p.pool.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	t := new(transaction)
	t.transaction = tx
	t.repo = repo{tx}

	return t, nil
}

func (p *pool) Close() error {

	err := p.pool.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var _ contract.TransactionManager = &transaction{}

type transaction struct {
	repo
	transaction *sql.Tx
	committed   bool
	rolledback  bool
}

func (t *transaction) Rollback() error {

	if !t.committed && !t.rolledback {

		err := t.transaction.Rollback()
		if err != nil {
			return errors.WithStack(err)
		}

		t.rolledback = true
	}

	return nil
}

func (t *transaction) Commit() error {

	err := t.transaction.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	t.committed = true
	return nil
}

var _ contract.RepoManager = repo{}

type repo struct {
	ex executor
}

func (r repo) Champ() contract.ChampRepo {
	return newChampRepo(r.ex)
}

func (r repo) Team() contract.TeamRepo {
	return newTeamRepo(r.ex)
}

func (r repo) Trophy() contract.TrophyRepo {
	return newTrophyRepo(r.ex)
}
