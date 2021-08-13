package main

import (
	"os"

	"github.com/paemuri/gorduchinha/app"
)

func main() {

	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	app := app.New(env)
	err := app.Services().NewScraper().ScrapeAndUpdate()
	app.EndAsErr(err, "Could not execute service.", app.Logger.InfoWriter(), app.Logger.ErrorWriter())
}
