package main

import (
	"fmt"
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

	config, err := rawConfig.ToConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}
