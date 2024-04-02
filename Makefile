all:
	buf generate --exclude-path google,openapiv3
	go run main.go