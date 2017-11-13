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

package istanbul

import (
	"context"
	"math/big"
	"sort"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
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

// Propose injects a new authorization candidate that the validator will attempt to push through.
func (c *client) ProposeValidator(ctx context.Context, address common.Address, auth bool) error {
	var r []byte
	err := c.rpc.CallContext(ctx, &r, "istanbul_propose", address, auth)
	if err != nil {
		return ethereum.NotFound
	}
	return err
}

type addresses []common.Address

func (addrs addresses) Len() int {
	return len(addrs)
}

func (addrs addresses) Less(i, j int) bool {
	return strings.Compare(addrs[i].String(), addrs[j].String()) < 0
}

func (addrs addresses) Swap(i, j int) {
	addrs[i], addrs[j] = addrs[j], addrs[i]
}

// GetValidators retrieves the list of authorized validators at the specified block.
func (c *client) GetValidators(ctx context.Context, blockNumbers *big.Int) ([]common.Address, error) {
	var r []common.Address
	err := c.rpc.CallContext(ctx, &r, "istanbul_getValidators", toNumArg(blockNumbers))
	if err == nil && r == nil {
		return nil, ethereum.NotFound
	}

	sort.Sort(addresses(r))

	return r, err
}

func toNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}
