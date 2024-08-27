SHELL := /bin/bash

.PHONY: clean
clean:
	rm bin/*

.PHONY: run
run:
	go run ./cmd/web_server_demo/

.PHONY: build
build:
	go build -o bin/web_server_demo ./cmd/web_server_demo/

.PHONY: build_and_run
build_and_run: build
	./bin/web_server_demo

.PHONY: clean_db
clean_db:
	find  -type f -name "chirp_db*.json" -delete
