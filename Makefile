all: buf go

buf:
	buf generate --exclude-path google,openapiv3

go:
	go run main.go

PHONY: buf go