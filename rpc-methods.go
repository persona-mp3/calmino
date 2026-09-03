package main

import "log/slog"

func (s *Server) SnapshotRPC(req SnapshotRequest, reply *SnapshotReply) error {
	res := make(chan RPCReply, 1)
	payload := RPCPayload{kind: RPCKindSnapshot, payload: req, reply: res}
	s.network <- payload

	data := <-res
	switch data := data.payload.(type) {
	case *SnapshotReply:
		reply = data
	default:
		s.logger.Error(
			"invalid payload from node expected Snapshot", slog.Any("got", data),
		)

		panic("invalid payload recvd")
	}

	return nil
}

func (s *Server) VoteRPC(req VoteRequest, reply *VoteReply) error {
	res := make(chan RPCReply, 1)
	payload := RPCPayload{kind: RPCKindVote, payload: req, reply: res}
	s.network <- payload

	data := <-res
	switch data := data.payload.(type) {
	case *VoteReply:
		reply = data
	default:
		s.logger.Error(
			"invalid payload from node expected Vote", slog.Any("got", data),
		)

		panic("invalid payload recvd")
	}

	return nil
}
