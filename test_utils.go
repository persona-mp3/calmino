package main

import (
	"fmt"
	db "jraft/database"

	"pgregory.net/rapid"
)

func generateLogEntries(rt *rapid.T, min, max uint64) []*Log {
	logSize := rapid.Uint64Range(min, max).Draw(rt, "logSize")
	logEntries := make([]*Log, logSize)
	for i := range logSize {
		logEntries[i] = &Log{
			Index: i,
			Term:  rapid.Uint64().Draw(rt, "logTerm"),
			Data: db.KV{
				Command: rapid.SampledFrom(
					[]db.Command{db.CommandGet, db.CommandSet, db.CommandRemove},
				).Draw(rt, fmt.Sprintf("cmd-%d", i)),
				Key:   rapid.StringMatching(`[\x20-\x7E]+`).Draw(rt, fmt.Sprintf("key-%d", i)),
				Value: rapid.StringMatching(`[\x20-\x7E]+`).Draw(rt, fmt.Sprintf("val-%d", i)),
			},
		}
	}

	return logEntries
}
