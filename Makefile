up-traefik:
	docker compose -f ./traefik.docker-compose.yml up -d --build

down-traefik:
	docker compose -f ./traefik.docker-compose.yml down


up-gateway:
	docker compose -f ./api-gateway.docker-compose.yml up -d --build

down-gateway:
	docker compose -f ./api-gateway.docker-compose.yml down

# Protobuf generate
proto-gen:
	protoc --go_out=./internal/shared --go-grpc_out=./internal/shared --proto_path=./scripts/proto ./scripts/proto/*.proto