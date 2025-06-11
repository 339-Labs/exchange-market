package evm

import (
	"context"
	bindings "github.com/339-Labs/exchange-market/bindings/uniswapv2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func TestIsURLAvailablecall_test(t *testing.T) {

	// 1. 初始化客户端（ethclient 实现了 bind.ContractBackend 接口）
	ethClient, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	// 1. 实例化合约对象
	contract, err := bindings.NewUniswapV2Pair(common.HexToAddress("0x3139Ffc91B99aa94DA8A2dc13f1fC36F9BDc98eE"), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 调用只读方法（View 函数）
	reserves, err := contract.GetReserves(&bind.CallOpts{
		Context: context.Background(),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("reserve0: %v", reserves.Reserve0)
}

func TestIsURLAvailablecaell_test(t *testing.T) {

	evm, err := NewEvmClient(context.Background(), "https://eth-mainnet.g.alchemy.com/v2/", "0x3139Ffc91B99aa94DA8A2dc13f1fC36F9BDc98eE")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	reserves0, _, err := evm.GetReserves()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("reserves0: %v", reserves0)

}
