package resolver

import (
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/pkg/errors"
)

type queryResolver struct {
	teamService  contract.TeamService
	champService contract.ChampService
}

func NewQueryResolver(teamService contract.TeamService, champService contract.ChampService) *queryResolver {
	return &queryResolver{
		teamService:  teamService,
		champService: champService,
	}
}

type TeamArgs struct {
	Abbr string
}

func (r queryResolver) Team(args *TeamArgs) (*teamResolver, error) {

	team, err := r.teamService.FindByAbbr(args.Abbr)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return NewTeamResolver(team), nil
}

func (r queryResolver) Teams() ([]*teamResolver, error) {

	teams, err := r.teamService.Find()
	if err != nil {
		return nil, errors.Cause(err)
	}

	resolvers := make([]*teamResolver, len(teams))
	for i := range teams {
		resolvers[i] = &teamResolver{
			team: teams[i],
		}
	}

	return resolvers, nil
}

type ChampArgs struct {
	Slug string
}

func (r queryResolver) Champ(args *ChampArgs) (*champResolver, error) {

	champ, err := r.champService.FindBySlug(args.Slug)
	if err != nil {
		return nil, errors.Cause(err)
	}

	return NewChampResolver(champ), nil
}

func (r queryResolver) Champs() ([]*champResolver, error) {

	champs, err := r.champService.Find()
	if err != nil {
		return nil, errors.Cause(err)
	}

	resolvers := make([]*champResolver, len(champs))
	for i := range champs {
		resolvers[i] = &champResolver{
			champ: champs[i],
		}
	}

	return resolvers, nil
}
