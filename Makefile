install:
	go install -v ./cmd/buffalo

test:
	go test -failfast -short -cover ./...

cov:
	go test -short -coverprofile cover.out ./...
	go tool cover -html cover.out
