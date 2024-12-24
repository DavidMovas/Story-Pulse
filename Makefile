# Docker compose commands
## Traefik
up-traefik:
	docker compose -f ./traefik.docker-compose.yml up -d --build

down-traefik:
	docker compose -f ./traefik.docker-compose.yml down

## Gateway
up-gateway:
	docker compose -f ./api-gateway.docker-compose.yml up -d --build

down-gateway:
	docker compose -f ./api-gateway.docker-compose.yml down

## Clear network
net-create:
	docker network create gateway
	docker network create users_db_net

net-clear:
	docker network rm gateway
	docker network rm users_db_net

## Gateway splitted
full-up:
	make net-clear
	make net-create
	docker-compose -f ./deployments/compose/gateway.docker-compose.yml --env-file=./.env up -d --build
	docker-compose -f ./deployments/compose/auth-service.docker-compose.yml --env-file=./.env up -d --build
	docker-compose -f ./deployments/compose/users-service.docker-compose.yml --env-file=./.env up -d --build

full-down:
	docker-compose -f ./deployments/compose/users-service.docker-compose.yml --env-file=./.env down
	docker-compose -f ./deployments/compose/auth-service.docker-compose.yml --env-file=./.env down
	docker-compose -f ./deployments/compose/gateway.docker-compose.yml --env-file=./.env down
	make net-clear

gateway-up:
	docker-compose -f ./deployments/compose/gateway.docker-compose.yml --env-file=./.env up -d --build

gateway-down:
	docker-compose -f ./deployments/compose/gateway.docker-compose.yml --env-file=./.env down

auth-up:
	docker-compose -f ./deployments/compose/auth-service.docker-compose.yml --env-file=./.env up -d --build

auth-down:
	docker-compose -f ./deployments/compose/auth-service.docker-compose.yml --env-file=./.env downd

users-up:
	docker-compose -f ./deployments/compose/users-service.docker-compose.yml --env-file=./.env up -d --build

users-down:
	docker-compose -f ./deployments/compose/users-service.docker-compose.yml --env-file=./.env down


# Protobuf generate commands
GRPC_VERSION := "v1"

proto-gen:
	protoc \
      -I ./scripts/proto \
      -I /home/nofre/go/pkg/mod/github.com/googleapis/googleapis@v0.0.0-20241220203547-09b3c838b775 \
      --go_out=./internal/shared \
      --go-grpc_out=./internal/shared/grpc/"${GRPC_VERSION}" \
      --grpc-gateway_out ./internal/shared/grpc/"${GRPC_VERSION}" \
      --grpc-gateway_opt logtostderr=true \
      --go-grpc_opt paths=source_relative \
      --grpc-gateway_opt paths=source_relative \
      --experimental_allow_proto3_optional=true \
      ./scripts/proto/*.proto