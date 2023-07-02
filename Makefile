build:
	go build -o bin/api

run: build
	./bin/api	

seed:
	go run scripts/seed.go

test:
	go test ./... -count=1

mongodb:
	docker run --name mongodb -p 27017:27017 -d mongo:latest 

.PHONY: build run test