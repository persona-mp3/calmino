package jraft

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
	Append(log Log) error

	// Commited returns the index that has be written to the database
	Commited() uint64
}

type Logs struct {
	logs      []*Log
	commitIdx uint64
}

