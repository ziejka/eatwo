dev:
	air

dev-s:
	go run cmd/server/main.go

templ:
	templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/server/main.go"

templ-s:
	templ generate --watch --proxy="http://localhost:8080" 

test:
	go test ./...

sqlc: 
	sqlc generate

tailwind:
	npm run watch
