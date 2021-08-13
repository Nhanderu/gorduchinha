package resolver

import (
	"github.com/paemuri/gorduchinha/app/entity"
)

type trophyResolver struct {
	trophy entity.Trophy
}

func NewTrophyResolver(trophy entity.Trophy) *trophyResolver {
	return &trophyResolver{
		trophy: trophy,
	}
}

func (r trophyResolver) Year() int32 {
	return int32(r.trophy.Year)
}

func (r trophyResolver) Champ() *champResolver {
	return NewChampResolver(r.trophy.Champ)
}
