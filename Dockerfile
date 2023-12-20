#build stage
FROM golang:1.21.5 AS builder
WORKDIR /src

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

RUN CGO_ENABLED=0 go build -o microservice cmd/main.go
RUN CGO_ENABLED=0 go build -o migrator cmd/migrator/main.go

#final stage
FROM alpine:3.16.2
WORKDIR /bin/
COPY --from=builder /src/migrator /bin/migrator
COPY --from=builder /src/microservice /bin/microservice
COPY --from=builder /src/migrations /bin/migrations
ENV GIN_MODE=release
CMD /bin/migrator && /bin/microservice
