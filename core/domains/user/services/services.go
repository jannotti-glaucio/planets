package services

import (
	"github.com/jannotti-glaucio/planets/core/domains/user/repositories"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/authenticate"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/create"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/destroy"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/index"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/show"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/update"
	"github.com/jannotti-glaucio/planets/core/tools/providers/hash"
	"github.com/jannotti-glaucio/planets/core/tools/providers/jwt"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
)

type Dependences struct {
	Repository repositories.IUserRepository
	Logger     logger.ILoggerProvider
	Hash       hash.IHashProvider
	Jwt        jwt.IJwtProvider
}

type Services struct {
	Create       create.Service
	Update       update.Service
	Index        index.Service
	Show         show.Service
	Destroy      destroy.Service
	Authenticate authenticate.Service
}

func New(dep Dependences) *Services {
	return &Services{
		Authenticate: authenticate.Service{
			Repository: dep.Repository,
			Hash:       dep.Hash,
			Jwt:        dep.Jwt,
			Logger:     dep.Logger,
		},
		Create: create.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
			Hash:       dep.Hash,
		},
		Update: update.Service{
			Repository: dep.Repository,
			Logger:     dep.Logger,
			Hash:       dep.Hash,
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
