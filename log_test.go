package main

import (
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

func TestLogStoreAppendAndPreviousLogEntry(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		logEntries := generateLogEntries(rt, 0, 100)
		logSize := uint64(len(logEntries))

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

func TestLogStoreSnapshotFrom(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		logEntries := generateLogEntries(rt, 0, 100)
		logSize := uint64(len(logEntries))

		logStore := NewLogStore()
		for _, log := range logEntries {
			logStore.Append(log)
		}

		uint64NumberGen := rapid.Uint64Range(0, 99)
		startIdx := uint64NumberGen.Draw(rt, "startIdxForSnapshotFrom")
		var expectedSnapshot []Log
		var expectedErr error

		clonedLogs := make([]Log, 0, logSize)
		for _, log := range logEntries {
			clonedLogs = append(clonedLogs, *log)
		}
		if startIdx > logSize {
			expectedErr = ErrIndexNotFound
			expectedSnapshot = clonedLogs
		} else {
			expectedSnapshot = clonedLogs[startIdx:]
		}

		actualSnapshot, actualErr := logStore.SnapshotFrom(startIdx)
		assert.Equal(t, actualSnapshot, expectedSnapshot, "snapshots dont match")
		assert.Equal(t, actualErr, expectedErr, "snapshots errors dont match")

	})
}
