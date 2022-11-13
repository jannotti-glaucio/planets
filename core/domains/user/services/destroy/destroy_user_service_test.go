package destroy

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jannotti-glaucio/planets/core/domains/user/entities"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/mocks"
)

func TestService(t *testing.T) {
	comm := communication.New()

	useCases := map[string]struct {
		expectedResponse communication.Response
		inputData        string
		prepare          func(repostitoryMock *mocks.MockIUserRepository, loggerMock *mocks.MockILoggerProvider)
	}{
		"success": {
			inputData: "userID",
			expectedResponse: communication.Response{
				Status:  200,
				Code:    comm.Mapping["success"].Code,
				Message: comm.Mapping["success"].Message,
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.User{UUID: "userID"}, nil)
				repostitoryMock.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		"error: on validation": {
			expectedResponse: communication.Response{
				Status:  400,
				Code:    comm.Mapping["validate_failed"].Code,
				Message: comm.Mapping["validate_failed"].Message,
			},
			prepare: func(repostitoryMock *mocks.MockIUserRepository, loggerMock *mocks.MockILoggerProvider) {
				loggerMock.EXPECT().Info(gomock.Any(), gomock.Any())
			},
		},
		"error: on repository FindByUUID": {
			expectedResponse: communication.Response{
				Status:  500,
				Code:    comm.Mapping["error_delete"].Code,
				Message: comm.Mapping["error_delete"].Message,
			},
			inputData: "userID",
			prepare: func(repostitoryMock *mocks.MockIUserRepository, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.User{}, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
		"error: on repository Destroy": {
			expectedResponse: communication.Response{
				Status:  500,
				Code:    comm.Mapping["error_delete"].Code,
				Message: comm.Mapping["error_delete"].Message,
			},
			inputData: "userID",
			prepare: func(repostitoryMock *mocks.MockIUserRepository, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.User{UUID: "userID"}, nil)
				repostitoryMock.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},

		"error: on repository": {
			expectedResponse: communication.Response{
				Status:  500,
				Code:    comm.Mapping["error_delete"].Code,
				Message: comm.Mapping["error_delete"].Message,
			},
			inputData: "userID",
			prepare: func(repostitoryMock *mocks.MockIUserRepository, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.User{}, nil)
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any())
			},
		},
	}

	for name, useCase := range useCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			defer ctrl.Finish()

			repository := mocks.NewMockIUserRepository(ctrl)
			logger := mocks.NewMockILoggerProvider(ctrl)
			useCase.prepare(repository, logger)

			service := Service{
				Repository: repository,
				Logger:     logger,
			}
			response := service.Execute(ctx, useCase.inputData)

			if response.Status != useCase.expectedResponse.Status {
				t.Errorf("Expected %d, but got %d", useCase.expectedResponse.Status, response.Status)
			}

			if response.Message != useCase.expectedResponse.Message {
				t.Errorf("Expected %s, but got %s", useCase.expectedResponse.Message, response.Message)
			}
		})
	}
}
