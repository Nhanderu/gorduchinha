package resolver

import (
	"github.com/paemuri/gorduchinha/app/entity"
)

type teamResolver struct {
	team entity.Team
}

func NewTeamResolver(team entity.Team) *teamResolver {
	return &teamResolver{
		team: team,
	}
}

func (r teamResolver) Abbr() string {
	return r.team.Abbr
}

func (r teamResolver) Name() string {
	return r.team.Name
}

func (r teamResolver) FullName() string {
	return r.team.FullName
}

func (r teamResolver) Trophies(args *TrophyArgs) []*trophyResolver {

	resolvers := make([]*trophyResolver, 0)
	for _, trophy := range r.team.Trophies {
		if args.ChampSlug == nil || *args.ChampSlug == trophy.Champ.Slug {
			resolvers = append(resolvers, NewTrophyResolver(trophy))
		}
	}

	return resolvers
}

type TrophyArgs struct {
	ChampSlug *string
}
