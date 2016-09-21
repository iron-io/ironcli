build:
	go build main.go

.PHONY: test_docker
test_docker:
	go test -v ./cmd/docker

.PHONY: test_lambda
test_lambda:
	go test -v ./cmd/lambda
