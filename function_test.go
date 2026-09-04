package main

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestLogCompletionTermsMatchWithHigerIndex(t *testing.T) {
	prevLogEntry := Log{Index: 0, Term: 9}
	reqTerm := uint64(9)
	reqIndex := uint64(8)

	status := logCompletion(prevLogEntry, reqIndex, reqTerm)
	assert.Equal(t, status, true)
}

func TestLogCompletionTermsMatchWithLowerIndex(t *testing.T) {
	prevLogEntry := Log{Index: 100, Term: 9}
	reqTerm := uint64(9)
	reqIndex := uint64(8)

	status := logCompletion(prevLogEntry, reqTerm, reqIndex)

	assert.Equal(t, status, false)
}

func TestLogCompletionWithLowerTerm(t *testing.T) {
	prevLogEntry := Log{Index: 100, Term: 3}
	reqTerm := uint64(9)
	reqIndex := uint64(8)

	status := logCompletion(prevLogEntry, reqTerm, reqIndex)

	assert.Equal(t, status, true)
}

func TestLogCompletionFromLowerTerm(t *testing.T) {
	prevLogEntry := Log{Index: 100, Term: 13}
	reqIndex := uint64(10)
	reqTerm := uint64(2)

	status := logCompletion(prevLogEntry, reqTerm, reqIndex)

	assert.Equal(t, status, false)
}
