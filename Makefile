dating-app:
	go run cmd/dating-app/main.go

test:
	go test ./... -coverprofile coverage.cov
	go tool cover -func coverage.cov
	rm coverage.cov