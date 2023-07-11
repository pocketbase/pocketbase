lint:
	golangci-lint run -c ./golangci.yml ./...

test:
	go test ./... -v --cover

jsvmdocs:
	go run ./plugins/jsvm/internal/docs/docs.go

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out
