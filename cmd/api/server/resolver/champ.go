package resolver

import (
	"github.com/paemuri/gorduchinha/app/entity"
)

type champResolver struct {
	champ entity.Champ
}

func NewChampResolver(champ entity.Champ) *champResolver {
	return &champResolver{
		champ: champ,
	}
}

func (r champResolver) Slug() string {
	return r.champ.Slug
}

func (r champResolver) Name() string {
	return r.champ.Name
}
