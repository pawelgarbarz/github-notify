build:
	go mod tidy && go build -o bin/github-notify

run:
	bin/github-notify get

test:
	go mod tidy && go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.48.0

lint:
	golangci-lint run

staticcheck-install:
	go install honnef.co/go/tools/cmd/staticcheck@latest

staticcheck:
	staticcheck ./...