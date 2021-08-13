package data

import (
	"database/sql"
	"fmt"

	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/entity"
	"github.com/pkg/errors"
)

type teamRepo struct {
	ex           executor
	entity       string
	selectFields string
}

func newTeamRepo(ex executor) contract.TeamRepo {
	return teamRepo{
		ex:     ex,
		entity: "team",
		selectFields: `
			c.id
			, c.created_at
			, c.updated_at
			, c.abbr
			, c.name
		`,
	}
}

func (r teamRepo) parseEntities(rows *sql.Rows, err error) ([]entity.Team, error) {
	if err != nil {
		return nil, errors.WithStack(err)
	}

	teams := make([]entity.Team, 0)
	for rows.Next() {

		team, err := r.parseEntity(rows)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (r teamRepo) parseEntity(s scanner) (entity.Team, error) {

	var team entity.Team
	err := s.Scan(
		&team.ID,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.Abbr,
		&team.Name,
	)
	if err != nil {
		return entity.Team{}, errors.WithStack(err)
	}

	return team, nil
}

func (r teamRepo) Find() ([]entity.Team, error) {
	const query = `
		SELECT %s
			FROM tb_team AS c
			WHERE c.deleted_at IS NULL
		;
	`

	q := fmt.Sprintf(query, r.selectFields)
	teams, err := r.parseEntities(r.ex.Query(q))
	if err != nil {
		return nil, errors.WithStack(parseError(err, r.entity))
	}

	return teams, nil
}

func (r teamRepo) FindByAbbr(abbr string) (entity.Team, error) {
	const query = `
		SELECT %s
			FROM tb_team AS c
			WHERE c.deleted_at IS NULL
				AND c.abbr = $1
		;
	`

	q := fmt.Sprintf(query, r.selectFields)
	team, err := r.parseEntity(r.ex.QueryRow(
		q,
		abbr,
	))
	if err != nil {
		return entity.Team{}, errors.WithStack(parseError(err, r.entity))
	}

	return team, nil
}
