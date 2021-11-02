NAME := kramer-microservice
OWNER := byuoitav
PKG := github.com/${OWNER}/${NAME}
DOCKER_URL := ghcr.io
DOCKER_PKG := ${DOCKER_URL}/${OWNER}/${NAME}
BRANCH:= $(shell git rev-parse --abbrev-ref HEAD)

# go
VENDOR=gvt fetch -branch $(BRANCH)

# docker
UNAME=$(shell echo $(DOCKER_USERNAME))
EMAIL=$(shell echo $(DOCKER_EMAIL))
PASS=$(shell echo $(DOCKER_PASSWORD))

build: deps
	@mkdir -p dist
	@env GOOS=linux CGO_ENABLED=0 go build -o ./dist/$(NAME)-linux-amd64 -v
	@env GOOS=linux GOARCH=arm go build -o ./dist/$(NAME)-arm -v
	
test: 
	@go test -v -race $(go list ./... | grep -v /vendor/) 

clean: 
	@go clean
	@rm -rf dist/

deps: 
	@go get -d -v

docker: clean build
	@echo Building arm container
	@docker build --build-arg NAME=$(NAME)-arm -f dockerfile -t $(DOCKER_PKG)/$(NAME)-arm:$(BRANCH) dist

deploy: docker
	@echo logging in to dockerhub...
	@docker login $(DOCKER_URL) -u $(UNAME) -p $(PASS)

	@docker push ${DOCKER_PKG}/$(NAME)-arm:$(BRANCH)
