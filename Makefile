BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG := "calendar:develop"
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_DSN := user=user password=password dbname=calendar host=localhost port=5432 sslmode=disable

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) ./cmd/calendar_sender

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--build-arg=DB_DSN="$(DB_DSN)" \
		-t $(DOCKER_IMG) \
		-f deployments/build/Dockerfile .

run-db:
	docker-compose -f deployments/docker-compose.yaml up -d db

wait-db:
	@echo "Waiting for database to be ready..."
	@until docker exec postgres_db pg_isready -U user; do \
		echo "Waiting for database..."; \
		sleep 2; \
	done

run-rabbitmq:
	docker-compose -f deployments/docker-compose.yaml up -d rabbitmq

wait-rabbitmq:
	@echo "Waiting for RabbitMQ to be ready..."
	@until [ "`docker inspect -f {{.State.Health.Status}} rabbitmq`" == "healthy" ]; do \
       sleep 2; \
    done
	@echo "RabbitMQ is ready!"

run-scheduler:
	docker run --network calendar_network --name calendar_scheduler_service --rm calendar:develop /opt/calendar/calendar-scheduler-app -config /app/configs/scheduler_config.yaml

run-sender:
	docker run --network calendar_network --name calendar_sender_service --rm calendar:develop /opt/calendar/calendar-sender-app -config /app/configs/sender_config.yaml

run: run-db wait-db run-rabbitmq wait-rabbitmq stop-old-container build-img
	@trap 'docker stop calendar_service' SIGINT; \
	if [ "$(docker ps -q -f name=calendar_service)" ]; then \
		docker stop calendar_service; \
		docker rm calendar_service; \
	fi; \
	docker run --network calendar_network -p 8080:8080 -p 50051:50051 --name calendar_service --rm $(DOCKER_IMG)

stop-old-container:
	@docker ps -q -f name=calendar_service | xargs -r docker stop
	@docker ps -aq -f name=calendar_service | xargs -r docker rm

run-img: run-db build-img
	docker run -p 8080:8080 -p 50051:50051 --name calendar_service --rm $(DOCKER_IMG)

generate:
	protoc -I=api --go_out=api/pb --go-grpc_out=api/pb api/*.proto

up: build
	@docker-compose -f deployments/docker-compose.yaml up --build

down:
	@docker-compose -f deployments/docker-compose.yaml down --remove-orphans

version: build
	$(BIN) version

test:
	go test -race -cover ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest

lint: install-lint-deps
	golangci-lint run ./...
	go fmt ./...
	go vet ./...

migrate:
	docker exec -it calendar_service goose -dir /migrations postgres "postgres://user:password@db:5432/calendar?sslmode=disable" up

integration-tests:
	docker-compose -f deployments/docker-compose.test.yaml up --build -d
	@sleep 10
	@EXIT_CODE=0; \
	docker-compose -f deployments/docker-compose.test.yaml run --rm integration_tests || EXIT_CODE=$$?; \
	docker-compose -f deployments/docker-compose.test.yaml down --volumes; \
	echo "Integration tests exited with code: $$EXIT_CODE"; \
	exit $$EXIT_CODE


.PHONY: build build-img run run-db wait-db run-rabbitmq wait-rabbitmq run-scheduler run-sender run-img generate up down version test install-lint-deps lint migrate integration-tests

