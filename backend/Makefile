.PHONY: proto sqlc build test run

proto:
	mkdir -p proto/gen
	protoc -I. -I./third_party/googleapis \
		--go_out=. --go_opt=module=github.com/mehmetymw/search-aggregation-service/backend \
		--go-grpc_out=. --go-grpc_opt=module=github.com/mehmetymw/search-aggregation-service/backend \
		--grpc-gateway_out=. --grpc-gateway_opt=module=github.com/mehmetymw/search-aggregation-service/backend \
		--grpc-gateway_opt=generate_unbound_methods=true \
		proto/content.proto

sqlc:
	cd db && sqlc generate

build:
	go build -o bin/api ./cmd/api

test:
	go test ./... -v

run:
	go run ./cmd/api
