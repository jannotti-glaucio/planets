package show

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
func (service *Service) Execute(ctx context.Context, uuid string) (planet entities.Planet, response communication.Response) {
	identifierTracer := "show.planet.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	planet, err := service.Repository.FindByUUID(ctx, uuid)
	comm := communication.New()

	if err != nil {
		service.Logger.Error(ctx, "domain.planet.service.show.show_planet_service.Repository.FindByUUID", err)
		response = comm.Response(404, "error_list")
		return
	}

	response = comm.Response(200, "success")
	return
}
