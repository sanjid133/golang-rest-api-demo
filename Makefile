IMAGE_REPO = user
IMAGE_VERSION ?= v1
IMAGE_NAME = sanjid133/$(IMAGE_REPO):$(IMAGE_VERSION)

.PHONY: run build push

run:
	docker-compose up --build user

stop:
	docker-compose down

build:
	docker build --no-cache -t $(IMAGE_NAME) .

push:
	docker push $(IMAGE_NAME)

default: run