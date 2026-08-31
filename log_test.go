package main

import (
	"fmt"
	db "jraft/database"
	"testing"

	"github.com/magiconair/properties/assert"
	"pgregory.net/rapid"
)

func TestNewLogStore(t *testing.T) {
	logStore := NewLogStore()
	assert.Equal(
		t,
		logStore.CommitIndex(),
		uint64(0), "commit indexes don't match for initialized logStore")
}

func TestLogStorePreviousLogEntry(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		logSize := rapid.Uint64Range(0, 100).Draw(rt, "logSize")
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

		var prevLogEntry Log
		if len(logEntries) < 2 {
			prevLogEntry = Log{}
		} else {
			prevLogEntry = *logEntries[logSize-2]
		}

		logStore := NewLogStore()
		for _, log := range logEntries {
			logStore.Append(log)
		}

		actualPrevLogEntry := logStore.PreviousEntry()
		assert.Equal(
			t,
			prevLogEntry,
			actualPrevLogEntry,
			"prevLogEntry and actualPrevLogEntry do not match",
		)

	})
}
