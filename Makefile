
build: cmd/sample-go-app/main.go
	go build ./...

clean:
	rm sample-go-app

security: trivy

trivy:
	trivy fs .

gofmt:
	gofmt -w .
