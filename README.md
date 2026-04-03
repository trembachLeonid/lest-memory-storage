### lest-memory-storage

#### Benchmarking along with Redis using `redis-benchmark`

Running Lest:
```redis-benchmark -h 127.0.0.1 -p 8090 -t set -c 1 -n 10000 -d 64 -q```

Running Redis:
```redis-benchmark -h 127.0.0.1 -p 6379 -t set -c 1 -n 10000 -d 64 -q```