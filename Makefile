dev:
	air

templ:
	templ generate --watch --proxy="http://localhost:8080"

test:
	go test ./...