package main

import "fmt"

type RaftResult string

const (
	// RaftResultAcked s sent when the Follower successfully acknowledges the Leader
	RaftResultAcked RaftResult = "RaftResultAcked"

	// RaftResultLowerTerm means the Sender was rejected because they had a lower term
	RaftResultLowerTerm RaftResult = "RaftResultLowerTerm"

	// RaftResultStaleLeader is sent when the Follower refuses to acknowledege the RPC request from
	// a leader of a previous term.
	RaftResultStaleLeader RaftResult = "RaftResultStaleLeader"

	// RaftResultRejectedLeader means that the Sender was not acknowledged by a node (Follower, Leader, Candidate)
	// as the leader of a current term. This can happen if there was a Split brain or network partition
	// due to the sender having logs that are not up to date
	RaftResultRejectedLeader RaftResult = "RaftResultRejectedLeader"

	// RaftResultLogsOutOfSync is sent when the previousLogIndex of the Follower and Leader do not match
	RaftResultLogsOutOfSync RaftResult = "RaftResultLogsOutOfSync"

	// RaftResultUnknownUnhandled accounts for situations that are unexpected or unhandled
	RaftResultUnknownUnhandled RaftResult = "RaftResultUnknownUnhandled"

	RaftResultSnapshot RaftResult = "RaftResultSnapshot"

	// RaftResultVoteDenied means that a Candidate was denied a vote
	RaftResultVoteDenied  RaftResult = "RaftResultVoteDenied"

	// RaftResultVoteGranted means that a Candidate was granted a vote
	RaftResultVoteGranted RaftResult = "RaftResultVoteGranted"
)

type LogStatus int

const (
	// LogStatusOutOfSync signifies that our local logs don't match with the leaders, which could
	// be by log index or log term
	LogStatusOutOfSync LogStatus = iota

	// LogStatusMatch signifies that our local logs match  with the leaders including the leader commits
	LogStatusMatch

	// LogStatusUpdateCommit signifies that we need to update our commited logs to match with the leaders'
	LogStatusUpdateCommit
)

func (rr RaftResult) String() string {
	switch rr {
	case RaftResultAcked:
		return "RaftResultAcked"
	case RaftResultLowerTerm:
		return "RaftResultLowerTerm"
	case RaftResultStaleLeader:
		return "RaftResultStaleLeader"
	case RaftResultRejectedLeader:
		return "RaftResultRejectedLeader"
	case RaftResultLogsOutOfSync:
		return "RaftResultLogsOutOfSync"
	case RaftResultUnknownUnhandled:
		return "RaftResultUnknownUnhandled"
	case RaftResultSnapshot:
		return "RaftResultSnapshot"
	case RaftResultVoteDenied:
		return "RaftResultVoteDenied"
	case RaftResultVoteGranted:
		return "RaftResultVoteGranted"
	default:
		msg := fmt.Sprintf("unexpected RaftResult: %s", string(rr))
    return msg
	}
}

func (ll LogStatus) String() string {
	switch ll {
	case LogStatusMatch:
		return "LogStatusMatch"
	case LogStatusOutOfSync:
		return "LogStatusOutOfSync"
	case LogStatusUpdateCommit:
		return "LogStatusUpdateCommit"
	default:
		panic(fmt.Sprintf("unexpected main.LogStatus: %#v", ll))
	}
}
