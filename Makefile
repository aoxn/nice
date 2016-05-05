#!/usr/bin/env bash

REPO=133.130.96.46:5000
PROJECT=spacexnice
APP=nice
VERSION=1.0.0

IMAGE=${REPO}/${PROJECT}/${APP}:${VERSION}


all:
	bash build.sh
	docker build -t ${IMAGE} build
	docker push ${IMAGE}

tag:
	docker build -t ${IMAGE} build
	docker push ${IMAGE}
push:
	docker push ${IMAGE}
run:
	docker run -d -p 8000:8000 -v /data/registry:/data/registry -e DB_DATA_PATH=/data/registry/ -e GIN_MODE=release ${IMAGE}
.PHONY: all