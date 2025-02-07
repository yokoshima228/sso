APP_NAME=sso

GO_FILES=main.go

build:
	go build -o sso cmd/sso/main.go

run: build
	./sso