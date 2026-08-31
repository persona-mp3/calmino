package main

import "fmt"

func RunCluster(rc *RawConfig) error {
	if rc.Mode == singleProcess {
		singleProcessCluster(rc)
		return nil
	}
	return fmt.Errorf("%s not yet implemented. please use single_process", string(rc.Mode))
}

func singleProcessCluster(rc *RawConfig) {
	// 1. extract peers and addresses, possibly a map[id][]string{}
	idPeers := buildIndexedPeerMap(rc.Addrs)
	for idx, info := range idPeers {
		for addr, peer := range info {
			fmt.Printf("(%d) %s ->  %v\n", idx, addr, peer)
		}
	}
	// CURRENTLY: 2. Basically create new nodes, but also create their configs

}

func singleNodeCluster(rc *RawConfig) {}
