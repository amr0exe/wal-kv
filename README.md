# wal-kv

A minimal distributed key-value store built from scratch to explore the challenges of distributed systems. Uses Write-Ahead Logging (WAL) for persistence and a primary-replica replication model, no heavy dependencies, no external databases.

## Architecture
A single primary node handles all write operations (SET, DEL). Replica nodes serve read requests only. On startup, each replica fetches a full snapshot of the key-value state from the primary via gRPC streaming. Thereafter, every mutation performed on the primary is pushed in real time to all registered replicas over the same gRPC layer. Sequence numbers assigned by the primary guarantee that operations are applied in order across the cluster.

## Workflow
<img width="1281" height="490" alt="Image" src="https://github.com/user-attachments/assets/3af2af89-9705-4cd4-8146-f27757571418" />


## WAL Format

Each log entry is persisted in the following binary layout:

```
Seq | Op | key_len | key | value_len | value | checksum
```

- **Seq**: sequence number for global ordering across nodes
- **Op**: operation type (`SET` or `DEL`)
- **key_len/key**: variable length key
- **value_len/value**: variable-length value
- **checksum**: basic integrity check for corruption detection

Writing is straightforward: collect the fields and write them sequentially. Reading is more involved, parse Seq, Op, key_len and key bytes, value_len and value bytes, then recompute the checksum over all preceding data to verify integrity before handing the entry to the caller for replay.

Only the primary maintains a WAL on disk. Replicas hold state entirely in memory, avoiding file conflicts when multiple nodes run on the same machine.

## Backstory

I wanted persistence for a KV store without pulling in a full database like Postgres. This was overkill for the scope of the project. Having recently learned about Write-Ahead Logging, I decided to implement one from scratch rather than rely on an existing library.

The hardest part was the transition from a single-node to a multi-node design. The project was originally conceived as a single-node store, and I only later realised how deeply that assumption was baked into the code. Decoupling the instantiation logic, separating concerns, and redistributing responsibilities that were once comfortably interleaved proved to be a significant effort.

I knew I would use gRPC for inter-node communication and that a primary-replica topology was the right fit, the primary handles mutations while replicas serve reads. Knowing the shape, however, and actually getting there were two very different things.

## Specific challenges

- Running on multiple nodes required more than just assigning different ports. Handler logic and storage had to be fully split between primary and replica roles.
- Factory code (`NewRouter`, `NewStore`, etc.) had to be redesigned. Only the primary performs WAL recovery. Replicas accept only GET requests, rather than explicitly rejecting others, the replica router simply does not serve them.

## Replication

Two gRPC services drive the replication mechanism:

- **GetSnapshot**: invoked by a replica during boot. The primary streams its entire key-value state to the joining replica.
- **ApplyMutation**: invoked by the primary after every `SET` or `DEL` operation. Sends the mutation to every registered replica, which applies it locally while retaining the primary's sequence number.

Replica addresses are currently provided to the primary via the `--replicas` command-line flag.

## Test Procedure
```
make build          # Build executable

make primary        # Start Primary Node

make r1             # Start Replica Node
make r2
make r3

# TEST
sudo snap install httpie --classic  (Optional) for curl

http PUT :8080/kv/foo value=bar
http GET :8080/kv/foo

# Try read on replica nodes
http GET :8081/kv/foo
http GET :8082/kv/foo
http GET :8083/kv/foo

```

## Demo

<video src="https://github.com/user-attachments/assets/c7f27f2d-7cdc-4afa-8748-1a7867202fcd" width="600" autoplay muted />
