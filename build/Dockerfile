# Общий Dockerfile
FROM golang:1.22.2-alpine3.19 as build

ARG MAIN_PATH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ${MAIN_PATH} .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app

FROM scratch AS final

COPY --from=build /bin/app /app

EXPOSE 8080
ENTRYPOINT ["/app"]
