package services

import (
	"github.com/jannotti-glaucio/planets/core/domains/planet/repositories"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/create"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/destroy"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/index"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/show"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/update"
	"github.com/jannotti-glaucio/planets/core/tools/providers/cache"
	client "github.com/jannotti-glaucio/planets/core/tools/providers/http_client"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
)

type Dependences struct {
	Repository repositories.IPlanetRepository
	Logger     logger.ILoggerProvider
	HTTPClient client.IHTTPClientProvider
	Cache      cache.ICacheProvider
}

type Services struct {
	Create  create.Service
	Update  update.Service
	Index   index.Service
	Show    show.Service
	Destroy destroy.Service
	Cache   cache.ICacheProvider
}

func New(dep Dependences) *Services {
	return &Services{
		Create: create.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
			HTTPClient: dep.HTTPClient,
			Cache:      dep.Cache,
		},
		Update: update.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
			HTTPClient: dep.HTTPClient,
			Cache:      dep.Cache,
		},
		Index: index.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
		},
		Show: show.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
		},
		Destroy: destroy.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
		},
	}
}
