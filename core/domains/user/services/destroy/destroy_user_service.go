package destroy

import (
	"context"

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

//Execute responsible for deleting register
func (service *Service) Execute(ctx context.Context, uuid string) (response communication.Response) {
	identifierTracer := "destroy.user.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".uuid", Value: uuid})
	defer span.Finish()

	comm := communication.New()

	if uuid == "" {
		response = comm.Response(400, "validate_failed")
		response.Fields = append(response.Fields, comm.Fields("uuid", "validate_required"))
		service.Logger.Info(ctx, "domain.user.service.destroy.destroy_user_service.Validation")
		return
	}

	user, err := service.Repository.FindByUUID(ctx, uuid)
	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.destroy.destroy_user_service.Repository.FindByUUID", err)
		response = comm.Response(500, "error_delete")
		return
	}

	if user.UUID != "" {
		err = service.Repository.Destroy(ctx, uuid)
		if err != nil {
			service.Logger.Error(ctx, "domain.user.service.destroy.destroy_user_service.Repository.Destroy", err)
			response = comm.Response(500, "error_delete")
			return
		}

		response = comm.Response(200, "success")
		return
	}

	service.Logger.Error(ctx, "domain.user.service.destroy.destroy_user_service.Repository", err)
	response = comm.Response(500, "error_delete")
	return
}
