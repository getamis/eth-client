package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/filters"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

//go:generate mockgen -source=eth.go -destination=mock_eth.go -package=rpc
type Eth interface {
	PublicTransactionPool
	PublicEthereum
	PublicBlockChain
	PublicAccount
	PublicMiner
	PublicDownloader
	PublicFilter
}

type eth struct {
	PublicTransactionPool
	PublicEthereum
	PublicBlockChain
	PublicAccount
	PublicMiner
	PublicDownloader
	PublicFilter
}

func NewEth(client *ethrpc.Client) Eth {
	return &eth{
		PublicTransactionPool: NewPublicTransactionPool(client),
		PublicEthereum:        NewPublicEthereum(client),
		PublicBlockChain:      NewPublicBlockChain(client),
		PublicAccount:         NewPublicAccount(client),
		PublicMiner:           NewPublicMiner(client),
		PublicDownloader:      NewPublicDownloader(client),
		PublicFilter:          NewPublicFilter(client),
	}
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      hexutil.Big    `json:"gas"`
	GasPrice hexutil.Big    `json:"gasPrice"`
	Value    hexutil.Big    `json:"value"`
	Data     hexutil.Bytes  `json:"data"`
	Nonce    hexutil.Uint64 `json:"nonce"`
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        common.Hash     `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              *hexutil.Big    `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex hexutil.Uint    `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

type PublicTransactionPool interface {
	// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
	GetBlockTransactionCountByNumber(ctx context.Context, blockNr string) (*hexutil.Uint, error)
	// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
	GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error)
	// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
	GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (*RPCTransaction, error)
	// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
	GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (*RPCTransaction, error)
	// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
	GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (hexutil.Bytes, error)
	// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
	GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (hexutil.Bytes, error)
	// GetTransactionCount returns the number of transactions the given address has sent for the given block number
	GetTransactionCount(ctx context.Context, address common.Address, blockNr string) (*hexutil.Uint64, error)
	// GetTransactionByHash returns the transaction for the given hash
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error)
	// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
	GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error)
	// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error)
	// SendTransaction creates a transaction for the given argument, sign it and submit it to the
	// transaction pool.
	SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error)
	// SendRawTransaction will add the signed transaction to the transaction pool.
	// The sender is responsible for signing the transaction and using the correct nonce.
	SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error)
	// Sign calculates an ECDSA signature for:
	// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
	//
	// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
	// where the V value will be 27 or 28 for legacy reasons.
	//
	// The account associated with addr must be unlocked.
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
	Sign(ctx context.Context, addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error)
	// SignTransaction will sign the given transaction with the from account.
	// The node needs to have the private key of the account corresponding with
	// the given from address and it needs to be unlocked.
	SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error)
	// PendingTransactions returns the transactions that are in the transaction pool and have a from address that is one of
	// the accounts this node manages.
	PendingTransactions(ctx context.Context) ([]*RPCTransaction, error)
	// Resend accepts an existing transaction and a new gas price and limit. It will remove
	// the given transaction from the pool and reinsert it with the new gas price and limit.
	Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice, gasLimit hexutil.Big) (common.Hash, error)
}

type publicTransactionPool struct {
	client *ethrpc.Client
}

func NewPublicTransactionPool(client *ethrpc.Client) PublicTransactionPool {
	return &publicTransactionPool{
		client: client,
	}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (pub *publicTransactionPool) GetBlockTransactionCountByNumber(ctx context.Context, blockNr string) (*hexutil.Uint, error) {
	var r *hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_getBlockTransactionCountByNumber", blockNr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (pub *publicTransactionPool) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error) {
	var r *hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_getBlockTransactionCountByHash", blockHash)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (pub *publicTransactionPool) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (*RPCTransaction, error) {
	var r *RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionByBlockNumberAndIndex", blockNr, index)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (pub *publicTransactionPool) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (*RPCTransaction, error) {
	var r *RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionByBlockNumberAndIndex", blockHash, index)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (pub *publicTransactionPool) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByBlockNumberAndIndex", blockNr, index)
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (pub *publicTransactionPool) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByBlockHashAndIndex", blockHash, index)
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (pub *publicTransactionPool) GetTransactionCount(ctx context.Context, address common.Address, blockNr string) (*hexutil.Uint64, error) {
	var r *hexutil.Uint64
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionCount", address, blockNr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionByHash returns the transaction for the given hash
func (pub *publicTransactionPool) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	var r *RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (pub *publicTransactionPool) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByHash", hash)
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (pub *publicTransactionPool) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionReceipt", hash)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (pub *publicTransactionPool) SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error) {
	var r common.Hash
	err := pub.client.CallContext(ctx, &r, "eth_sendTransaction", args)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (pub *publicTransactionPool) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	var r common.Hash
	err := pub.client.CallContext(ctx, &r, "eth_sendRawTransaction", encodedTx)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Sign calculates an ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (pub *publicTransactionPool) Sign(ctx context.Context, addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_sign", addr, data)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (pub *publicTransactionPool) Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice, gasLimit hexutil.Big) (common.Hash, error) {
	var r common.Hash
	err := pub.client.CallContext(ctx, &r, "eth_resend", sendArgs, gasPrice, gasLimit)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (pub *publicTransactionPool) SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {
	var r *SignTransactionResult
	err := pub.client.CallContext(ctx, &r, "eth_signTransaction", args)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TODO: how to implement submitTransaction?

// GetRawTransaction returns the bytes of the transaction for the given hash.
func (pub *publicTransactionPool) GetRawTransaction(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByHash", hash)
	if err != nil {
		return r, err
	}
	return r, nil
}

// PendingTransactions returns the transactions that are in the transaction pool and have a from address that is one of
// the accounts this node manages.
func (pub *publicTransactionPool) PendingTransactions(ctx context.Context) ([]*RPCTransaction, error) {
	var r []*RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_pendingTransactions")
	if err != nil {
		return r, err
	}
	return r, nil
}

type PublicEthereum interface {
	// GasPrice returns a suggestion for a gas price.
	GasPrice(ctx context.Context) (*big.Int, error)

	// ProtocolVersion returns the current Ethereum protocol version this node supports
	ProtocolVersion(ctx context.Context) (hexutil.Uint, error)

	// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
	// yet received the latest block headers from its pears. In case it is synchronizing:
	// - startingBlock: block number this node started to synchronise from
	// - currentBlock:  block number this node is currently importing
	// - highestBlock:  block number of the highest block header this node has received from peers
	// - pulledStates:  number of state entries processed until now
	// - knownStates:   number of known state entries that still need to be pulled
	Syncing(ctx context.Context) (interface{}, error)

	// Etherbase is the address that mining rewards will be send to
	Etherbase(ctx context.Context) (common.Address, error)

	// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
	Coinbase(ctx context.Context) (common.Address, error)

	// Hashrate returns the POW hashrate
	Hashrate(ctx context.Context) (hexutil.Uint64, error)
}

type publicEthereum struct {
	client *ethrpc.Client
}

func NewPublicEthereum(client *ethrpc.Client) PublicEthereum {
	return &publicEthereum{
		client: client,
	}
}

// GasPrice returns a suggestion for a gas price.
func (pub *publicEthereum) GasPrice(ctx context.Context) (*big.Int, error) {
	var r *big.Int
	err := pub.client.CallContext(ctx, &r, "eth_gasPrice")
	if err != nil {
		return r, err
	}
	return r, nil

}

// ProtocolVersion returns the current Ethereum protocol version this node supports
func (pub *publicEthereum) ProtocolVersion(ctx context.Context) (hexutil.Uint, error) {
	var r hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_protocolVersion")
	if err != nil {
		return r, err
	}
	return r, nil
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (pub *publicEthereum) Syncing(ctx context.Context) (interface{}, error) {
	var r interface{}
	err := pub.client.CallContext(ctx, &r, "eth_syncing")
	if err != nil {
		return r, err
	}
	return r, nil
}

// Etherbase is the address that mining rewards will be send to
func (pub *publicEthereum) Etherbase(ctx context.Context) (common.Address, error) {
	var r common.Address
	err := pub.client.CallContext(ctx, &r, "eth_etherbase")
	if err != nil {
		return r, err
	}
	return r, nil
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (pub *publicEthereum) Coinbase(ctx context.Context) (common.Address, error) {
	var r common.Address
	err := pub.client.CallContext(ctx, &r, "eth_coinbase")
	if err != nil {
		return r, err
	}
	return r, nil
}

// Hashrate returns the POW hashrate
func (pub *publicEthereum) Hashrate(ctx context.Context) (hexutil.Uint64, error) {
	var r hexutil.Uint64
	err := pub.client.CallContext(ctx, &r, "eth_hashrate")
	if err != nil {
		return r, err
	}
	return r, nil
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      hexutil.Big    `json:"gas"`
	GasPrice hexutil.Big    `json:"gasPrice"`
	Value    hexutil.Big    `json:"value"`
	Data     hexutil.Bytes  `json:"data"`
}

type PublicBlockChain interface {
	// BlockNumber returns the block number of the chain head.
	BlockNumber(ctx context.Context) (*big.Int, error)

	// GetBalance returns the amount of wei for the given address in the state of the
	// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
	// block numbers are also allowed.
	GetBalance(ctx context.Context, address common.Address, blockNr string) (*big.Int, error)

	// GetBlockByNumber returns the requested block. When blockNr is -1 the chain head is returned. When fullTx is true all
	// transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
	GetBlockByNumber(ctx context.Context, blockNr string, fullTx bool) (map[string]interface{}, error)

	// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
	// detail, otherwise only the transaction hash is returned.
	GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (map[string]interface{}, error)

	// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index. When fullTx is true
	// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
	GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (map[string]interface{}, error)

	// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index. When fullTx is true
	// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
	GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error)

	// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number
	GetUncleCountByBlockNumber(ctx context.Context, blockNr string) (*hexutil.Uint, error)

	// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash
	GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error)

	// GetCode returns the code stored at the given address in the state for the given block number.
	GetCode(ctx context.Context, address common.Address, blockNr string) (hexutil.Bytes, error)

	// GetStorageAt returns the storage from the state at the given address, key and
	// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
	// numbers are also allowed.
	GetStorageAt(ctx context.Context, address common.Address, key string, blockNr string) (hexutil.Bytes, error)

	// Call executes the given transaction on the state for the given block number.
	// It doesn't make and changes in the state/blockchain and is useful to execute and retrieve values.
	Call(ctx context.Context, args CallArgs, blockNr string) (hexutil.Bytes, error)

	// EstimateGas returns an estimate of the amount of gas needed to execute the
	// given transaction against the current pending block.
	EstimateGas(ctx context.Context, args CallArgs) (*hexutil.Big, error)
}

type publicBlockChain struct {
	client *ethrpc.Client
}

func NewPublicBlockChain(client *ethrpc.Client) PublicBlockChain {
	return &publicBlockChain{
		client: client,
	}
}

// BlockNumber returns the block number of the chain head.
func (pub *publicBlockChain) BlockNumber(ctx context.Context) (*big.Int, error) {
	var r *big.Int
	err := pub.client.CallContext(ctx, &r, "eth_blockNumber")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (pub *publicBlockChain) GetBalance(ctx context.Context, address common.Address, blockNr string) (*big.Int, error) {
	var r *big.Int
	err := pub.client.CallContext(ctx, &r, "eth_getBalance")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetBlockByNumber returns the requested block. When blockNr is -1 the chain head is returned. When fullTx is true all
// transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (pub *publicBlockChain) GetBlockByNumber(ctx context.Context, blockNr string, fullTx bool) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := pub.client.CallContext(ctx, &r, "eth_getBlockByNumber")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (pub *publicBlockChain) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := pub.client.CallContext(ctx, &r, "eth_getBlockByHash")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (pub *publicBlockChain) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := pub.client.CallContext(ctx, &r, "eth_getUncleByBlockNumberAndIndex")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (pub *publicBlockChain) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := pub.client.CallContext(ctx, &r, "eth_getUncleByBlockHashAndIndex")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number
func (pub *publicBlockChain) GetUncleCountByBlockNumber(ctx context.Context, blockNr string) (*hexutil.Uint, error) {
	var r *hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_getUncleCountByBlockNumber")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash
func (pub *publicBlockChain) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error) {
	var r *hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_getUncleCountByBlockHash")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetCode returns the code stored at the given address in the state for the given block number.
func (pub *publicBlockChain) GetCode(ctx context.Context, address common.Address, blockNr string) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getCode")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (pub *publicBlockChain) GetStorageAt(ctx context.Context, address common.Address, key string, blockNr string) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getStorageAt")
	if err != nil {
		return r, err
	}
	return r, nil
}

// Call executes the given transaction on the state for the given block number.
// It doesn't make and changes in the state/blockchain and is useful to execute and retrieve values.
func (pub *publicBlockChain) Call(ctx context.Context, args CallArgs, blockNr string) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_call")
	if err != nil {
		return r, err
	}
	return r, nil
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (pub *publicBlockChain) EstimateGas(ctx context.Context, args CallArgs) (*hexutil.Big, error) {
	var r *hexutil.Big
	err := pub.client.CallContext(ctx, &r, "eth_estimateGas")
	if err != nil {
		return r, err
	}
	return r, nil
}

type PublicAccount interface {
	// Accounts returns the collection of accounts this node manages
	Accounts(context.Context) ([]common.Address, error)
}

type publicAccount struct {
	client *ethrpc.Client
}

func NewPublicAccount(client *ethrpc.Client) PublicAccount {
	return &publicAccount{
		client: client,
	}
}

// Accounts returns the collection of accounts this node manages
func (pub *publicAccount) Accounts(ctx context.Context) ([]common.Address, error) {
	var r []common.Address
	err := pub.client.CallContext(ctx, &r, "eth_accounts")
	if err != nil {
		return r, err
	}
	return r, nil
}

type PublicMiner interface {
	// Mining returns an indication if this node is currently mining.
	Mining(ctx context.Context) (bool, error)

	// SubmitWork can be used by external miner to submit their POW solution. It returns an indication if the work was
	// accepted. Note, this is not an indication if the provided work was valid!
	SubmitWork(ctx context.Context, nonce types.BlockNonce, solution, digest common.Hash) (bool, error)

	// GetWork returns a work package for external miner. The work package consists of 3 strings
	// result[0], 32 bytes hex encoded current block header pow-hash
	// result[1], 32 bytes hex encoded seed hash used for DAG
	// result[2], 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
	GetWork(ctx context.Context) ([3]string, error)

	// SubmitHashrate can be used for remote miners to submit their hash rate. This enables the node to report the combined
	// hash rate of all miners which submit work through this node. It accepts the miner hash rate and an identifier which
	// must be unique between nodes.
	SubmitHashrate(ctx context.Context, hashrate hexutil.Uint64, id common.Hash) (bool, error)
}

type publicMiner struct {
	client *ethrpc.Client
}

func NewPublicMiner(client *ethrpc.Client) PublicMiner {
	return &publicMiner{
		client: client,
	}
}

// Mining returns an indication if this node is currently mining.
func (pub *publicMiner) Mining(ctx context.Context) (bool, error) {
	var r bool
	err := pub.client.CallContext(ctx, &r, "eth_mining")
	if err != nil {
		return r, err
	}
	return r, nil
}

// SubmitWork can be used by external miner to submit their POW solution. It returns an indication if the work was
// accepted. Note, this is not an indication if the provided work was valid!
func (pub *publicMiner) SubmitWork(ctx context.Context, nonce types.BlockNonce, solution, digest common.Hash) (bool, error) {
	var r bool
	err := pub.client.CallContext(ctx, &r, "eth_submitWork")
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetWork returns a work package for external miner. The work package consists of 3 strings
// result[0], 32 bytes hex encoded current block header pow-hash
// result[1], 32 bytes hex encoded seed hash used for DAG
// result[2], 32 bytes hex encoded boundary condition ("target"), 2^256/difficulty
func (pub *publicMiner) GetWork(ctx context.Context) ([3]string, error) {
	var r [3]string
	err := pub.client.CallContext(ctx, &r, "eth_getWork")
	if err != nil {
		return r, err
	}
	return r, nil
}

// SubmitHashrate can be used for remote miners to submit their hash rate. This enables the node to report the combined
// hash rate of all miners which submit work through this node. It accepts the miner hash rate and an identifier which
// must be unique between nodes.
func (pub *publicMiner) SubmitHashrate(ctx context.Context, hashrate hexutil.Uint64, id common.Hash) (bool, error) {
	var r bool
	err := pub.client.CallContext(ctx, &r, "eth_submitHashrate")
	if err != nil {
		return r, err
	}
	return r, nil
}

type PublicDownloader interface {
	// Duplicate function with PublicEthereum
	// Syncing provides information when this nodes starts synchronising with the Ethereum network and when it's finished.
	// Syncing(ctx context.Context) (*ethrpc.Subscription, error)

	// SubscribeSyncStatus creates a subscription that will broadcast new synchronisation updates.
	// The given channel must receive interface values, the result can either
	SubscribeSyncStatus(ctx context.Context, status chan interface{}) (*downloader.SyncStatusSubscription, error)
}

type publicDownloader struct {
	client *ethrpc.Client
}

func NewPublicDownloader(client *ethrpc.Client) PublicDownloader {
	return &publicDownloader{
		client: client,
	}
}

// // Syncing provides information when this nodes starts synchronising with the Ethereum network and when it's finished.
// func (pub *publicDownloader) Syncing(ctx context.Context) (*ethrpc.Subscription, error) {
// 	var r *ethrpc.Subscription
// 	err := pub.client.CallContext(ctx, &r, "eth_syncing")
// 	if err != nil {
// 		return r, err
// 	}
// 	return r, nil
// }

// SubscribeSyncStatus creates a subscription that will broadcast new synchronisation updates.
// The given channel must receive interface values, the result can either
func (pub *publicDownloader) SubscribeSyncStatus(ctx context.Context, status chan interface{}) (*downloader.SyncStatusSubscription, error) {
	var r *downloader.SyncStatusSubscription
	err := pub.client.CallContext(ctx, &r, "eth_subscribeSyncStatus")
	if err != nil {
		return r, err
	}
	return r, nil
}

type PublicFilter interface {
	// NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
	// as transactions enter the pending state.
	//
	// It is part of the filter package because this filter can be used throug the
	// `eth_getFilterChanges` polling method that is also used for log filters.
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newpendingtransactionfilter
	NewPendingTransactionFilter(ctx context.Context) (ethrpc.ID, error)

	// NewPendingTransactions creates a subscription that is triggered each time a transaction
	// enters the transaction pool and was signed from one of the transactions this nodes manages.
	NewPendingTransactions(ctx context.Context) (*ethrpc.Subscription, error)

	// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
	// It is part of the filter package since polling goes with eth_getFilterChanges.
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter
	NewBlockFilter(ctx context.Context) (ethrpc.ID, error)

	// NewHeads send a notification each time a new (header) block is appended to the chain.
	NewHeads(ctx context.Context) (*ethrpc.Subscription, error)

	// Logs creates a subscription that fires for all new log that match the given filter criteria.
	Logs(ctx context.Context, crit filters.FilterCriteria) (*ethrpc.Subscription, error)
}

type publicFilter struct {
	client *ethrpc.Client
}

func NewPublicFilter(client *ethrpc.Client) PublicFilter {
	return &publicFilter{
		client: client,
	}
}

// NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
// as transactions enter the pending state.
//
// It is part of the filter package because this filter can be used throug the
// `eth_getFilterChanges` polling method that is also used for log filters.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newpendingtransactionfilter
func (pub *publicFilter) NewPendingTransactionFilter(ctx context.Context) (ethrpc.ID, error) {
	var r ethrpc.ID
	err := pub.client.CallContext(ctx, &r, "eth_newPendingTransactionFilter")
	if err != nil {
		return r, err
	}
	return r, nil
}

// NewPendingTransactions creates a subscription that is triggered each time a transaction
// enters the transaction pool and was signed from one of the transactions this nodes manages.
func (pub *publicFilter) NewPendingTransactions(ctx context.Context) (*ethrpc.Subscription, error) {
	var r *ethrpc.Subscription
	err := pub.client.CallContext(ctx, &r, "eth_newPendingTransactions")
	if err != nil {
		return r, err
	}
	return r, nil
}

// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
// It is part of the filter package since polling goes with eth_getFilterChanges.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter
func (pub *publicFilter) NewBlockFilter(ctx context.Context) (ethrpc.ID, error) {
	var r ethrpc.ID
	err := pub.client.CallContext(ctx, &r, "eth_newBlockFilter")
	if err != nil {
		return r, err
	}
	return r, nil
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (pub *publicFilter) NewHeads(ctx context.Context) (*ethrpc.Subscription, error) {
	var r *ethrpc.Subscription
	err := pub.client.CallContext(ctx, &r, "eth_newHeads")
	if err != nil {
		return r, err
	}
	return r, nil
}

// Logs creates a subscription that fires for all new log that match the given filter criteria.
func (pub *publicFilter) Logs(ctx context.Context, crit filters.FilterCriteria) (*ethrpc.Subscription, error) {
	var r *ethrpc.Subscription
	err := pub.client.CallContext(ctx, &r, "eth_logs")
	if err != nil {
		return r, err
	}
	return r, nil
}
