proto:
	protoc -I api/proto api/proto/auth/auth.proto \
		--go_out=./api/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=./api/gen/go \
		--go-grpc_opt=paths=source_relative \
		--validate_out="lang=go,paths=source_relative:./api/gen/go" \
		-I ./api/proto/validate/validate.proto

.PHONY: cover
cover:
	go test -short -count=1 -race -coverpofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
