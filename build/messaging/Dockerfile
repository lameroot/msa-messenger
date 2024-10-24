FROM golang:1.22.2-alpine3.19 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/messaging ./cmd/messaging
COPY ./internal/messaging ./internal/messaging
COPY ./pkg ./pkg
COPY ./docs ./docs
COPY ./docs/arch /dev/null

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/messaging

FROM scratch AS final

COPY --from=build /bin/app /app

EXPOSE 8080
ENTRYPOINT ["/app"]
