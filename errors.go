package main

import (
	"errors"
	"fmt"
)

type CandidateError struct {
	term uint64
	err  error
	info any
}

var ErrElectingForPastTerm = errors.New("running electing for past term. this should never happen")

func NewCandidateError(term uint64, err error, info ...any) CandidateError {
	return CandidateError{
		term: term,
		err:  err,
		info: info,
	}
}

func (ce CandidateError) Error() string {
	return fmt.Sprintf("%s. currentTerm: %d diagnostics: %+v", ce.err, ce.term, ce.info)
}
