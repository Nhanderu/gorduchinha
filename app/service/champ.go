package service

import (
	"time"

	"github.com/paemuri/gorduchinha/app/contract"
	"github.com/paemuri/gorduchinha/app/entity"
	"github.com/pkg/errors"
)

type champService struct {
	data  contract.DataManager
	cache contract.CacheManager
}

func NewChampService(
	data contract.DataManager,
	cache contract.CacheManager,
) contract.ChampService {

	return champService{
		data:  data,
		cache: cache,
	}
}

func (s champService) Find() ([]entity.Champ, error) {

	cacheKey := "champ"

	var champs []entity.Champ
	err := s.cache.GetJSON(cacheKey, &champs)
	if err != nil {

		champs, err = s.data.Champ().Find()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		s.cache.SetJSON(cacheKey, champs)
		s.cache.SetExpiration(cacheKey, time.Hour*24*30)
	}

	return champs, nil
}

func (s champService) FindBySlug(slug string) (entity.Champ, error) {

	cacheKey := "champ:slug:" + slug

	var champ entity.Champ
	err := s.cache.GetJSON(cacheKey, &champ)
	if err != nil {

		champ, err = s.data.Champ().FindBySlug(slug)
		if err != nil {
			return champ, errors.WithStack(err)
		}

		s.cache.SetJSON(cacheKey, champ)
		s.cache.SetExpiration(cacheKey, time.Hour*24*30)
	}

	return champ, nil
}
