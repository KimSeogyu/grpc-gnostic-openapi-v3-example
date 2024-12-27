all: buf go

buf:
	@buf dep update
	@buf generate
	@npx @redocly/cli build-docs docs/openapi.yaml -o ./docs/spec.html

go:
	@go run main.go

PHONY: buf go