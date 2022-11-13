package index

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
func (service *Service) Execute(ctx context.Context) (users entities.Users, response communication.Response) {
	identifierTracer := "index.user.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer)
	defer span.Finish()

	users, err := service.Repository.All(ctx)
	comm := communication.New()

	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.index.index_user_service.Repository.All", err)
		response = comm.Response(404, "error_list")
		return
	}

	response = comm.Response(200, "success")
	return
}
