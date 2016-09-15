build:
	go build main.go

.PHONY: test_docker
test_docker:
	go test -v ./cmd/docker --email="${email}" --password="${password}" --username="${username}"
