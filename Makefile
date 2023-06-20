include config.env

swagger-run:
	which swagger || alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOCACHE=/tmp -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'


swagger-check: swagger-run
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: swagger-check
	swagger generate spec -o ./swagger.yaml --scan-models

swag-install:
	which swag || go get github.com/swaggo/swag/cmd/swag@latest

swag: swag-install
	swag init -g server.go -o "./" --outputTypes "yaml"

docker-compose-up:
	docker-compose up -d

docker-compose-stop:
	docker-compose stop

docker-build:
	docker build -t luizalabs-order:latest .

docker-run:
	docker run \
	-p 9000:9000 \
	--net luizalabs-challenge-backend_luizalabs-network \
	--env SERVER_LOG_JSON_FORMAT=false \
	--env DB_DRIVER=postgres \
	--env DB_URL=postgres://userluizalabs:luizaLABS@123@luizalabs-postgres:5432/db_luizalabs?sslmode=disable \
	--env DB_MIGRATION_URL=file://migration \
	--env CACHE_URL=redis://:@luizalabs-redis:6379/0 \
	luizalabs-order

migrate-up:
	migrate -source ${DB_MIGRATION_URL} -database "${DB_URL}" up
	
migrate-down:
	migrate -source ${DB_MIGRATION_URL} -database "${DB_URL}" down
	
go-test: 
	go test -v -cover ./...

go-run:
	go run server.go

.PHONY: swagger swagger-check docker-compose-up docker-compose-stop docker-build docker-run migrate-up migrate-down go-test go-run