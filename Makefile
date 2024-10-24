build-all:
	docker build --tag msa-messenger-db --file build/db/Dockerfile .
	docker build --tag msa-messenger-auth --file build/auth/Dockerfile .
	docker build --tag msa-messenger-user --file build/user/Dockerfile .
	docker build --tag msa-messenger-messaging --file build/messaging/Dockerfile .
delete-all:
	docker image rm --force msa-messenger-auth msa-messenger-user msa-messenger-messaging
run-all:
	docker-compose -f build/docker-compose.yml up --build
minikube-docker:
	eval $(minikube docker-env)
	minikube ssh -- docker rmi -f msa-messenger-db:latest
	minikube ssh -- docker rmi -f msa-messenger-auth:latest
	minikube ssh -- docker rmi -f msa-messenger-user:latest
	minikube ssh -- docker rmi -f msa-messenger-messaging:latest
	minikube image load msa-messenger-db:latest
	minikube image load msa-messenger-auth:latest
	minikube image load msa-messenger-user:latest
	minikube image load msa-messenger-messaging:latest
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

# Сборка приложения
build:
	CGO_ENABLED=0 go build -v -o $(LOCAL_BIN) ./cmd/user

# Объявляем, что текущие команды не являются файлами и
# интсрументируем Makefile не искать изменения в файловой системе
.PHONY: \
	.bin-deps \
	.protoc-generate \
	.tidy \
	generate \
	build