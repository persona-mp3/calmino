package main

import (
	"log"
)

func main() {
	// node := NewNode(
	// 	"",
	// 	"localhost:8990",
	// 	[]string{},
	// 	nil,
	// 	nil,
	// 	nil,
	// 	Configuration{Out: os.Stdout, LogFormat: 1})
	// _ = node

	rawConfig, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := RunCluster(rawConfig); err != nil {
		log.Fatal(err)
	}

}
