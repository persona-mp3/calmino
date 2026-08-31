package main

import (
	"log"
)

func main() {
	rawConfig, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := RunCluster(rawConfig); err != nil {
		log.Fatal(err)
	}

}
