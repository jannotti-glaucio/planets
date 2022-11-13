package update

import (
	"context"
	"encoding/json"

	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/domains/user/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/hash"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
	"github.com/jannotti-glaucio/planets/core/tools/validations"
)

//Dto object receiver
type Dto struct {
	UUID     string `json:"uuid" validate:"isRequired"`
	Name     string `json:"name" validate:"isRequired"`
	Email    string `json:"email" validate:"isRequired|isEmail"`
	Password string `json:"password" validate:"isPassword"`
}

//Service ...
type Service struct {
	Repository repositories.IUserRepository
	Hash       hash.IHashProvider
	Logger     logger.ILoggerProvider
}

//Execute responsável por atualizar registros
func (service *Service) Execute(ctx context.Context, dto Dto) (updated entities.User, response communication.Response) {
	identifierTracer := "update.user.service"
	span := tracer.New(identifierTracer).StartSpanWidthContext(ctx, identifierTracer, tracer.Options{Key: identifierTracer + ".dto", Value: dto})
	defer span.Finish()

	response.Fields = validations.ValidateStruct(&dto, "")
	comm := communication.New()

	//Check e-mail in use
	userByEmail, err := service.Repository.FindByEmail(ctx, dto.Email)
	if err != nil {
		service.Logger.Info(ctx, "domain.user.service.update.update_user_service.Repository.FindByEmail", err)
	}

	if userByEmail.UUID != "" && userByEmail.UUID != dto.UUID {
		response.Fields = append(response.Fields, comm.Fields("email", "already_exists"))
	}

	//Check exists user with this identifier
	userFinderUUID, err := service.Repository.FindByUUID(ctx, dto.UUID)
	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.update.update_user_service.Repository.FindByUUID", err)
		response = comm.Response(500, "error_update")
		return
	}

	if userFinderUUID.UUID == "" {
		response.Fields = append(response.Fields, comm.Fields("uuid", "validate_invalid"))
	}

	if len(response.Fields) > 0 {
		service.Logger.Info(ctx, "domain.user.service.update.update_user_service.ValidationError")
		resp := comm.Response(400, "validate_failed")
		resp.Fields = response.Fields
		response = resp
		return
	}

	//Apply security hash in password
	if dto.Password != "" {
		dto.Password = service.Hash.Create(dto.Password)
	} else {
		dto.Password = userFinderUUID.Password
	}

	//Mergin entity and DTO
	toMerge, _ := json.Marshal(dto)
	json.Unmarshal(toMerge, &userFinderUUID)

	updated, err = service.Repository.Save(ctx, userFinderUUID)

	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.update.update_user_service.Repository.Save", err)
		response = comm.Response(500, "error_update")
		return
	}

	response = comm.Response(200, "success")
	return
}
