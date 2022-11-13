package planets

import (
	"context"

	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb/repositories/planets"
	"github.com/jannotti-glaucio/planets/adapters/rest/util"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/create"
	"github.com/jannotti-glaucio/planets/core/domains/planet/services/update"
	"github.com/jannotti-glaucio/planets/core/tools/providers/cache"
	client "github.com/jannotti-glaucio/planets/core/tools/providers/http_client"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func service(ctx context.Context) *services.Services {
	mongoRepo := planets.Setup(ctx)
	memcacheHost := viper.GetString("memcache.host") + ":" + viper.GetString("memcache.port")
	memcache, _ := cache.New(memcacheHost)

	return services.New(services.Dependences{
		Repository: mongoRepo,
		Logger:     logger.New(),
		HTTPClient: client.New(),
		Cache:      memcache,
	})
}

//Create ...
func Create() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		dto := create.Dto{}
		util.Parser(echoCtx.Request(), &dto)

		span, ctx := util.TraceRestEndpoint(echoCtx, "create-planet")
		defer span.Finish()

		created, response := service(ctx).Create.Execute(ctx, dto)
		if created.UUID != "" {
			response.Data = created.PublicPlanet()
		}

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Index ...
func Index() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		filters := map[string]interface{}{}
		name := echoCtx.QueryParam("name")
		if name != "" {
			filters["name"] = name
		}

		span, ctx := util.TraceRestEndpoint(echoCtx, "index-planet")
		defer span.Finish()

		result, response := service(ctx).Index.Execute(ctx, &filters)
		response.Data = result.PublicPlanets()

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Show ...
func Show() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		uuid := echoCtx.Param("UUID")

		span, ctx := util.TraceRestEndpoint(echoCtx, "show-planet")
		defer span.Finish()

		result, response := service(ctx).Show.Execute(ctx, uuid)
		if result.UUID != "" {
			response.Data = result.PublicPlanet()
		}

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Update ...
func Update() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		dto := update.Dto{}
		dto.UUID = echoCtx.Param("UUID")
		util.Parser(echoCtx.Request(), &dto)

		span, ctx := util.TraceRestEndpoint(echoCtx, "update-planet")
		defer span.Finish()

		result, response := service(ctx).Update.Execute(ctx, dto)
		if result.UUID != "" {
			response.Data = result.PublicPlanet()
		}

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Destroy ...
func Destroy() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		uuid := echoCtx.Param("UUID")

		span, ctx := util.TraceRestEndpoint(echoCtx, "destroy-planet")
		defer span.Finish()

		response := service(ctx).Destroy.Execute(ctx, uuid)

		echoCtx.JSON(response.Status, response)
		return nil
	}
}
