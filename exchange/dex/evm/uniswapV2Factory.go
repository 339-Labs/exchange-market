package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UniswapV2Factory interface {
	AllPairsLength() (*big.Int, error)
	GetPairs(index *big.Int) (common.Address, error)
}

func (c Client) AllPairsLength() (*big.Int, error) {
	length, err := c.UniswapV2Factory.AllPairsLength(&bind.CallOpts{
		Context: context.Background(),
	})
	return length, err
}

func (c Client) GetPairs(index *big.Int) (common.Address, error) {
	pairs, err := c.UniswapV2Factory.AllPairs(&bind.CallOpts{
		Context: context.Background(),
	}, index)
	return pairs, err
}
