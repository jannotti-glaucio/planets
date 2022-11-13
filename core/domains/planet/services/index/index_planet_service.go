package index

import (
	"context"

	"github.com/jannotti-glaucio/planets/core/domains/planet/entities"
	"github.com/jannotti-glaucio/planets/core/domains/planet/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
)

//Service ...
type Service struct {
	Repository repositories.IPlanetRepository
	Logger     logger.ILoggerProvider
}

//Execute service responsible for find one register
func (service *Service) Execute(ctx context.Context, filter *map[string]interface{}) (planets entities.Planets, response communication.Response) {
	identifierTracer := "index.planet.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".flter", Value: filter})
	defer span.Finish()

	planets, err := service.Repository.All(ctx, filter)
	comm := communication.New()

	if err != nil {
		service.Logger.Error(ctx, "domain.planet.service.index.index_planet_service.Repository.All", err)
		response = comm.Response(404, "error_list")
		return
	}

	response = comm.Response(200, "success")
	return
}
