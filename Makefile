build:
	@go build -o app cmd/server/main.go

primary: build
	@./app --node=primary --http-port=:8080 --port=:5001 --replicas=localhost:5002,localhost:5003,localhost:5004

r1: build
	@./app --node=replica --http-port=:8081 --port=:5002 --primary=localhost:5001

r2: build
	@./app --node=replica --http-port=:8082 --port=:5003 --primary=localhost:5001

r3: build
	@./app --node=replica --http-port=:8083 --port=:5004 --primary=localhost:5001

kill:
	@pkill -f "./app" 2>/dev/null; exit 0
