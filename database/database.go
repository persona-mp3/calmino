package database

import "fmt"

type Command string

const (
	CommandGet    Command = "get"
	CommandSet    Command = "set"
	CommandRemove Command = "rm"
)

// KV holds the data that is sent to the KV Storage engine
type KV struct {
	Command Command
	Key     string
	Value   string
}

func (k KV) String() string {
	return fmt.Sprintf("KV {command: %s, Key: %s, Value: %s}", string(k.Command), k.Key, k.Value)
}
