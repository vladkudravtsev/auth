#build stage
FROM golang:1.21.5 AS builder
ARG SERVICE_NAME
WORKDIR /src

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o microservice cmd/$SERVICE_NAME/main.go

#final stage
FROM alpine:3.16.2
WORKDIR /bin/
COPY --from=builder /src/microservice /bin/microservice
ENV GIN_MODE=release
CMD /bin/microservice
