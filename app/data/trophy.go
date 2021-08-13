package data

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/entity"
	"github.com/pkg/errors"
)

type trophyRepo struct {
	ex           executor
	entity       string
	selectFields string
}

func newTrophyRepo(ex executor) contract.TrophyRepo {
	return trophyRepo{
		ex:     ex,
		entity: "trophy",
		selectFields: `
			t.id
			, t.created_at
			, t.updated_at
			, t.uuid
			, t.year
			, c.id
			, c.created_at
			, c.updated_at
			, c.slug
			, c.name
		`,
	}
}

func (r trophyRepo) parseEntities(rows *sql.Rows, err error) ([]entity.Trophy, error) {
	if err != nil {
		return nil, errors.WithStack(err)
	}

	trophies := make([]entity.Trophy, 0)
	for rows.Next() {

		trophy, err := r.parseEntity(rows)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		trophies = append(trophies, trophy)
	}

	return trophies, nil
}

func (r trophyRepo) parseEntity(s scanner) (entity.Trophy, error) {

	var trophy entity.Trophy
	err := s.Scan(
		&trophy.ID,
		&trophy.CreatedAt,
		&trophy.UpdatedAt,
		&trophy.UUID,
		&trophy.Year,
		&trophy.Champ.ID,
		&trophy.Champ.CreatedAt,
		&trophy.Champ.UpdatedAt,
		&trophy.Champ.Slug,
		&trophy.Champ.Name,
	)
	if err != nil {
		return entity.Trophy{}, errors.WithStack(err)
	}

	return trophy, nil
}

func (r trophyRepo) FindByTeamID(teamID uint32) ([]entity.Trophy, error) {
	const query = `
		SELECT %s
			FROM tb_trophy AS t
			JOIN tb_champ AS c
				ON c.deleted_at IS NULL
				AND t.champ_id = c.id
			WHERE t.deleted_at IS NULL
				AND t.team_id = $1
		;
	`

	q := fmt.Sprintf(query, r.selectFields)
	trophies, err := r.parseEntities(r.ex.Query(
		q,
		teamID,
	))
	if err != nil {
		return nil, errors.WithStack(parseError(err, r.entity))
	}

	return trophies, nil
}

func (r trophyRepo) BulkInsertByTeams(teams []entity.Team) error {
	const query = `
		INSERT INTO tb_trophy
			( uuid
			, year
			, champ_id
			, team_id
			)
			VALUES %s
		;
	`
	const value = `(UUID_GENERATE_V4(), $%d, $%d, $%d),`

	var count int
	params := []interface{}{}
	values := bytes.NewBuffer(nil)

	for i := range teams {
		for j := range teams[i].Trophies {
			params = append(
				params,
				teams[i].Trophies[j].Year,
				teams[i].Trophies[j].Champ.ID,
				teams[i].ID,
			)
			position := count * 3
			fmt.Fprintf(values, value, position+1, position+2, position+3)
			count++
		}
	}

	if count == 0 {
		return nil
	}

	q := fmt.Sprintf(query, strings.TrimSuffix(values.String(), ","))
	_, err := r.ex.Exec(q, params...)
	if err != nil {
		return errors.WithStack(parseError(err, r.entity))
	}

	return nil
}

func (r trophyRepo) Delete() error {
	const query = `
		DELETE FROM tb_trophy
			WHERE deleted_at IS NULL
		;
	`

	_, err := r.ex.Exec(query)
	if err != nil {
		return errors.WithStack(parseError(err, r.entity))
	}

	return nil
}
