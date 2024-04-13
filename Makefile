lint:
	golangci-lint run

test: lint
	go test -cover -race ./internal/adapters/left/cli

build: test
	go build -o diego main.go

run: build
	./diego

deploy: build
	goreleaser build --snapshot --clean

debug: test
	go build -gcflags=all="-N -l" -o diego main.go


update:
	go get -u ./...
	go mod tidy
	make run

golden: lint
	go test -cover ./internal/adapters/left/cli -update
	make run
