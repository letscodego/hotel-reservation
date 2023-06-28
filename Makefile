build:
	go build -o bin/api

run: build
	./bin/api	

test:
	go test ./... -count=1

.PHONY: build run test