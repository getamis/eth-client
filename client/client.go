// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

type client struct {
	*ethclient.Client
	c *rpc.Client
}

func Dial(rawurl string) (Client, error) {
	c, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}
	return &client{
		Client: ethclient.NewClient(c),
		c:      c,
	}, nil
}

func (c *client) Close() {
	c.c.Close()
}

// ----------------------------------------------------------------------------
// eth

// SendRawTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (c *client) SendRawTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.SendTransaction(ctx, tx)
}

func (c *client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var r string
	err := c.c.CallContext(ctx, &r, "eth_blockNumber")
	if err != nil {
		return nil, err
	}
	h, err := hexutil.DecodeBig(r)
	return h, err
}

// ----------------------------------------------------------------------------
// admin

func (c *client) AddPeer(ctx context.Context, nodeURL string) error {
	var r bool
	// TODO: Result needs to be verified
	// The response data type are bytes, but we cannot parse...
	err := c.c.CallContext(ctx, &r, "admin_addPeer", nodeURL)
	if err != nil {
		return err
	}
	return err
}

func (c *client) AdminPeers(ctx context.Context) ([]*p2p.PeerInfo, error) {
	var r []*p2p.PeerInfo
	// The response data type are bytes, but we cannot parse...
	err := c.c.CallContext(ctx, &r, "admin_peers")
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c *client) NodeInfo(ctx context.Context) (*p2p.PeerInfo, error) {
	var r *p2p.PeerInfo
	err := c.c.CallContext(ctx, &r, "admin_nodeInfo")
	if err != nil {
		return nil, err
	}
	return r, err
}

// ----------------------------------------------------------------------------
// miner

func (c *client) StartMining(ctx context.Context) error {
	var r []byte
	// TODO: Result needs to be verified
	// The response data type are bytes, but we cannot parse...
	err := c.c.CallContext(ctx, &r, "miner_start", nil)
	if err != nil {
		return err
	}
	return err
}

func (c *client) StopMining(ctx context.Context) error {
	err := c.c.CallContext(ctx, nil, "miner_stop", nil)
	if err != nil {
		return err
	}
	return err
}
