package users

import (
	"context"

	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb/repositories/users"
	"github.com/jannotti-glaucio/planets/adapters/rest/util"
	"github.com/jannotti-glaucio/planets/core/domains/user/services"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/authenticate"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/create"
	"github.com/jannotti-glaucio/planets/core/domains/user/services/update"
	"github.com/jannotti-glaucio/planets/core/tools/providers/hash"
	"github.com/jannotti-glaucio/planets/core/tools/providers/jwt"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func service(ctx context.Context) *services.Services {
	mongoRepo := users.Setup(ctx)

	return services.New(services.Dependences{
		Repository: mongoRepo,
		Logger:     logger.New(),
		Hash:       hash.New(),
		Jwt:        jwt.New(viper.GetString("jwt.secretKey")),
	})
}

//Auth ...
func Auth() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		dto := authenticate.Dto{}
		util.Parser(echoCtx.Request(), &dto)

		span, ctx := util.TraceRestEndpoint(echoCtx, "auth-user")
		defer span.Finish()

		token, response := service(ctx).Authenticate.Execute(ctx, dto)
		if token != "" {
			response.Data = token
		}

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Create ...
func Create() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		dto := create.Dto{}
		util.Parser(echoCtx.Request(), &dto)

		span, ctx := util.TraceRestEndpoint(echoCtx, "create-user")
		defer span.Finish()

		created, response := service(ctx).Create.Execute(ctx, dto)
		if created.UUID != "" {
			response.Data = created.PublicUser()
		}

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Index ...
func Index() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		span, ctx := util.TraceRestEndpoint(echoCtx, "index-user")
		defer span.Finish()

		result, response := service(ctx).Index.Execute(ctx)
		response.Data = result.PublicUsers()

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Show ...
func Show() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		uuid := echoCtx.Param("UUID")

		span, ctx := util.TraceRestEndpoint(echoCtx, "show-user")
		defer span.Finish()

		result, response := service(ctx).Show.Execute(ctx, uuid)
		if result.UUID != "" {
			response.Data = result.PublicUser()
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

		span, ctx := util.TraceRestEndpoint(echoCtx, "update-user")
		defer span.Finish()

		result, response := service(ctx).Update.Execute(ctx, dto)
		if result.UUID != "" {
			response.Data = result.PublicUser()
		}

		echoCtx.JSON(response.Status, response)
		return nil
	}
}

//Destroy ...
func Destroy() echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		uuid := echoCtx.Param("UUID")

		span, ctx := util.TraceRestEndpoint(echoCtx, "destroy-user")
		defer span.Finish()

		response := service(ctx).Destroy.Execute(ctx, uuid)

		echoCtx.JSON(response.Status, response)
		return nil
	}
}
