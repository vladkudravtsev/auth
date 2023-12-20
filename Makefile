proto:
	protoc -I api/proto api/proto/auth/auth.proto \
		--go_out=./api/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=./api/gen/go \
		--go-grpc_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:./api/gen/go" \
		-I ./api/proto/validate/validate.proto
.PHONY: test-ci
test-ci:
	@docker compose -f test/docker-compose.yml down -v
	@docker compose -f test/docker-compose.yml up --build --abort-on-container-exit --remove-orphans --force-recreate
	@docker compose -f test/docker-compose.yml down -v
.PHONY: test
test:
	go test ./...
run:
	go run cmd/main.go -config auth.env
migrate:
	go run cmd/migrator/main.go -config auth.env
