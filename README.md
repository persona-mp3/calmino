## 🚀 calmino

calmino is an ongoing refactor of [fsm](https://github.com/persona-mp3/fsm.git)


![Tests](https://github.com/persona-mp3/calmino/actions/workflows/test.yml/badge.svg)

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


### Run tests
calmino makes use of property based testing and normal unit testing. It is a goal 
to acheive DST in some parts of the codebase. The library used is rapid, and it 
offers reproduceablity if errors are found. If you run a test and they fail, rapid 
allows you to replay that same test with all data, please raise an issue or PR about 
it
```
go test
```

# Configuration
To configure the Raft Cluster please see [config.sample.toml](./config.sample.toml)
