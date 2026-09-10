By default, calmino looks for a toml file called `config.toml` in the current 
directory it's being executed, otherwise, it refuses to start. The configuration 
is limited by default but allows for creating a cluster and setting some raft 
configurations, along with some other things

- Mode: Allows for different ways to run a cluster.

1.  To run a whole cluster in a single process, set mode to `single_process`
2. To run a whole cluster in different processes, set mode to `multi_process`

By default this is set to `single_process`, which allows for easier profiling 
and trouble shooting when diagnosing problems

```toml
mode = "single_process"
```

Note that by default, you cannot specify to run just a single node with this setup.
This has been addressed by providing the `nodeId` flag to calmino
```bash
./calmino --nodeId <index-of-node>
```

The `index-of-node` is one based index, so providing `1` to calmino will run the 
create a server at the very first address provided in the `addr` field, and also 
start a pprof server  at the very first address provided in  `http_pprof_addr`.
To eliminate edgecases, both the `addr` and `http_pprof_addr` must have the same 
amount of addresses provided. NO CHECKS ARE PROVIDED IF THEY SHARE THE SAME ADDR

Also note that, unless you run all other nodes individually, consensus will not 
be acheived. 
This is the same thing `multi_process` does under the hood, but takes care of the 
processes


- Addrs: These are addresses that each node spawned will listen on. Each 
node will also try to connect to other servers provided in this field
```toml
addrs = ["localhost:4000", "localhost:4001", "localhost:4002"]
```

- Log Format: You can specify what format you want the logs to be, either normal 
logging or JSON structured logging
```toml
# 0 represents normal text logging and 1 represents json structued logging
log_format = 0
```

- Log Persistence: By default each node spawned writes to it's own log file, this 
can be configured with  the `persist` flag. If this is set to false, all logs 
will all be redirected to stdout. Except when the cluster is ran in `multi_process`
mode. Then all logs will be lost. It's recommended to leave it as `true`
```toml
# persist writes logs to file named log-file-<node-id>
persist = true
```

- Log Level: Follows Go's default logging levels, 
```toml
# log_level starts at -4 DEBUG, 0 INFO, 4 WARN, 8 ERROR
log_level = -4
```
