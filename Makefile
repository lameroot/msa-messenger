build-all:
	docker build --tag msa-messanger-user --build-arg MAIN_PATH=cmd/user/*.go --file build/Dockerfile .
	docker build --tag msa-messanger-auth --build-arg MAIN_PATH=cmd/auth/*.go --file build/Dockerfile .
	docker build --tag msa-messanger-messaging --build-arg MAIN_PATH=cmd/messaging/*.go --file build/Dockerfile .
delete-all:
	docker image rm --force msa-messanger-auth msa-messanger-user msa-messanger-messaging
run-all:
	docker-compose -f build/docker-compose.yml up --build
minikube-docker:
	eval $(minikube docker-env)
	minikube image load msa-messanger-user:latest
	minikube image load msa-messanger-auth:latest
	minikube image load msa-messanger-messaging:latest
	minikube image ls --format table

swag-gen:
	swag init -g cmd/auth/main.go
	swag init -g cmd/user/main.go
	swag init -g cmd/messaging/main.go

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin

# Добавляем bin в текущей директории в PATH при запуске protoc
PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

# Путь до protobuf файлов
PROTO_PATH := $(CURDIR)/api/auth

# Путь до сгенеренных .pb.go файлов
PKG_PROTO_PATH := $(CURDIR)/pkg/

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# генерация .go файлов с помощью protoc
.protoc-generate:
	mkdir -p $(PKG_PROTO_PATH)
	$(PROTOC) --proto_path=$(CURDIR) \
	--go_out=$(PKG_PROTO_PATH) --go_opt paths=source_relative \
	--go-grpc_out=$(PKG_PROTO_PATH) --go-grpc_opt paths=source_relative \
	$(PROTO_PATH)/*.proto

# go mod tidy
.tidy:
	GOBIN=$(LOCAL_BIN) go mod tidy

# Генерация кода из protobuf
generate: .bin-deps .protoc-generate .tidy		

# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.protoc-generate \
	.tidy \
	generate \
	build