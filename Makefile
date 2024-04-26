LOCAL_MIGRATION_DIR=./migrations
POSTGRES = ${POSTGRES_URL}

# Cекция настройки gRPC
install-grpc-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

# Устанавливаем proto описания google/googleapis
vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis &&\
 	cd vendor-proto/googleapis &&\
	git sparse-checkout set --no-cone google/api &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

# Устанавливаем proto описания google/protobuf
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

# Устанавливаем proto описания validate
vendor-proto/validate:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/validate-repo &&\
	cd vendor-proto/validate-repo &&\
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p  vendor-proto
	mv vendor-proto/validate-repo/validate vendor-proto
	rm -rf vendor-proto/validate-repo

generate_proto:
	mkdir -p generated/task
	protoc -I proto/task -I vendor-proto \
	--go_out generated/task --go_opt paths=source_relative \
	--go-grpc_out generated/task --go-grpc_opt paths=source_relative \
	--grpc-gateway_out generated/task --grpc-gateway_opt paths=source_relative \
	--validate_out="lang=go,paths=source_relative:generated/task" \
	proto/task/service.proto


generate: install-grpc-deps vendor-proto/google/api vendor-proto/google/protobuf vendor-proto/validate generate_proto

# Команды для работы с проектом
lint:
	golangci-lint run

install-linters:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

run-worker:
	go run cmd/worker/main.go

test:
	go test -v ./...

# работа с базой данных
run-database:
	docker-compose -f ./deployments/docker-compose.yaml up --build

# Секция работы с миграциями
migrate:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres $(POSTGRES) status && \
	goose -dir ${LOCAL_MIGRATION_DIR} postgres $(POSTGRES) up
