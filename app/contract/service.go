package contract

import (
	"github.com/paemuri/gorduchinha/app/entity"
)

type TeamService interface {
	Find() ([]entity.Team, error)
	FindByAbbr(abbr string) (entity.Team, error)
}

type ChampService interface {
	Find() ([]entity.Champ, error)
	FindBySlug(slug string) (entity.Champ, error)
}

type ScraperService interface {
	ScrapeAndUpdate() error
}
