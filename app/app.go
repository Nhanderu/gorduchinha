package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/paemuri/gorduchinha/app/cache"
	"github.com/paemuri/gorduchinha/app/config"
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/data"
	"github.com/paemuri/gorduchinha/app/logger"
)

type App struct {
	Config       config.Config
	Logger       logger.Logger
	DataManager  contract.DataManager
	CacheManager contract.CacheManager
	HTTPClient   *http.Client
}

func New(env string) App {

	cfg, err := config.ReadConfig(env)
	endAsErr(err, "Could not read configuration values.", os.Stdout, os.Stderr)

	log, err := logger.New(
		cfg.App.Name,
		cfg.App.Debug,
	)
	endAsErr(err, "Could not create logging structure.", os.Stdout, os.Stderr)

	db, err := data.Connect(cfg.DB.URL)
	endAsErr(err, "Could not connect to database.", log.InfoWriter(), log.ErrorWriter())

	atInterruption(func() {
		log.Infof("Closing DB Connection.")
		db.Close()
	})

	cache, err := cache.New(
		cfg.Cache.URL,
		cfg.Cache.DB,
		cfg.Cache.Prefix,
		cfg.Cache.DefaultExpiration,
	)
	endAsErr(err, "Could not connect to cache.", log.InfoWriter(), log.ErrorWriter())

	httpClient := &http.Client{Timeout: cfg.HTTPClient.Timeout}

	return App{
		Config:       cfg,
		Logger:       log,
		DataManager:  db,
		CacheManager: cache,
		HTTPClient:   httpClient,
	}
}

func (App) AtInterruption(fn func()) {
	atInterruption(fn)
}

func (App) EndAsErr(err error, message string, infow io.Writer, errorw io.Writer) {
	endAsErr(err, message, infow, errorw)
}

func atInterruption(fn func()) {
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt)
		<-sc

		fn()
		os.Exit(0)
	}()
}

func endAsErr(err error, message string, infow io.Writer, errorw io.Writer) {
	if err != nil {
		fmt.Fprintln(errorw, "Error:", err)
		fmt.Fprintln(infow, message)
		os.Exit(1)
	}
}
