build:
	@go build -o app cmd/server/main.go

primary:
	@go run cmd/server/main.go --node=primary --http-port=:8080 --port=:5001

r1:
	@go run cmd/server/main.go --node=replica --http-port=:8081 --primary=localhost:5001

r2:
	@go run cmd/server/main.go --node=replica --http-port=:8082 --primary=localhost:5001

r3:
	@go run cmd/server/main.go --node=replica --http-port=:8083 --primary=localhost:5001

kill:
	@pkill -f "cmd/server/main.go" 2>/dev/null || true
