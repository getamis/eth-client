package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/p2p"
	client "github.com/ethereum/go-ethereum/rpc"
)

//go:generate mockgen -source=admin.go -destination=mock_admin.go -package=rpc
type Admin interface {
	PrivateAdmin
	PublicAdmin
}

type admin struct {
	PrivateAdmin
	PublicAdmin
}

func NewAdmin(client *client.Client) Admin {
	return &admin{
		PrivateAdmin: NewPrivateAdmin(client),
		PublicAdmin:  NewPublicAdmin(client),
	}
}

type PrivateAdmin interface {
	// AddPeer connects to the given nodeURL.
	AddPeer(ctx context.Context, nodeURL string) (bool, error)
	// RemovePeer disconnects from a remote node if the connection exists
	RemovePeer(ctx context.Context, nodeURL string) (bool, error)
	// ImportChain imports a blockchain from a local file.
	ImportChain(ctx context.Context, file string) (bool, error)
	// ExportChain exports the current blockchain into a local file.
	ExportChain(ctx context.Context, file string) (bool, error)
	// StartRPC starts the HTTP RPC API server.
	StartRPC(ctx context.Context, host string, port int, cors string, apis string) (bool, error)
	// StopRPC terminates an already running HTTP RPC API endpoint.
	StopRPC(ctx context.Context) (bool, error)
	// StartWS starts the websocket RPC API server.
	StartWS(ctx context.Context, host string, port int, allowedOrigins string, apis string) (bool, error)
	// StopRPC terminates an already running websocket RPC API endpoint.
	StopWS(ctx context.Context) (bool, error)
}

type privateAdmin struct {
	client *client.Client
}

func NewPrivateAdmin(client *client.Client) PrivateAdmin {
	return &privateAdmin{
		client: client,
	}
}

// AddPeer connects to the given nodeURL.
func (pri *privateAdmin) AddPeer(ctx context.Context, nodeURL string) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_addPeer", nodeURL)
	if err != nil {
		return false, err
	}
	return r, nil
}

// RemovePeer disconnects from a remote node if the connection exists
func (pri *privateAdmin) RemovePeer(ctx context.Context, nodeURL string) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_removePeer", nodeURL)
	if err != nil {
		return false, err
	}
	return r, nil
}

// ImportChain imports a blockchain from a local file.
func (pri *privateAdmin) ImportChain(ctx context.Context, file string) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_importChain", file)
	if err != nil {
		return false, err
	}
	return r, nil
}

// ExportChain exports the current blockchain into a local file.
func (pri *privateAdmin) ExportChain(ctx context.Context, file string) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_exportChain", file)
	if err != nil {
		return false, err
	}
	return r, nil
}

// StartRPC starts the HTTP RPC API server.
func (pri *privateAdmin) StartRPC(ctx context.Context, host string, port int, cors string, apis string) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_startRPC", host, port, cors, apis)
	if err != nil {
		return false, err
	}
	return r, nil
}

// StopRPC terminates an already running HTTP RPC API endpoint.
func (pri *privateAdmin) StopRPC(ctx context.Context) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_stopRPC")
	if err != nil {
		return false, err
	}
	return r, nil

}

// StartWS starts the websocket RPC API server.
func (pri *privateAdmin) StartWS(ctx context.Context, host string, port int, allowedOrigins string, apis string) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_startWS", host, port, allowedOrigins, apis)
	if err != nil {
		return false, err
	}
	return r, nil
}

// StopRPC terminates an already running websocket RPC API endpoint.
func (pri *privateAdmin) StopWS(ctx context.Context) (bool, error) {
	var r bool
	err := pri.client.CallContext(ctx, &r, "admin_stopWS")
	if err != nil {
		return false, err
	}
	return r, nil
}

type PublicAdmin interface {
	// Peers returns the connected peers.
	Peers(ctx context.Context) ([]*p2p.PeerInfo, error)
	// NodeInfo gathers and returns a collection of metadata known about the host.
	NodeInfo(ctx context.Context) (*p2p.NodeInfo, error)
	// Datadir retrieves the current data directory the node is using.
	Datadir(ctx context.Context) (string, error)
}

type publicAdmin struct {
	client *client.Client
}

func NewPublicAdmin(client *client.Client) PublicAdmin {
	return &publicAdmin{
		client: client,
	}
}

// Peers returns the connected peers.
func (pub *publicAdmin) Peers(ctx context.Context) ([]*p2p.PeerInfo, error) {
	var r []*p2p.PeerInfo
	err := pub.client.CallContext(ctx, &r, "admin_peers")
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NodeInfo gathers and returns a collection of metadata known about the host.
func (pub *publicAdmin) NodeInfo(ctx context.Context) (*p2p.NodeInfo, error) {
	var r *p2p.NodeInfo
	err := pub.client.CallContext(ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Datadir retrieves the current data directory the node is using.
func (pub *publicAdmin) Datadir(ctx context.Context) (string, error) {
	var r string
	err := pub.client.CallContext(ctx, &r, "admin_datadir")
	if err != nil {
		return "", err
	}
	return r, nil
}
