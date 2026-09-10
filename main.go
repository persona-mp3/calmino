package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var nodeId int = 1

func main() {
	flag.IntVar(
		&nodeId,
		"nodeId", nodeId,
		`
		nodeId is used to run a single node from the addr. 
		By default it selects the first one and is 0 based
		`,
	)
	flag.Parse()
	rawConfig, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, syscall.SIGKILL)
	defer cancel()

	if isFlagPassed("nodeId") {
		nodeId = nodeId - 1
		if nodeId < 0 {
			// since we're using 1based index, automatically run the first one if user sepcified 0
			nodeId += 1
			log.Println("[warn] index is not 0 based, running first node instead")
		}
		if err := runSingleNodeById(ctx, nodeId, rawConfig); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := RunCluster(ctx, rawConfig); err != nil {
			log.Fatal(err)
		}
	}

}

func runSingleNodeById(ctx context.Context, id int, rc *RawConfig) error {
	if id >= len(rc.Addrs) || id >= len(rc.HTTPPprofAddr) {
		return fmt.Errorf(
			`make sure addrs or http_pprof_addr have the same amount of values set or nodeId is not larger or equal to thier size`,
		)
	}

	if len(rc.Addrs) != len(rc.HTTPPprofAddr) {
		return fmt.Errorf("addrs and pproff addr must be of equal length")
	}

	allNodes := createNodes(rc)
	node := allNodes[id]
	pprofAddr := rc.HTTPPprofAddr[id]
	go func() {
		log.Printf("[pprof] (%d) starting pprof server at http://%s\n", id+1, pprofAddr)
		err := startPprofServer(pprofAddr)
		if err != nil {
			log.Println("[error] could not start pprof server ", err)
		}
	}()

	log.Println("starting single node process")
	return node.Start(ctx)
}
