package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/logger"
	"github.com/paemuri/gorduchinha/cmd/api/server/handler"
	"github.com/paemuri/gorduchinha/cmd/api/server/middleware"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func Run(
	port int,
	prefix string,
	authClientsURLs []string,
	rateLimitPeriod time.Duration,
	rateLimitLimit int64,
	routeKeysUpdateTrophies string,
	log logger.Logger,
	cache contract.CacheManager,
	teamService contract.TeamService,
	champService contract.ChampService,
	scraperService contract.ScraperService,
) error {

	router := newRouter(middleware.CORS(authClientsURLs))

	open := router.group(
		prefix,
		middleware.Logger(log),
		middleware.BodyLimit(),
		middleware.CORS(authClientsURLs),
		// middleware.RateLimit(cache, rateLimitPeriod, rateLimitLimit),
	)

	open.handle(http.MethodGet, "/health", handler.HealthCheck())
	open.handle(http.MethodPost, "/graphql", handler.HandleGraphql(teamService, champService))
	open.handle(http.MethodGet, "/teams", handler.ListTeams(teamService))
	open.handle(http.MethodGet, "/teams/{abbr}", handler.FindTeamByAbbr(teamService))
	open.handle(http.MethodGet, "/champs", handler.ListChamps(champService))
	open.handle(http.MethodGet, "/champs/{slug}", handler.FindChampBySlug(champService))
	open.handle(http.MethodPut, "/trophies", handler.UpdateTrophies(scraperService), middleware.QueryKeyValidation(routeKeysUpdateTrophies))

	address := fmt.Sprintf(":%d", port)
	err := fasthttp.ListenAndServe(address, router.requestHandler())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
