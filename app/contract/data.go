package contract

import (
	"github.com/paemuri/gorduchinha/app/entity"
)

type DataManager interface {
	RepoManager
	Begin() (TransactionManager, error)
	Close() error
}

type TransactionManager interface {
	RepoManager
	Rollback() error
	Commit() error
}

type RepoManager interface {
	Champ() ChampRepo
	Team() TeamRepo
	Trophy() TrophyRepo
}

type ChampRepo interface {
	Find() ([]entity.Champ, error)
	FindBySlug(slug string) (entity.Champ, error)
}

type TeamRepo interface {
	Find() ([]entity.Team, error)
	FindByAbbr(abbr string) (entity.Team, error)
}

type TrophyRepo interface {
	FindByTeamID(teamID uint32) ([]entity.Trophy, error)
	BulkInsertByTeams(teams []entity.Team) error
	Delete() error
}
