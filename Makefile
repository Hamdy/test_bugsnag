# Settings come from .env (copy .env.example): BUGSNAG_API_KEY above all.
-include .env
export

.PHONY: run test build

run:
	go run .

test:
	go test ./...

build:
	go build -o test_bugsnag .
