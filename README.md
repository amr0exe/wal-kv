# wal_kv

`wal_kv` is an in-memory Key-Value store, which mimics Write-Ahead-Log(wal), resulting in persistence through wal.log file.

---

```bash
sudo apt install httpie         # opt, for curl

make                            # run the app

http PUT :8080/kv/foo value=bar # SET api
http GET :8080/kv/foo           # GET api
http DELETE :8080/kv/foo        # DEL api
```
