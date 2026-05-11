### lest-memory-storage

#### Benchmarking along with Redis using `redis-benchmark`

Running Lest:
```redis-benchmark -h 127.0.0.1 -p 8090 -t set -c 1 -n 10000 -d 64 -q```

```memtier_benchmark -s 127.0.0.1 -p 8090 --protocol=redis --command="SET __key__ __data__" -t 1 -c 1```

Running Redis:
```redis-benchmark -h 127.0.0.1 -p 6379 -t set -c 1 -n 10000 -d 64 -q```

```memtier_benchmark -s 127.0.0.1 -p 6379 --protocol=redis --command="SET __key__ __data__" -t 1 -c 1```

# FIRST WORKING VERSION BENCHMARK TESTS:
## 1 connection 1 thread SET commands
### Redis
~45000op/s
### Lest
~7500op/s
#### While the result of first version is quite impressive, the program does too much logging, creating Correlation Id per Request + Connection Id per connection.
### Lest. Logging with Correlation Id removed
~39000op/s

### Conclusion
First implementation handles an impressive amount of load. While having only one basic Map storage with single Mutex lock it does a good work with 1 concurrent connection and shows a solid degradation of it's operation speed with multiple concurrent connections.
Also, it's really important to do wise, quiet logging of exclusively the most important information
