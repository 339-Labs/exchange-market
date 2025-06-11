package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"log"
	"math/big"
)

type UniswapV2Pair interface {
	GetReserves() (*big.Int, *big.Int, error)
}

func (c Client) GetReserves() (*big.Int, *big.Int, error) {
	reserves, err := c.NewUniswapV2Pair.GetReserves(&bind.CallOpts{
		Context: context.Background(),
	})
	if err != nil {
		log.Fatal(err)
	}

	return reserves.Reserve0, reserves.Reserve1, nil
}
