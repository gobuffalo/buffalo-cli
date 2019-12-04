install:
	packr2
	go install -v ./cmd/buffalo

test:
	go test -short -cover ./...
