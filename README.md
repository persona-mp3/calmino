## 🚀 calmino
![Tests](https://github.com/persona-mp3/calmino/actions/workflows/test.yml/badge.svg)

calmino is an ongoing refactor of [fsm](https://github.com/persona-mp3/fsm.git), 
a Raft Consensus implementation.



### Prerequistes
- [Go](https://golang.org/) installed on your machine.

### Clone the repository
```bash
git clone https://github.com/persona-mp3/calmino.git
```

### Running the application
```bash
go build .
./calmino
```

### Low level profiling
Go's http pprof has already been integrated by default and is running on [http://localhost:6061/debug](http://localhost:6061/debug).

To change what address the server runs on, you can change this via the config file,
```bash
http_pprof_addr = ["localhost:port_number"]
```


### Run tests
calmino uses a mixture of property based testing and normal unit testing. It is a goal 
to acheive some sort DST in some parts of the codebase. 


The Go library used is rapid, which means that failed tests can be reproduced. If
tests fail on you run, please feel free raise an issue or PR about 

```bash
go test
```

# Configuration
To configure the Raft Cluster please see [config.sample.toml](./config.sample.toml)
and [docs/configuration.md](./docs/configuration.md)
