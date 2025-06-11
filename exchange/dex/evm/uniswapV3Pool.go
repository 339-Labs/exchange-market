package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UniswapV3Pool interface {
	Token0() (common.Address, error)
	Token1() (common.Address, error)
	GetFee() (*big.Int, error)
	Slot0() (*big.Float, *big.Int, error)
}

func (c Client) Token0() (common.Address, error) {
	token0, err := c.UniswapV3Pool.Token0(&bind.CallOpts{
		Context: context.Background(),
	})
	return token0, err
}
func (c Client) Token1() (common.Address, error) {
	token1, err := c.UniswapV3Pool.Token1(&bind.CallOpts{
		Context: context.Background(),
	})
	return token1, err
}

func (c Client) GetFee() (*big.Int, error) {
	fee, err := c.UniswapV3Pool.Fee(&bind.CallOpts{
		Context: context.Background(),
	})
	return fee, err
}

var Q96 = new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

func (c Client) Slot0() (*big.Float, *big.Int, error) {
	rsp, err := c.UniswapV3Pool.Slot0(&bind.CallOpts{
		Context: context.Background(),
	})
	if err != nil {
		return nil, nil, err
	}
	// 转为 big.Float
	sqrtPrice := new(big.Float).SetInt(rsp.SqrtPriceX96)
	// 除以 2^96
	ratio := new(big.Float).Quo(sqrtPrice, Q96)
	// 平方得到 token1/token0 的价格
	price := new(big.Float).Mul(ratio, ratio)

	return price, rsp.Tick, nil
}
