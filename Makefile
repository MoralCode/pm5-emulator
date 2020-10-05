test:
	@go test -v ./...

build:
	@GOOS=linux GOARCH=arm GOARM=5 go build -o pm5-emulator cmd/pm5-emulator/main.go
	@echo "build complete use 'sudo ./pm5-emulator' to run"

all: test build
