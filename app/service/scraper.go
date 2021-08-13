package service

import (
	"bytes"
	"net/http"
	"regexp"
	"strconv"

	"github.com/paemuri/gorduchinha/app/constant"
	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/entity"
	"github.com/paemuri/gorduchinha/app/logger"
	"github.com/andybalholm/cascadia"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type scraperService struct {
	data         contract.DataManager
	log          logger.Logger
	httpClient   *http.Client
	teamService  contract.TeamService
	champService contract.ChampService
}

func NewScraperService(
	data contract.DataManager,
	log logger.Logger,
	httpClient *http.Client,
	teamService contract.TeamService,
	champService contract.ChampService,
) contract.ScraperService {

	return scraperService{
		data:         data,
		log:          log,
		httpClient:   httpClient,
		teamService:  teamService,
		champService: champService,
	}
}

func (s scraperService) ScrapeAndUpdate() error {

	teams, err := s.scrapeAll()
	if err != nil {
		return errors.WithStack(err)
	}

	tx, err := s.data.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer tx.Rollback()

	err = tx.Trophy().Delete()
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Trophy().BulkInsertByTeams(teams)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type scraperFn func() (map[string][]int, error)

func (s scraperService) scrapeAll() ([]entity.Team, error) {

	scrapers := map[string]scraperFn{
		constant.ChampSlugNationalLeague1Div:    s.scrapeNationalLeague1Div,
		constant.ChampSlugNationalLeague2Div:    s.scrapeNationalLeague2Div,
		constant.ChampSlugNationalCup:           s.scrapeNationalCup,
		constant.ChampSlugWorldCup:              s.scrapeWorldCup,
		constant.ChampSlugIntercontinentalCup:   s.scrapeIntercontinentalCup,
		constant.ChampSlugSouthAmericanCupA:     s.scrapeSouthAmericanCupA,
		constant.ChampSlugSouthAmericanCupB:     s.scrapeSouthAmericanCupB,
		constant.ChampSlugSouthAmericanSupercup: s.scrapeSouthAmericanSupercup,
		constant.ChampSlugSPStateCup:            s.scrapeSPStateCup,
		constant.ChampSlugRJStateCup:            s.scrapeRJStateCup,
		constant.ChampSlugRSStateCup:            s.scrapeRSStateCup,
		constant.ChampSlugMGStateCup:            s.scrapeMGStateCup,
	}

	allTrophies := make(map[string][]entity.Trophy)
	for champSlug, scraper := range scrapers {

		champ, err := s.champService.FindBySlug(champSlug)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		champTrophies, err := scraper()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for team, years := range champTrophies {
			for _, year := range years {
				allTrophies[team] = append(allTrophies[team], entity.Trophy{
					Year:  year,
					Champ: champ,
				})
			}
		}

	}

	teams := make([]entity.Team, 0)
	for teamAbbr, trophies := range allTrophies {

		team, err := s.teamService.FindByAbbr(teamAbbr)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		team.Trophies = trophies
		teams = append(teams, team)
	}

	return teams, nil
}

var (
	reYears  = regexp.MustCompile(`\d{4}`)
	mapTeams = map[string]string{
		"Corinthians":      constant.TeamAbbrSCCP,
		"Palmeiras":        constant.TeamAbbrSEP,
		"São Paulo":        constant.TeamAbbrSPFC,
		"Santos":           constant.TeamAbbrSFC,
		"Flamengo":         constant.TeamAbbrCRF,
		"Vasco da Gama":    constant.TeamAbbrCRVG,
		"Vasco":            constant.TeamAbbrCRVG,
		"Fluminense":       constant.TeamAbbrFFC,
		"Botafogo":         constant.TeamAbbrBFR,
		"Atlético Mineiro": constant.TeamAbbrCAM,
		"Cruzeiro":         constant.TeamAbbrCEC,
		"Grêmio":           constant.TeamAbbrGFBPA,
		"Internacional":    constant.TeamAbbrSCI,
	}
	allTeamsAbbrs = []string{
		constant.TeamAbbrSCCP,
		constant.TeamAbbrSEP,
		constant.TeamAbbrSPFC,
		constant.TeamAbbrSFC,
		constant.TeamAbbrCRF,
		constant.TeamAbbrCRVG,
		constant.TeamAbbrCRVG,
		constant.TeamAbbrFFC,
		constant.TeamAbbrBFR,
		constant.TeamAbbrCAM,
		constant.TeamAbbrCEC,
		constant.TeamAbbrGFBPA,
		constant.TeamAbbrSCI,
	}
)

func (s scraperService) scrape(
	url string,
	linesSel, teamSel, yearsSel cascadia.Selector,
	possibleTeamsAbbrs ...string,
) (map[string][]int, error) {

	if len(possibleTeamsAbbrs) == 0 {
		possibleTeamsAbbrs = allTeamsAbbrs
	}

	possibleTeams := make(map[string]struct{}, len(possibleTeamsAbbrs))
	for _, abbr := range possibleTeamsAbbrs {
		possibleTeams[abbr] = struct{}{}
	}

	res, err := s.httpClient.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	trophies := make(map[string][]int)
	lines := linesSel.MatchAll(doc)
	for _, line := range lines {

		rawTeam := teamSel.MatchFirst(line)
		if rawTeam == nil {
			continue
		}

		teamName := innerText(rawTeam)
		teamAbbr, found := mapTeams[teamName]
		if !found {
			continue
		}

		_, possible := possibleTeams[teamAbbr]
		if !possible {
			continue
		}

		rawYears := innerText(yearsSel.MatchFirst(line))
		years := reYears.FindAllString(rawYears, -1)
		if len(years) < 1 {
			continue
		}

		teamTrophies := make([]int, len(years))
		for i := range teamTrophies {
			teamTrophies[i], err = strconv.Atoi(years[i])
			if err != nil {
				s.log.Errorf(
					"Error scraping title for %s: %s.",
					teamAbbr,
					err.Error(),
				)
				continue
			}
		}

		trophies[teamAbbr] = teamTrophies
	}

	return trophies, nil
}

func (s scraperService) scrapeNationalLeague1Div() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Lista_de_campe%C3%B5es_do_Campeonato_Brasileiro_de_Futebol"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Por_clube) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > span > a:last-child")
		years = cascadia.MustCompile("td:nth-child(2)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeNationalLeague2Div() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Campeonato_Brasileiro_de_Futebol_-_S%C3%A9rie_B"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Títulos_por_clube) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a:last-child")
		years = cascadia.MustCompile("td:nth-child(2)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeNationalCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Lista_de_campe%C3%B5es_da_Copa_do_Brasil_de_Futebol"
	)

	var (
		lines = cascadia.MustCompile("h2:has(#Resultados_por_clube) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a:last-child")
		years = cascadia.MustCompile("td:nth-child(4)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeWorldCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Copa_do_Mundo_de_Clubes_da_FIFA"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Por_clube) ~ table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a:last-child")
		years = cascadia.MustCompile("td:nth-child(2)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeIntercontinentalCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Copa_Intercontinental"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Por_clube) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a:last-child")
		years = cascadia.MustCompile("td:nth-child(2)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeSouthAmericanCupA() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Lista_de_campe%C3%B5es_da_Copa_Libertadores_da_Am%C3%A9rica"
	)

	var (
		lines = cascadia.MustCompile("h2:has(#Títulos_e_vice_por_equipe) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a")
		years = cascadia.MustCompile("td:nth-child(3)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeSouthAmericanCupB() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Copa_Sul-Americana"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Títulos_por_clube) ~ table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a")
		years = cascadia.MustCompile("td:nth-child(3)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeSouthAmericanSupercup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Recopa_Sul-Americana"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Por_equipe) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a")
		years = cascadia.MustCompile("td:nth-child(3)")
	)

	trophies, err := s.scrape(url, lines, team, years)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeSPStateCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Campeonato_Paulista_de_Futebol"
	)

	var (
		lines = cascadia.MustCompile("h4:has(#Por_clube) + table > tbody > tr:not(:first-child)")
		team  = cascadia.MustCompile("td:first-child > b > a:last-child")
		years = cascadia.MustCompile("td:nth-child(3)")
	)

	trophies, err := s.scrape(
		url,
		lines,
		team,
		years,
		constant.TeamAbbrSCCP,
		constant.TeamAbbrSEP,
		constant.TeamAbbrSPFC,
		constant.TeamAbbrSFC,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeRJStateCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Campeonato_Carioca_de_Futebol"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Títulos_por_clube) ~ table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > a")
		years = cascadia.MustCompile("td:nth-child(2)")
	)

	trophies, err := s.scrape(
		url,
		lines,
		team,
		years,
		constant.TeamAbbrCRF,
		constant.TeamAbbrCRVG,
		constant.TeamAbbrFFC,
		constant.TeamAbbrBFR,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeRSStateCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Campeonato_Ga%C3%BAcho_de_Futebol"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Por_clube) + table > tbody > tr:not(:first-child)")
		team  = cascadia.MustCompile("td:first-child > b > a:last-child")
		years = cascadia.MustCompile("td:nth-child(3)")
	)

	trophies, err := s.scrape(
		url,
		lines,
		team,
		years,
		constant.TeamAbbrGFBPA,
		constant.TeamAbbrSCI,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func (s scraperService) scrapeMGStateCup() (map[string][]int, error) {

	const (
		url = "https://pt.wikipedia.org/wiki/Campeonato_Mineiro_de_Futebol"
	)

	var (
		lines = cascadia.MustCompile("h3:has(#Por_equipe) + table > tbody > tr")
		team  = cascadia.MustCompile("td:first-child > b > a")
		years = cascadia.MustCompile("td:nth-child(3)")
	)

	trophies, err := s.scrape(
		url,
		lines,
		team,
		years,
		constant.TeamAbbrCAM,
		constant.TeamAbbrCEC,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return trophies, nil
}

func innerText(node *html.Node) string {

	if node.FirstChild == nil {
		return node.Data
	}

	buffer := bytes.NewBufferString("")
	child := node.FirstChild
	for {

		if child.FirstChild == nil {
			buffer.WriteString(child.Data)
		} else {
			buffer.WriteString(innerText(child))
		}

		if child.NextSibling == nil || child == node.LastChild {
			break
		}

		child = child.NextSibling
	}

	return buffer.String()
}
