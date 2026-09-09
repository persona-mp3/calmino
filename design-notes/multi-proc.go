package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill)
	defer cancel()
	cmds := []*exec.Cmd{}
	for i := range 8 {
		cmd := exec.CommandContext(ctx, "./calmino", "--id", strconv.Itoa(i))
		if err := cmd.Start(); err != nil {
			fmt.Println("err::", i, err)
			continue
		}
		cmds = append(cmds, cmd)
	}
	// so we either want to wait for each one to return to us
	// or we want to easily cancel each one, should we use wg?
	// so a user can simply just ctrlC from here, and all others will obey, (hopefully)
	// so instead, we simply just want to leverage the control-plane here

	// collect stdin
	// cmds:: kill node 3
	// cmds:: restart node 3

}
