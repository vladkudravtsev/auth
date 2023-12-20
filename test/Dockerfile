FROM golang:1.21.5-alpine3.18
WORKDIR /root

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY api ./api
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

CMD CGO_ENABLED=0 go run cmd/migrator/main.go && go test ./...
