install:
	packr2
	go install -v ./cmd/buffalo

test:
	go test -short -cover ./...

cov:
	go test -short -coverprofile cover.out ./...
	go tool cover -html cover.out
