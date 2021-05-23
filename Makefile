GOCMD=go
BINARY_NAME=gateway
VERSION?=0.0.0
SERVICE_PORT?=8080
DOCKER_REGISTRY?=scarletfairy/#if set it should finished by /
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true
BIN_FOLDER?=bin/
MAIN_PATH?=cmd/gateway/main.go

build:
	@echo 'Building ${BINARY_NAME}'
	@mkdir -p bin
	@$(GOCMD) build -o $(BIN_FOLDER)$(BINARY_NAME) $(MAIN_PATH)

run:
	@$(GOCMD) run $(MAIN_PATH)

docker-run: docker-build
	docker run --network host $(BINARY_NAME)

docker-build:
	docker build --rm --tag $(BINARY_NAME) .

docker-release:
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)

proto-build:
	buf generate

