build:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/app main.go

run: build
	@./bin/app

docker: build
	@docker compose up --build -d

.PHONY: build