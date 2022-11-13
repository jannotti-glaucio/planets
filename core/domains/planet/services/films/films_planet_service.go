package films

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jannotti-glaucio/planets/core/tools/providers/cache"
	client "github.com/jannotti-glaucio/planets/core/tools/providers/http_client"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
)

//Service ...
type Service struct {
	Logger     logger.ILoggerProvider
	HTTPClient client.IHTTPClientProvider
	Cache      cache.ICacheProvider
}

//ResponseAPI ...
type ResponseAPI struct {
	Results []Result `json:"results"`
}

//Result ...
type Result struct {
	Films []string `json:"films"`
}

//Execute service
func (service *Service) Execute(ctx context.Context, name string) int {
	identifierTracer := "films.planet.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".name", Value: name})
	defer span.Finish()

	var quantityFilms int
	if service.Cache.Get(ctx, name, &quantityFilms) == nil {
		return quantityFilms
	}

	apiBase := "https://swapi.dev/api/planets/?search="
	response, err := service.HTTPClient.Get(apiBase + name)
	if err != nil {
		service.Logger.Error(ctx, "domain.movies.service.movies_planet_service.HTTPClient.Get", err)
		return 0
	}

	var films ResponseAPI
	err = json.NewDecoder(response.Body).Decode(&films)
	if err != nil {
		service.Logger.Error(ctx, "domain.movies.service.movies_planet_service.json.NewDecoder.Decode", err)
		return 0
	}

	if len(films.Results) > 0 {
		for _, film := range films.Results {
			quantityFilms += len(film.Films)
		}

		cacheExpireTime := time.Minute * 20
		service.Cache.WithExpiration(cacheExpireTime).Set(ctx, name, quantityFilms)
	}

	return quantityFilms
}
