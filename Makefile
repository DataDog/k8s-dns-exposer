PROJECT_NAME=k8s-dns-exposer
ARTIFACT=controller
ARTIFACT_PLUGIN=kubectl-${PROJECT_NAME}

# 0.0 shouldn't clobber any released builds
TAG=v0.0.1
DOCKER_REGISTRY =
PREFIX = ${DOCKER_REGISTRY}datadog/${PROJECT_NAME}-${ARTIFACT}
SOURCEDIR = "."

SOURCES := $(shell find $(SOURCEDIR) ! -name "*_test.go" -name '*.go')

BUILDINFOPKG=github.com/DataDog/k8s-dns-exposer/version
GIT_TAG?=$(shell git tag|tail -1)
GIT_COMMIT?=$(shell git rev-parse HEAD)
DATE=$(shell date +%Y-%m-%d/%H:%M:%S )
GOMOD=
LDFLAGS= -ldflags "-w -X ${BUILDINFOPKG}.Tag=${GIT_TAG} -X ${BUILDINFOPKG}.Commit=${GIT_COMMIT} -X ${BUILDINFOPKG}.Version=${TAG} -X ${BUILDINFOPKG}.BuildTime=${DATE} -s"

all: build

vendor:
	go mod vendor

tidy:
	go mod tidy -v

build: ${ARTIFACT}

${ARTIFACT}: ${SOURCES}
	CGO_ENABLED=0 go build ${GOMOD} -i -installsuffix cgo ${LDFLAGS} -o ${ARTIFACT} ./cmd/manager/main.go

build-plugin: ${ARTIFACT_PLUGIN}

${ARTIFACT_PLUGIN}: ${SOURCES}
	CGO_ENABLED=0 go build -i -installsuffix cgo ${LDFLAGS} -o ${ARTIFACT_PLUGIN} ./cmd/${ARTIFACT_PLUGIN}/main.go

container:
	./bin/operator-sdk build $(PREFIX):$(TAG)
    ifeq ($(KINDPUSH), true)
	 kind load docker-image $(PREFIX):$(TAG)
    endif

container-ci:
	docker build -t $(PREFIX):$(TAG) --build-arg  "VERSION=$(TAG)" .

test:
	./go.test.sh

e2e: container	
	./bin/operator-sdk test local ./test/e2e --image $(PREFIX):$(TAG)

push: container
	docker push $(PREFIX):$(TAG)

clean:
	rm -f ${ARTIFACT}

validate:
	./bin/golangci-lint run ./...

generate:
	./bin/operator-sdk generate k8s
	./bin/operator-sdk generate openapi

CRDS = $(wildcard deploy/crds/*_crd.yaml)
local-load: $(CRDS)
		for f in $^; do kubectl apply -f $$f; done
		kubectl apply -f deploy/
		kubectl delete pod -l name=${PROJECT_NAME}

$(filter %.yaml,$(files)): %.yaml: %yaml
	kubectl apply -f $@

install-tools:
	./hack/golangci-lint.sh
	./hack/install-operator-sdk.sh

.PHONY: vendor build push clean test e2e validate local-load install-tools list container container-ci
