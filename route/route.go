package route

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/router"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/cache"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/database/repository"
	"github.com/hashicorp/go-hclog"
)

type RouteParameters struct {
	AppRouter  router.Router
	Log        hclog.Logger
	Repository repository.Repository
	Cache      cache.Cache
}
