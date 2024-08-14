SHELL := /bin/bash

.PHONY: run
run:
	go run cmd/web_server_demo/main.go

.PHONY: build
build:
	go build -o bin/web_server_demo cmd/web_server_demo/main.go

.PHONY: build_and_run
build_and_run: build
	./bin/web_server_demo
