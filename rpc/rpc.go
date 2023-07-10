package rpc

type RPC struct {
	Module string
}

func NewRPC(module string) RPC {
	return RPC{
		Module: module,
	}
}