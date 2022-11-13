CURRENTNAME = $(shell pwd | sed 's!.*/!!')
NAME = $(shell echo $(CURRENTNAME) | sed 's/\(.\)\([A-Z]\)/\1-\2/g' | tr '[:upper:]' '[:lower:]')
TAG = $(shell /bin/date "+%Y%m%d%H%M%S")
VERSION = $(shell git branch --show-current)

# comandos para execução a api echo
run-api:
	VERSION=$(VERSION) go run ./adapters/rest/main.go

run-seeds:
	VERSION=$(VERSION) go run ./seeds/main.go

# comandos para gerar mock
mock:
	rm -rf ./mocks

	mockgen -source=./core/tools/providers/cache/cache_provider.go -destination=./mocks/cache_provider_mock.go -package=mocks	
	mockgen -source=./core/tools/providers/hash/hash_provider.go -destination=./mocks/hash_provider_mock.go -package=mocks	
	mockgen -source=./core/tools/providers/logger/logger_provider.go -destination=./mocks/logger_provider_mock.go -package=mocks	
	mockgen -source=./core/tools/providers/http_client/http_client_provider.go -destination=./mocks/http_client_provider_mock.go -package=mocks	
	mockgen -source=./core/tools/providers/jwt/jwt_provider.go -destination=./mocks/jwt_provider_mock.go -package=mocks	
	mockgen -source=./core/domains/user/repositories/user_repository.go -destination=./mocks/user_repository_mock.go -package=mocks	
	mockgen -source=./core/domains/planet/repositories/planet_repository.go -destination=./mocks/planet_repository_mock.go -package=mocks	
	
# comandos para gerar testes
test:	
	go test -v -p 1 -cover ./... -coverprofile=coverage.out
	@go tool cover -func coverage.out | awk 'END{print sprintf("coverage: %s", $$3)}'

test-cover: test
	go tool cover -html=coverage.out