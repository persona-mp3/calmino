package main

import "net/rpc"

// RPCConn provides an abstraction over rpc.Client
type RPCConn interface {
	Id() string
	Addr() string
	Close() error
	Call(serviceName string, req any, reply any) error
}

type RPCPeer struct {
	// id holds the current id of the connection
	id string
	// addr holds the address this peer is connected to
	addr string
	// conn is the underlying *rpc.Client connection
	conn *rpc.Client
}

func (rp RPCPeer) Id() string {
	return rp.id
}

func (rp RPCPeer) Addr() string {
	return rp.addr
}

func (rp RPCPeer) Call(serviceMethodName string, req any, reply any) error {
	return rp.conn.Call(serviceMethodName, req, reply)
}

func (rp RPCPeer) Close() error {
	return rp.conn.Close()
}
func NewRPCPeer(id, addr string, conn *rpc.Client) RPCConn {
	if conn == nil {
		panic("cannot pass in nil conn to NewRPCPeer")
	}
	return RPCPeer{id: id, addr: addr, conn: conn}
}
