NAME=sailfish

build:
	@go build -o $(NAME)

run: build
	@./$(NAME)

dev:
	nodemon --exec go run main.go --signal SIGTERM

.PHONY: dev build run

