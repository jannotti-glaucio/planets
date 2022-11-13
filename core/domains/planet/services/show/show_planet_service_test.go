package show

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jannotti-glaucio/planets/core/domains/planet/entities"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/mocks"
)

func TestService(t *testing.T) {
	comm := communication.New()

	expectedData := entities.Planet{
		UUID:    "uuid",
		Name:    "name",
		Terrain: "terrain",
		Climate: "climate",
	}

	useCases := map[string]struct {
		expectedResponse communication.Response
		expectedData     entities.Planet
		inputData        string
		prepare          func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider)
	}{
		"success": {
			inputData:    "planetID",
			expectedData: expectedData,
			expectedResponse: communication.Response{
				Status:  200,
				Code:    comm.Mapping["success"].Code,
				Message: comm.Mapping["success"].Message,
			},
			prepare: func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(expectedData, nil)
			},
		},
		"error: on repository": {
			expectedResponse: communication.Response{
				Status:  404,
				Code:    comm.Mapping["error_list"].Code,
				Message: comm.Mapping["error_list"].Message,
			},
			prepare: func(repostitoryMock *mocks.MockIPlanetRepository, loggerMock *mocks.MockILoggerProvider) {
				repostitoryMock.EXPECT().FindByUUID(gomock.Any(), gomock.Any()).Return(entities.Planet{}, errors.New("error"))
				loggerMock.EXPECT().Error(gomock.Any(), gomock.Any())
			},
		},
	}

	for name, useCase := range useCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()
			defer ctrl.Finish()

			repository := mocks.NewMockIPlanetRepository(ctrl)
			logger := mocks.NewMockILoggerProvider(ctrl)
			useCase.prepare(repository, logger)

			service := Service{
				Repository: repository,
				Logger:     logger,
			}
			data, response := service.Execute(ctx, useCase.inputData)

			if response.Status != useCase.expectedResponse.Status {
				t.Errorf("Expected %d, but got %d", useCase.expectedResponse.Status, response.Status)
			}

			if response.Message != useCase.expectedResponse.Message {
				t.Errorf("Expected %s, but got %s", useCase.expectedResponse.Message, response.Message)
			}

			if data != useCase.expectedData {
				t.Errorf("Expected %q, but got %q", useCase.expectedData, data)
			}
		})
	}
}
