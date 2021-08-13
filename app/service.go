package app

import (
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/service"
)

func (app App) Services() AppService {
	return AppService{
		app: app,
	}
}

type AppService struct {
	app App
}

func (svc AppService) NewTeam() contract.TeamService {
	return service.NewTeamService(svc.app.DataManager, svc.app.CacheManager)
}

func (svc AppService) NewChamp() contract.ChampService {
	return service.NewChampService(svc.app.DataManager, svc.app.CacheManager)
}

func (svc AppService) NewScraper() contract.ScraperService {
	return service.NewScraperService(
		svc.app.DataManager,
		svc.app.Logger,
		svc.app.HTTPClient,
		svc.NewTeam(),
		svc.NewChamp(),
	)
}
