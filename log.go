package main

import (
	db "calmino/database"
	"errors"
	"sync"
)

var (
	ErrIndexNotFound = errors.New("Log index not found")
)

type Log struct {
	// Index holds the position of the entry in the logs
	Index uint64
	Term  uint64
	Data  db.KV
}

// LogStore is used to provide an interface for storing and retrieving logs
type LogStore interface {
	// PreviousLogEntry returns the second-to-last log. Returns an empty log
	// if there are no more than 2 logs
	PreviousEntry() Log

	// SnapshotFrom returns all the logs starting from startIndex. If not found,
	// returns all logs
	SnapshotFrom(startIndex uint64) ([]Log, error)

	// Append stores a log entry
	Append(log *Log) uint64

	// Commited returns the index that has be written to the database
	CommitIndex() uint64

	// Apply applies all the logs to the database and returns the index of the most
	// recently commited log
	Apply() (uint64, error)

	// FlushTill applies all logs till stopIndex
	FlushTill(stopIndex uint64) error

	LogAt(targetIdx uint64) (Log, error)
}

type Logs struct {
	mu        sync.Mutex
	logs      []*Log
	commitIdx uint64
}

func NewLogStore() LogStore {
	logStore := &Logs{
		mu:        sync.Mutex{},
		commitIdx: 0,
	}

	return logStore
}

// PreviousLogEntry returns the second-to-last log. Returns an empty log if there
// are no more than 2 logs
func (l *Logs) PreviousEntry() Log {
	l.mu.Lock()
	defer l.mu.Unlock()
	logSize := len(l.logs)
	if logSize < 2 {
		return Log{}
	}

	return *l.logs[logSize-2]
}

// FlushTill applies all logs till stopIndex and updates the commitIdx. If
// stopIndex is not found, it returns ErrIndexNotFound
func (l *Logs) FlushTill(stopIndex uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	panic("FlushTill not implemented yet")
}

func (l *Logs) Apply() (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	panic("Apply not implemented yet")
}

// SnapshotFrom returns all the logs starting from startIndex. If not found,
// returns all logs
func (l *Logs) SnapshotFrom(startIdx uint64) ([]Log, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	logSize := len(l.logs)
	if startIdx > uint64(logSize) {
		clonedLogs := make([]Log, 0, logSize)
		for _, log := range l.logs {
			clonedLogs = append(clonedLogs, *log)
		}
		return clonedLogs, ErrIndexNotFound
	}
	clonedLogs := make([]Log, 0, logSize-int(startIdx))
	for _, log := range l.logs[startIdx:] {
		clonedLogs = append(clonedLogs, *log)
	}
	return clonedLogs, nil

}

// Append takes appends a new log into the store log entries. It returns the
// index of the assigned to the log provided
func (l *Logs) Append(log *Log) uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	// since len() is one based we just assign the logSize as the index position
	// for the new log
	logSize := uint64(len(l.logs))
	log.Index = logSize
	l.logs = append(l.logs, log)
	return logSize
}

func (l *Logs) CommitIndex() uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.commitIdx
}

func (l *Logs) LogAt(targetIndex uint64) (Log, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if targetIndex >= uint64(len(l.logs)) {
		return Log{}, ErrIndexNotFound
	}
	return *l.logs[targetIndex], nil
}
