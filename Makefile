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
.PHONY: all