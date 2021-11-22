GOLANG_VERSION := 1.15.6
ALPINE_VERSION := 3.13

DOCKER_REPO := michalsw

PORT ?= 8080
# APORT ?= 8080
PPORT ?= 5050

.DEFAULT_GOAL := help
.PHONY: go-build docker-build docker-run-webapp docker-run-proxy

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

go-build: ## Build binary
	CGO_ENABLED=0 \
	go build \
	-v \
	-ldflags "-s -w" \
	-o $(APPNAME) $(APPNAME).go

docker-build: ## Build docker image
	docker build \
	--file $(APPNAME)_Dockerfile \
	--pull \
	--build-arg APPNAME="$(APPNAME)" \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--tag="$(DOCKER_REPO)/$(APPNAME):latest" .

docker-run-webapp: ## Run webapp in docker
	docker run -d --rm \
	--name $(APPNAME) \
	-e PORT=$(PORT) \
	-p $(PORT):$(PORT) \
	$(DOCKER_REPO)/$(APPNAME):latest

docker-run-proxy: ## Run proxy in docker
	docker run -d --rm \
	--name $(APPNAME) \
	-e APORT=$(PORT) \
	-e PPORT=$(PPORT) \
	-e AIP=$(AIP) \
	-p $(PPORT):$(PPORT) \
	$(DOCKER_REPO)/$(APPNAME):latest
