package main

import (
	db "jraft/database"
)

type Log struct {
	// Index holds the position of the entry in the logs
	Index uint64
	Term  uint64
	Data  db.KV
}

// LogStore is used to provide an interface for storing and retreiving logs
type LogStore interface {
	// PreviousLogEntry returns the second-to-last log. Returns an empty log
	// if there are no more than 2 logs
	PreviousEntry() Log

	// SnapshotFrom returns all the logs starting from startIndex. If not found,
	// returns all logs
	SnapshotFrom(startIndex uint64) []Log

	// Apppend stores a log entry
	Append(log Log) uint64

	// Commited returns the index that has be written to the database
	LastCommitIndex() uint64
}

type Logs struct {
	logs      []*Log
	commitIdx uint64
}

func (l Logs) PreviousLogEntry() Log {
	return Log{}
}

func (l Logs) SnapshotFrom(startIdx uint64) []Log {
	return []Log{}
}

func (l Logs) Append(log Log) uint64 {
	return 0
}

func (l Logs) LastCommitIndex() uint64 {
	return 0
}
