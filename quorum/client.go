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

package quorum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	ethClient "github.com/getamis/eth-client/client"
)

// client defines typed wrappers for the eth-client.
type client struct {
	ethClient.Client
	rpc *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (Client, error) {
	rc, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	c := &client{
		Client: ethClient.NewClient(rc),
		rpc:    rc,
	}

	return c, nil
}

// CreateContract creates a contract with the given parameters.
func (c *client) CreateContract(ctx context.Context, from common.Address, bytecode string, gas *big.Int) (txHash string, err error) {
	var hex hexutil.Bytes
	arg := map[string]interface{}{
		"from": from,
		"gas":  (*hexutil.Big)(gas),
		"data": bytecode,
	}
	if err = c.rpc.CallContext(ctx, &hex, "eth_sendTransaction", arg); err == nil {
		txHash = hex.String()
	}
	return
}

// CreatePrivateContract creates a private contract with the given parameters. The related information can refer
// to https://github.com/jpmorganchase/quorum/wiki/Using-Quorum#creating-private-transactionscontracts.
func (c *client) CreatePrivateContract(ctx context.Context, from common.Address, bytecode string, gas *big.Int, privateFor []string) (txHash string, err error) {
	var hex hexutil.Bytes
	arg := map[string]interface{}{
		"from":       from,
		"gas":        (*hexutil.Big)(gas),
		"data":       bytecode,
		"privateFor": privateFor,
	}
	if err = c.rpc.CallContext(ctx, &hex, "eth_sendTransaction", arg); err == nil {
		txHash = hex.String()
	}
	return
}
