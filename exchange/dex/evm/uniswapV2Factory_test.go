package evm

import (
	"context"
	"math/big"
	"testing"
)

func TestClient_AllPairsLength(t *testing.T) {

	evm, err := NewEvmClient(context.Background(), "https://eth-mainnet.g.alchemy.com/v2/", "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")

	if err != nil {
		t.Error(err)
	}
	le, err := evm.AllPairsLength()

	if err != nil {
		t.Error(err)
	}
	t.Log(le)

	p, err := evm.GetPairs(big.NewInt(10))
	if err != nil {
		t.Error(err)
	}
	t.Log(p)

}
