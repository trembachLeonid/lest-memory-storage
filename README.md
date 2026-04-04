### lest-memory-storage

#### Benchmarking along with Redis using `redis-benchmark`

Running Lest:
```redis-benchmark -h 127.0.0.1 -p 8090 -t set -c 1 -n 10000 -d 64 -q```

```memtier_benchmark -s 127.0.0.1 -p 8090 --protocol=redis --command="SET __key__ __data__" -t 1 -c 1```

Running Redis:
```redis-benchmark -h 127.0.0.1 -p 6379 -t set -c 1 -n 10000 -d 64 -q```

```memtier_benchmark -s 127.0.0.1 -p 6379 --protocol=redis --command="SET __key__ __data__" -t 1 -c 1```