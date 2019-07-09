.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: env-up run

##### BUILD
build:
	@echo "Build ..."
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && docker-compose up --force-recreate -d
	@echo "Sleep 5 seconds in order to let the environment setup correctly"
	@sleep 5
	@echo "Environment up"

down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@go run main.go
