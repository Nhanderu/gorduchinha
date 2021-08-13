package service

import (
	"time"

	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/entity"
	"github.com/pkg/errors"
)

type teamService struct {
	data  contract.DataManager
	cache contract.CacheManager
}

func NewTeamService(
	data contract.DataManager,
	cache contract.CacheManager,
) contract.TeamService {

	return teamService{
		data:  data,
		cache: cache,
	}
}

func (s teamService) Find() ([]entity.Team, error) {

	cacheKey := "team"

	var teams []entity.Team
	err := s.cache.GetJSON(cacheKey, &teams)
	if err != nil {

		teams, err = s.data.Team().Find()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for i := range teams {
			teams[i].Trophies, err = s.data.Trophy().FindByTeamID(teams[i].ID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}

		s.cache.SetJSON(cacheKey, teams)
		s.cache.SetExpiration(cacheKey, time.Hour*24*30)
	}

	return teams, nil
}

func (s teamService) FindByAbbr(abbr string) (entity.Team, error) {

	cacheKey := "team:abbr:" + abbr

	var team entity.Team
	err := s.cache.GetJSON(cacheKey, &team)
	if err != nil {

		team, err = s.data.Team().FindByAbbr(abbr)
		if err != nil {
			return team, errors.WithStack(err)
		}

		team.Trophies, err = s.data.Trophy().FindByTeamID(team.ID)
		if err != nil {
			return team, errors.WithStack(err)
		}

		s.cache.SetJSON(cacheKey, team)
		s.cache.SetExpiration(cacheKey, time.Hour*24*30)
	}

	return team, nil
}
