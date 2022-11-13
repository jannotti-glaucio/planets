package update

import (
	"context"
	"encoding/json"

	"github.com/jannotti-glaucio/planets/core/domains/planet/entities"
	"github.com/jannotti-glaucio/planets/core/domains/planet/repositories"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/films"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/cache"
	client "github.com/jannotti-glaucio/planets/core/tools/providers/http_client"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
	"github.com/jannotti-glaucio/planets/core/tools/validations"
)

//Dto object receiver
type Dto struct {
	UUID    string `json:"uuid" validate:"isRequired"`
	Name    string `json:"name" validate:"isRequired"`
	Terrain string `json:"terrain" validate:"isRequired"`
	Climate string `json:"climate" validate:"isRequired"`
}

//Service ...
type Service struct {
	Repository repositories.IPlanetRepository
	Logger     logger.ILoggerProvider
	HTTPClient client.IHTTPClientProvider
	Cache      cache.ICacheProvider
}

//Execute responsável por atualizar registros
func (service *Service) Execute(ctx context.Context, dto Dto) (updated entities.Planet, response communication.Response) {
	identifierTracer := "update.planet.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".dto", Value: dto})
	defer span.Finish()

	response.Fields = validations.ValidateStruct(&dto, "")
	comm := communication.New()

	//Check exists planet with this identifier
	planet, err := service.Repository.FindByUUID(ctx, dto.UUID)
	if err != nil {
		service.Logger.Error(ctx, "domain.planet.service.update.update_planet_service.Repository.FindByUUID", err)
		response = comm.Response(500, "error_update")
		return
	}

	planetByName, err := service.Repository.FindByName(ctx, dto.Name)
	if err != nil {
		service.Logger.Info(ctx, "domain.planet.service.update.update_planet_service.Repository.All", err)
	}

	//Check planet already exists
	if planetByName.UUID != "" && planetByName.UUID != planet.UUID {
		response.Fields = append(response.Fields, comm.Fields("name", "already_exists"))
	}

	if planet.UUID == "" {
		response.Fields = append(response.Fields, comm.Fields("uuid", "validate_invalid"))
	}

	if len(response.Fields) > 0 {
		service.Logger.Info(ctx, "domain.planet.service.update.update_planet_service.ValidationError")
		resp := comm.Response(400, "validate_failed")
		resp.Fields = response.Fields
		response = resp
		return
	}

	//Mergin entity and DTO
	toMerge, _ := json.Marshal(dto)
	json.Unmarshal(toMerge, &planet)

	filmsService := films.Service{
		Logger:     service.Logger,
		HTTPClient: service.HTTPClient,
		Cache:      service.Cache,
	}
	planet.Films = filmsService.Execute(ctx, planet.Name)

	updated, err = service.Repository.Save(ctx, planet)

	if err != nil {
		service.Logger.Error(ctx, "domain.planet.service.update.update_planet_service.Repository.Save", err)
		response = comm.Response(500, "error_update")
		return
	}

	response = comm.Response(200, "success")
	return
}
