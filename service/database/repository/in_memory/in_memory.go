package repository

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/util"
)

type InMemory struct{}

func NewInMemory(config *util.Config) (repository.Repository, error) {
	return &InMemory{}, nil
}

func (inMemory *InMemory) Check() error {
	return nil
}

func (inMemory *InMemory) Close() error {
	return nil
}

func (inMemory *InMemory) Order() repository.Order {
	return NewOrder()
}
