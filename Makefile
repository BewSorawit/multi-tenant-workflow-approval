APP=multi-tenant-workflow-approval

format-check:
	gofmt -l .

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

docker-build:
	docker build -t $(APP) .

smoke-test:
	./scripts/smoke-test.sh

ci: format-check vet test build docker-build smoke-test