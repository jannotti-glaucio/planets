package show

import (
	"context"

	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/domains/user/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
)

//Service ...
type Service struct {
	Repository repositories.IUserRepository
	Logger     logger.ILoggerProvider
}

//Execute service responsible for find one register
func (service *Service) Execute(ctx context.Context, uuid string) (user entities.User, response communication.Response) {
	identifierTracer := "show.user.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	user, err := service.Repository.FindByUUID(ctx, uuid)
	comm := communication.New()

	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.show.show_user_service.Repository.FindByUUID", err)
		response = comm.Response(404, "error_list")
		return
	}

	response = comm.Response(200, "success")
	return
}
