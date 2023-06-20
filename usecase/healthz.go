package usecase

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
)

type Healthz interface {
	CheckRepository() error
	CheckCache() error
}

type UseCaseHealthz struct {
	Repository repository.Repository
	Cache      cache.Cache
}

func NewHealthz(repository repository.Repository, cache cache.Cache) Healthz {
	return &UseCaseHealthz{
		Repository: repository,
		Cache:      cache,
	}
}

func (usecaseHealthz *UseCaseHealthz) CheckRepository() error {
	return usecaseHealthz.Repository.Check()
}

func (usecaseHealthz *UseCaseHealthz) CheckCache() error {
	return usecaseHealthz.Cache.Check()
}
