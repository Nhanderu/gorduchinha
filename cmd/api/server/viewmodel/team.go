package viewmodel

import (
	"github.com/paemuri/gorduchinha/app/entity"
)

type TeamResponse struct {
	Abbr     string               `json:"abbr"`
	Name     string               `json:"name"`
	FullName string               `json:"full_name"`
	Trophies []TeamResponseTrophy `json:"trophies"`
}

type TeamResponseTrophy struct {
	Year  int `json:"year"`
	Champ struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
	} `json:"champ"`
}

func ParseTeamResponseList(teams []entity.Team) []TeamResponse {

	vm := make([]TeamResponse, len(teams))
	for i := range vm {
		vm[i] = ParseTeamResponse(teams[i])
	}

	return vm
}

func ParseTeamResponse(team entity.Team) TeamResponse {

	trophiesVM := make([]TeamResponseTrophy, len(team.Trophies))
	for i := range trophiesVM {
		trophiesVM[i].Year = team.Trophies[i].Year
		trophiesVM[i].Champ.Slug = team.Trophies[i].Champ.Slug
		trophiesVM[i].Champ.Name = team.Trophies[i].Champ.Name
	}

	return TeamResponse{
		Abbr:     team.Abbr,
		Name:     team.Name,
		FullName: team.FullName,
		Trophies: trophiesVM,
	}
}
