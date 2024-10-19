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