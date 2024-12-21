up-traefik:
	docker compose -f ./traefik.docker-compose.yml up -d --build

down-traefik:
	docker compose -f ./traefik.docker-compose.yml down


up-gateway:
	docker compose -f ./api-gateway.docker-compose.yml up -d --build

down-gateway:
	docker compose -f ./api-gateway.docker-compose.yml down

# Protobuf generate
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