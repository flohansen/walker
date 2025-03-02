DOCKER ?= docker
KCTL ?= kubectl

TAG ?= latest

.PHONY: generate
generate:
	go generate ./...

.PHONY: setup
setup: generate
	cp .githooks/pre-push .git/hooks/pre-push

.PHONY: test
test: generate
	go test ./... -cover -race

.PHONY: docker-build
docker-build:
	$(DOCKER) build -t api:$(TAG) -f build/api.Dockerfile .

.PHONY: deploy-local
deploy-local:
	minikube image rm api:$(TAG) || true
	minikube image load api:$(TAG)
	$(KCTL) apply -k deployment/k8s/

.PHONY: undeploy
undeploy:
	$(KCTL) delete -k deployment/k8s/
