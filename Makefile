# Define variables
REGISTRY             := raihankhanraka
PUBLISHER_REPO       := rabbitmq-publisher
CONSUMER_REPO        := rabbitmq-consumer
PUBLISHER_IMAGE_NAME := $(REGISTRY)/$(PUBLISHER_REPO)
CONSUMER_IMAGE_NAME  := $(REGISTRY)/$(CONSUMER_REPO)
PUBLISHER_IMAGE_TAG  := v1.0.0
CONSUMER_IMAGE_TAG   := v1.0.0
PUBLISHER_BUILD_DIR  := rabbit-publisher-go
CONSUMER_BUILD_DIR   := rabbit-consumer-go

# Default target
all: build push deploy

build-publisher:
	@cd $(PUBLISHER_BUILD_DIR) && docker build -t $(PUBLISHER_IMAGE_NAME):$(PUBLISHER_IMAGE_TAG) .

build-consumer:
	@cd $(CONSUMER_BUILD_DIR) && docker build -t $(CONSUMER_IMAGE_NAME):$(CONSUMER_IMAGE_TAG) .

push-publisher:
	@docker push $(PUBLISHER_IMAGE_NAME):$(PUBLISHER_IMAGE_TAG)

push-consumer:
	@docker push $(CONSUMER_IMAGE_NAME):$(CONSUMER_IMAGE_TAG)

deploy-publisher:
	@cd $(PUBLISHER_BUILD_DIR) && kubectl apply -f deployment.yaml

deploy-consumer:
	@cd $(CONSUMER_BUILD_DIR) && kubectl apply -f deployment.yaml

publisher: build-publisher push-publisher deploy-publisher

consumer: build-consumer push-consumer deploy-consumer

# Build the Docker image
build: build-publisher build-consumer

push: push-publisher push-consumer

deploy: deploy-publisher deploy-consumer

# Push Image to local Kind Cluster
push-to-kind: build
	@echo "Loading publisher docker image into kind cluster...."
	@kind load docker-image $(PUBLISHER_IMAGE_NAME):$(PUBLISHER_IMAGE_TAG)
	@echo "Publisher Image has been pushed successfully into kind cluster."
	@echo "Loading consumer docker image into kind cluster...."
	@kind load docker-image $(CONSUMER_IMAGE_NAME):$(CONSUMER_IMAGE_TAG)
	@echo "Consumer Image has been pushed successfully into kind cluster."