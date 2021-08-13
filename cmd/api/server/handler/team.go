package handler

import (
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/cmd/api/server/viewmodel"
	"github.com/valyala/fasthttp"
)

func ListTeams(teamService contract.TeamService) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {

		teams, err := teamService.Find()
		if err != nil {
			HandleError(ctx, err)
			return
		}

		respondOK(ctx, viewmodel.ParseTeamResponseList(teams))
	}
}

func FindTeamByAbbr(teamService contract.TeamService) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {

		abbr, _ := ctx.UserValue("abbr").(string)
		team, err := teamService.FindByAbbr(abbr)
		if err != nil {
			HandleError(ctx, err)
			return
		}

		respondOK(ctx, viewmodel.ParseTeamResponse(team))
	}
}
