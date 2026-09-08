## 🚀 calmino
![Tests](https://github.com/persona-mp3/calmino.git/actions/workflows/test.yml/badge.svg)
calmino is an ongoing refactor of [fsm](https://github.com/persona-mp3/fsm.git)

### Prerequistes
- [Go](https://golang.org/) installed on your machine.

### Running the application
```bash
go build .
./calmino
```

### Low level profiling
Go's http pprof has already been integrated by default and is running on [http://localhost:6061/debug](http://localhost:6061/debug) 
To change what port it connects to, you can edit the config file 
```bash
http_pprof_addr = "localhost:port_number"
```

# Configuration
To configure the Raft Cluster please see [config.sample.toml](./config.sample.toml)
