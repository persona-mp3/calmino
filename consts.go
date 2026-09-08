package main

import "time"

var (
	HeartbeatInterval   = uint(800)
	ElectionIntervalMin = uint(1900)
	ElectionIntervalMax = uint(2500)

	// SendChanTimeout is the maximum amount of time that a goroutine
	// should spend sending on a channel before dropping the value or assuming there's
	// no receiver on the other end
	SendChanTimeout = 400 * time.Millisecond
)
