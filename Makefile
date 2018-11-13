APP?=money.io
PROJECT=github.com/gromnsk/money.io
REGISTRY?=hub.docker.io
CA_DIR?=certs
ASSETS?=pkg/handlers/v1/v1.apib

# Use the 0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE?=0.0.1
GOOS?=linux
GOARCH?=amd64

MONEYIO_LOCAL_HOST?=0.0.0.0
MONEYIO_LOCAL_PORT?=1666
MONEYIO_LOG_LEVEL?=0

# Namespace: dev, prod, release, cte, username ...
NAMESPACE?=dev

# Infrastructure (dev, stable, test ...) and kube-context for helm
INFRASTRUCTURE?=dev

KUBE_CONTEXT?=${INFRASTRUCTURE}
VALUES?=values-${INFRASTRUCTURE}

CONTAINER_IMAGE?=${REGISTRY}/${APP}
CONTAINER_NAME?=${APP}-${NAMESPACE}

REPO_INFO=$(shell git config --get remote.origin.url)

ifndef COMMIT
	COMMIT := git-$(shell git rev-parse --short HEAD)
endif

BUILDTAGS=

.PHONY: all
all: build

.PHONY: vendor
vendor: clean bootstrap
	dep ensure

.PHONY: generate
generate: bootstrap
	-go-bindata -pkg assets -o pkg/assets/assets.go ${ASSETS}

.PHONY: build
build: vendor test certs generate
	@echo "+ $@"
	@CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo \
		-ldflags "-s -w -X ${PROJECT}/pkg/version.RELEASE=${RELEASE} -X ${PROJECT}/pkg/version.COMMIT=${COMMIT} -X ${PROJECT}/pkg/version.REPO=${REPO_INFO}" \
		-o bin/${GOOS}-${GOARCH}/${APP} ${PROJECT}/cmd
	docker build --pull -t $(CONTAINER_IMAGE):$(RELEASE) .

.PHONY: certs
certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name ${CONTAINER_NAME}-certs -d alpine:latest sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait ${CONTAINER_NAME}-certs
	@mkdir -p ${CA_DIR}
	@docker cp ${CONTAINER_NAME}-certs:/etc/ssl/certs/ca-certificates.crt ${CA_DIR}
	@docker rm -f ${CONTAINER_NAME}-certs
endif

.PHONY: push
push: build
	@echo "+ $@"
	@docker push $(CONTAINER_IMAGE):$(RELEASE)

.PHONY: run
run: build
	@echo "+ $@"
	@docker run --name ${CONTAINER_NAME} --restart on-failure:5 -p ${MONEYIO_LOCAL_PORT}:${MONEYIO_LOCAL_PORT} \
		-e "MONEYIO_LOCAL_HOST=${MONEYIO_LOCAL_HOST}" \
		-e "MONEYIO_LOCAL_PORT=${MONEYIO_LOCAL_PORT}" \
		-e "MONEYIO_LOG_LEVEL=${MONEYIO_LOG_LEVEL}" \
		-d $(CONTAINER_IMAGE):$(RELEASE)
	@sleep 1
	@docker logs ${CONTAINER_NAME}

HAS_RUNNED := $(shell docker ps | grep ${CONTAINER_NAME})
HAS_EXITED := $(shell docker ps -a | grep ${CONTAINER_NAME})

.PHONY: logs
logs:
	@echo "+ $@"
	@docker logs ${CONTAINER_NAME}

.PHONY: stop
stop:
ifdef HAS_RUNNED
	@echo "+ $@"
	@docker stop ${CONTAINER_NAME}
endif

.PHONY: start
start: stop
	@echo "+ $@"
	@docker start ${CONTAINER_NAME}

.PHONY: rm
rm:
ifdef HAS_EXITED
	@echo "+ $@"
	@docker rm ${CONTAINER_NAME}
endif

.PHONY: deploy
deploy: push
	helm upgrade ${CONTAINER_NAME} -f charts/${VALUES}.yaml charts --kube-context ${KUBE_CONTEXT} --namespace ${NAMESPACE} --version=${RELEASE} -i --wait

GO_LIST_FILES=$(shell go list ${PROJECT}/... | grep -v vendor)

.PHONY: fmt
fmt:
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"gofmt -s -l {{.Dir}}"{{end}}' ${GO_LIST_FILES} | xargs -L 1 sh -c

.PHONY: lint
lint: bootstrap
	@echo "+ $@"
	@go list -f '{{if len .TestGoFiles}}"golint -min_confidence=0.85 {{.Dir}}/..."{{end}}' ${GO_LIST_FILES} | xargs -L 1 sh -c

.PHONY: vet
vet:
	@echo "+ $@"
	@go vet ${GO_LIST_FILES}

.PHONY: charts
charts: bootstrap
	helm template charts -n ${NAMESPACE} -f charts/${VALUES}.yaml

.PHONY: test
test: vendor fmt lint vet charts
	@echo "+ $@"
	@go test -v -race -cover -tags "$(BUILDTAGS) cgo" ${GO_LIST_FILES}

.PHONY: cover
cover:
	@echo "+ $@"
	@> coverage.txt
	@go list -f '{{if len .TestGoFiles}}"go test -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}} && cat {{.Dir}}/.coverprofile  >> coverage.txt"{{end}}' ${GO_LIST_FILES} | xargs -L 1 sh -c

.PHONY: clean
clean: stop rm
	@rm -f bin/${GOOS}-${GOARCH}/${APP}

HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)
HAS_GO_BINDATA := $(shell command -v go-bindata;)
HAS_HELM_TEMPLATE := $(shell helm plugin list | grep "template";)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_LINT
	go get -u github.com/golang/lint/golint
endif
ifndef HAS_GO_BINDATA
	go get -u github.com/jteeuwen/go-bindata/...
endif
ifndef HAS_HELM_TEMPLATE
	helm plugin install https://github.com/technosophos/helm-template
endif

.PHONY: git-init
git-init: 
	git init
	git config --global user.name "Gromnsk"
	git config --global user.email "gromnsk@yandex.ru"
