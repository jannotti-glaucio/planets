package authenticate

import (
	"context"

	"github.com/jannotti-glaucio/planets/core/domains/user/repositories"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/hash"
	"github.com/jannotti-glaucio/planets/core/tools/providers/jwt"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/validations"
)

//Dto object receiver
type Dto struct {
	Email    string `json:"email" validate:"isRequired|isEmail"`
	Password string `json:"password" validate:"isRequired|isPassword"`
}

//Service ...
type Service struct {
	Repository repositories.IUserRepository
	Hash       hash.IHashProvider
	Jwt        jwt.IJwtProvider
	Logger     logger.ILoggerProvider
}

//Execute Serviço responsável pela inserção de registros
func (service *Service) Execute(ctx context.Context, dto Dto) (token string, response communication.Response) {
	response.Fields = validations.ValidateStruct(&dto, "")
	comm := communication.New()

	if len(response.Fields) > 0 {
		service.Logger.Info(ctx, "domain.user.service.authenticate.authenticate_user_service.ValidationError")
		resp := comm.Response(400, "validate_failed")
		resp.Fields = response.Fields
		response = resp
		return
	}

	user, err := service.Repository.FindByEmail(ctx, dto.Email)
	if err != nil {
		service.Logger.Error(ctx, "domain.user.service.authenticate.authenticate_user_service.Repository.FindByEmail", err)
		response = comm.Response(400, "authenticate_failed")
		return
	}

	if user.UUID != "" {
		//Apply security hash in password
		checkPassword := service.Hash.Compare(user.Password, dto.Password)

		if checkPassword {
			tokenDetails, err := service.Jwt.CreateToken(jwt.TokenPayload{UserUUID: user.UUID})
			if err != nil {
				service.Logger.Error(ctx, "domain.user.service.authenticate.authenticate_user_service.Jwt.CreateToken", err)
				response = comm.Response(400, "authenticate_failed")
				return
			}

			token = tokenDetails.AccessToken
			response = comm.Response(200, "authenticate_success")
			return
		}

	}

	response = comm.Response(400, "authenticate_failed")
	return
}
