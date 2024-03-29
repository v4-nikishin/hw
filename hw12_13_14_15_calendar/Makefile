BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_IMG_CALENDAR="calendar:develop"
DOCKER_IMG_SCHEDULER="scheduler:develop"
DOCKER_IMG_SENDER="sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml
	$(BIN_SENDER) -config ./configs/sender_config.yaml

build-img:
	# docker build \
	# 	--build-arg=LDFLAGS="$(LDFLAGS)" \
	# 	-t $(DOCKER_IMG_CALENDAR) \
	# 	-f build/Dockerfile.calendar .
	# docker build \
	# 	--build-arg=LDFLAGS="$(LDFLAGS)" \
	# 	-t $(DOCKER_IMG_SCHEDULER) \
	# 	-f build/Dockerfile.scheduler .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/Dockerfile.sender .

run-img: build-img
	docker run $(DOCKER_IMG_CALENDAR)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb
	go generate internal/server/grpc/server.go

up:
	docker-compose up -d --build

down:
	docker-compose down

integration-test:
	docker-compose -f docker-compose-testing.yaml up -d --build
	sleep 30 || echo "Waiting for prepear environment..."
	go test ./tests/... -count=1 -v || echo "Integration tests"
	docker-compose -f docker-compose-testing.yaml down

.PHONY: build run build-img run-img version test lint
