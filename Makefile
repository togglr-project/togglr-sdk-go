OGEN_IMAGE=ghcr.io/ogen-go/ogen:latest
SPEC=specs/sdk.yml
OUT=internal/generated/client

.PHONY: generate build test lint tidy clean
generate:
	@mkdir -p $(OUT)
	@docker run --rm \
      --volume ".:/workspace" \
      ghcr.io/ogen-go/ogen:latest --target workspace/internal/generated/client --config workspace/ogen-config.yml --clean workspace/specs/sdk.yml

build:
	go build ./...

test:
	go test ./... -race

lint:
	gofmt -l .
	go vet ./...
	staticcheck ./...

tidy:
	go mod tidy

clean:
	rm -rf $(OUT)
