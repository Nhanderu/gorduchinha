package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/paemuri/gorduchinha/app/constant"
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/cmd/api/server/resolver"
	"github.com/paemuri/gorduchinha/cmd/api/server/viewmodel"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func HandleGraphql(teamService contract.TeamService, champService contract.ChampService) func(*fasthttp.RequestCtx) {

	typedefs, _ := ioutil.ReadFile("static/graphql/schema.gql")
	queryResolver := resolver.NewQueryResolver(teamService, champService)
	schema := graphql.MustParseSchema(string(typedefs), queryResolver)

	return func(ctx *fasthttp.RequestCtx) {

		var request viewmodel.GraphQLQueryRequest
		err := json.Unmarshal(ctx.PostBody(), &request)
		if err != nil {
			HandleError(ctx, errors.WithStack(constant.NewErrorInvalidRequestBody()))
			return
		}

		response := schema.Exec(ctx, request.Query, request.OperationName, request.Variables)
		respondJSON(ctx, http.StatusOK, response)
	}
}
