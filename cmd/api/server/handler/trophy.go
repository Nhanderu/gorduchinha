package handler

import (
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/valyala/fasthttp"
)

func UpdateTrophies(scraperService contract.ScraperService) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {

		err := scraperService.ScrapeAndUpdate()
		if err != nil {
			HandleError(ctx, err)
			return
		}

		respondOK(ctx, nil)
	}
}
