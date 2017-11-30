package client

import (
	"context"

	"github.com/ethereum/go-ethereum/p2p"
)

// AddPeer connects to the given nodeURL.
func (c *client) AddPeer(ctx context.Context, nodeURL string) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_addPeer", nodeURL)
	if err != nil {
		return false, err
	}
	return r, nil
}

// RemovePeer disconnects from a remote node if the connection exists
func (c *client) RemovePeer(ctx context.Context, nodeURL string) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_removePeer", nodeURL)
	if err != nil {
		return false, err
	}
	return r, nil
}

// Peers return the connected peers.
func (c *client) Peers(ctx context.Context) ([]*p2p.PeerInfo, error) {
	var r []*p2p.PeerInfo
	err := c.rpc.CallContext(ctx, &r, "admin_peers")
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (c *client) NodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := c.rpc.CallContext(ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Datadir retrieves the current data directory the node is using.
func (c *client) Datadir(ctx context.Context) (string, error) {
	var r string
	err := c.rpc.CallContext(ctx, &r, "admin_datadir")
	if err != nil {
		return "", err
	}
	return r, nil
}

// ImportChain imports a blockchain from a local file.
func (c *client) ImportChain(ctx context.Context, file string) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_importChain", file)
	if err != nil {
		return false, err
	}
	return r, nil
}

// ExportChain exports the current blockchain into a local file.
func (c *client) ExportChain(ctx context.Context, file string) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_exportChain", file)
	if err != nil {
		return false, err
	}
	return r, nil
}

// StartRPC starts the HTTP RPC API server.
func (c *client) StartRPC(ctx context.Context, host *string, port *int, cors *string, apis *string) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_startRPC", host, port, cors, apis)
	if err != nil {
		return false, err
	}
	return r, nil
}

// StopRPC terminates an already running HTTP RPC API endpoint.
func (c *client) StopRPC(ctx context.Context) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_stopRPC")
	if err != nil {
		return false, err
	}
	return r, nil

}

// StartWS starts the websocket RPC API server.
func (c *client) StartWS(ctx context.Context, host *string, port *int, allowedOrigins *string, apis *string) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_startWS", host, port, allowedOrigins, apis)
	if err != nil {
		return false, err
	}
	return r, nil
}

// StopRPC terminates an already running websocket RPC API endpoint.
func (c *client) StopWS(ctx context.Context) (bool, error) {
	var r bool
	err := c.rpc.CallContext(ctx, &r, "admin_stopWS")
	if err != nil {
		return false, err
	}
	return r, nil
}
