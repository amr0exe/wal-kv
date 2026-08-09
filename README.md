# plainmap

It's a key-value store, designed to run on a multi-node setup, splitting it's reads between replica nodes while primary node handles read, write, delete. It also has custom Write-Ahead-Log(WAL) implementation for storing each action in a log file, which allows for recoverability of key-values in-memory.

## Workflow
<img width="1281" height="612" alt="Image" src="https://github.com/user-attachments/assets/3af2af89-9705-4cd4-8146-f27757571418" />


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

To Write: collect the fields and write them sequentially based on previous WAL protocol. 
To Read: reading is more involved, you have to parse based on each part of protocol, for now, ignore Seq, So, read Op hold it, then read key_len read that amount of bytes from there follow same for value, then recompute the checksum over all preceding data to verify integrity before handing the entry to the caller for replay.

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
