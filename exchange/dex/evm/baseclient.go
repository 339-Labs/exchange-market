package evm

import (
	"context"
	"fmt"
	uniswapv2 "github.com/339-Labs/exchange-market/bindings/uniswapv2"
	uniswapv3 "github.com/339-Labs/exchange-market/bindings/uniswapv3"
	"github.com/339-Labs/exchange-market/common/retry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"net"
	"net/url"
	"time"
)

const (
	defaultDialTimeout    = 5 * time.Second
	defaultDialAttempts   = 3
	defaultRequestTimeout = 100 * time.Second
)

type Client struct {
	EthClient        *ethclient.Client
	NewUniswapV2Pair *uniswapv2.UniswapV2Pair
	UniswapV2Factory *uniswapv2.UniswapV2Factory
	UniswapV3Pool    *uniswapv3.UniswapV3Pool
	UniswapV3Factory *uniswapv3.UniswapV3Factory
}

type EvmClient interface {
	UniswapV2Pair
	UniswapV2Factory
}

func NewEvmClient(ctx context.Context, rpcUrl string, contractAddress string) (EvmClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()
	off := retry.Exponential()
	ethClient, err := retry.Do(ctx, defaultDialAttempts, off, func() (*ethclient.Client, error) {
		if !IsURLAvailable(rpcUrl) {
			return nil, fmt.Errorf("rpc url %s is not available", rpcUrl)
		}
		client, err := ethclient.DialContext(ctx, rpcUrl)
		if err != nil {
			return nil, err
		}
		return client, nil
	})
	if err != nil {
		return nil, err
	}

	v2corecontract, err := uniswapv2.NewUniswapV2Pair(common.HexToAddress(contractAddress), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	v2factorycontract, err := uniswapv2.NewUniswapV2Factory(common.HexToAddress(contractAddress), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	v3corecontract, err := uniswapv3.NewUniswapV3Pool(common.HexToAddress(contractAddress), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	v3factorycontract, err := uniswapv3.NewUniswapV3Factory(common.HexToAddress(contractAddress), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	return Client{
		EthClient:        ethClient,
		NewUniswapV2Pair: v2corecontract,
		UniswapV2Factory: v2factorycontract,
		UniswapV3Pool:    v3corecontract,
		UniswapV3Factory: v3factorycontract,
	}, nil
}

func IsURLAvailable(address string) bool {
	u, err := url.Parse(address)
	if err != nil {
		return false
	}
	addr := u.Host
	if u.Port() == "" {
		switch u.Scheme {
		case "http", "ws":
			addr += ":80"
		case "https", "wss":
			addr += ":443"
		default:
			return true
		}
	}
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}
