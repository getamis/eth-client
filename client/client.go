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
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// client defines typed wrappers for the Ethereum RPC API.
type client struct {
	*ethclient.Client
	rpc *ethrpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (Client, error) {
	c, err := ethrpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(rpc *ethrpc.Client) Client {
	return &client{
		Client: ethclient.NewClient(rpc),
		rpc:    rpc,
	}
}

// Close closes an existing RPC connection.
func (c *client) Close() {
	c.rpc.Close()
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

// BlockNumber returns the current block number.
func (c *client) BlockNumber(ctx context.Context) (*big.Int, error) {
	var r string
	err := c.rpc.CallContext(ctx, &r, "eth_blockNumber")
	if err != nil {
		return nil, err
	}
	h, err := hexutil.DecodeBig(r)
	return h, err
}

// ----------------------------------------------------------------------------
// miner

// StartMining starts mining operation.
func (c *client) StartMining(ctx context.Context) error {
	var r []byte
	// TODO: Result needs to be verified
	// The response data type are bytes, but we cannot parse...
	err := c.rpc.CallContext(ctx, &r, "miner_start", nil)
	if err != nil {
		return err
	}
	return err
}

// StopMining stops mining.
func (c *client) StopMining(ctx context.Context) error {
	err := c.rpc.CallContext(ctx, nil, "miner_stop", nil)
	if err != nil {
		return err
	}
	return err
}
