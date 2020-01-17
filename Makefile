install:
	go install -v -tags "sqlite" ./cmd/buffalo
	go mod tidy -v

test:
	go test -failfast -short -cover -tags "sqlite" ./...
	cd ./plugins && go test -failfast -short -cover -tags "sqlite" ./...
	go mod tidy -v

cov:
	go test -short -coverprofile cover.out -tags "sqlite" ./...
	go tool cover -html cover.out
	go mod tidy -v
