APP_NAME := cloudevents-podgateway-client-provisioner
REGISTRY ?= localhost:5001
REPOSITORY ?= $(REGISTRY)/clastix
VERSION ?= latest
IMAGE := $(REPOSITORY)/$(APP_NAME):$(VERSION)

.DEFAULT_GOAL := oci/build

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP_NAME) .

.PHONY: oci
oci: oci/build oci/push

.PHONY: oci/build
oci/build:
	docker build . -t $(APP_NAME)

.PHONY: oci/push
oci/push: kind-registry
	docker tag $(APP_NAME) $(IMAGE)
	docker push $(IMAGE)

.PHONY: kind-registry
kind-registry:
	(curl -sL https://kind.sigs.k8s.io/examples/kind-with-registry.sh | bash -) || true

.PHONY: publish/local
publish/local: kind-registry oci

.PHONY: cleanup
cleanup:
	kind delete cluster
	docker stop kind-registry && docker rm kind-registry
