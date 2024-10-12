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
build-db:
	docker build -t my-postgres -f build/db/Dockerfile build/db
run-db:
	docker run -d --name my-postgres-container -p 5432:5432 my-postgres
stop-db:
	docker stop my-postgres-container
start-db:
	docker start my-postgres-container
remove-db:
	docker rm my-postgres-container	