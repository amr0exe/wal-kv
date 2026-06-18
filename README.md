# wal-kv

## Info
In-memory key-value store with WAL persistence and primary-replica replication.

**Primary** handles SET/GET/DEL and replicates mutations to replicas. **Replicas** serve GET requests and stay in sync via snapshot on boot + real-time replication from primary.

## Quick Start

```bash
make build           # build the binary

# Terminal 1 — Primary
make primary         # :8080 (HTTP) + :5001 (gRPC)

# Terminals 2-4 — Replicas
make r1              # :8081 (HTTP) + :5002 (gRPC) → primary :5001
make r2              # :8082 (HTTP) + :5003 (gRPC) → primary :5001
make r3              # :8083 (HTTP) + :5004 (gRPC) → primary :5001
```

## Test

```bash
# Write to primary — replicates to all replicas
curl -X PUT localhost:8080/kv/foo -d '{"value":"bar"}'

# Read from any replica
curl localhost:8081/kv/foo
curl localhost:8082/kv/foo
curl localhost:8083/kv/foo

# Delete from primary — also replicated
curl -X DELETE localhost:8080/kv/foo

# Replicas reject writes
curl -X PUT localhost:8081/kv/x -d '{"value":"y"}'   # 405

# Stop all
make kill
```

## Manual Node Setup

```bash
# Primary
./app --node=primary --http-port=:8080 --port=:5001 --replicas=localhost:5002,localhost:5003,localhost:5004

# Replica 1
./app --node=replica --http-port=:8081 --port=:5002 --primary=localhost:5001
```

## Demo
<video src="https://github.com/user-attachments/assets/c7f27f2d-7cdc-4afa-8748-1a7867202fcd" width="600" autoplay />
