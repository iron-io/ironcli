build:
	go build main.go

.PHONY: test_base
test_base:
	go test -v ./cmd

.PHONY: test_docker
test_docker:
	go test -v ./cmd/docker

.PHONY: test_lambda
test_lambda:
	go test -v ./cmd/lambda

.PHONY: test_mq
test_mq:
	go test -v ./cmd/mq

.PHONY: test_worker
test_worker:
	go test -v ./cmd/worker
