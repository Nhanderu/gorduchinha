package main

import (
	"os"

	"github.com/paemuri/gorduchinha/app"
	"github.com/paemuri/gorduchinha/cmd/api/server"
)

func main() {

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	app := app.New(env)
	app.Logger.Infof("Running server at port %d.", app.Config.Server.Port)
	err := server.Run(
		app.Config.Server.Port,
		app.Config.Server.Prefix,
		app.Config.Server.Auth.ClientsURLs,
		app.Config.Server.RateLimit.Period,
		app.Config.Server.RateLimit.Limit,
		app.Config.Server.RouteKeys.UpdateTrophies,
		app.Logger,
		app.CacheManager,
		app.Services().NewTeam(),
		app.Services().NewChamp(),
		app.Services().NewScraper(),
	)
	app.EndAsErr(err, "Could not run server.", app.Logger.InfoWriter(), app.Logger.ErrorWriter())
}
