package solana

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// 定义结构体用于解析 JSON 响应
type Pair struct {
	PairID          string  `json:"pair_id"`
	Name            string  `json:"name"`
	LpMint          string  `json:"lp_mint"`
	Official        string  `json:"official"`
	Liquidity       string  `json:"liquidity"`
	Market          string  `json:"market"`
	Volume          string  `json:"volume_24h"`
	VolumeQuote     string  `json:"volume_24h_quote"`
	Fee             string  `json:"fee_24h"`
	FeeQuote        string  `json:"fee_24h_quote"`
	VolumeD         string  `json:"volume_7d"`
	VolumeDQuote    string  `json:"volume_7d_quote"`
	FeeD            string  `json:"fee_7d"`
	FeeDQuote       string  `json:"fee_7d_quote"`
	Price           float64 `json:"price"`
	LpPrice         float64 `json:"lp_price"`
	AmmId           string  `json:"amm_id"`
	TokenAmountCoin float64 `json:"token_amount_coin"`
	TokenAmountPc   float64 `json:"token_amount_pc"`
	TokenAmountLp   float64 `json:"token_amount_lp"`
	Apy             float64 `json:"apy"`
}

func Test_raym(t *testing.T) {
	// 请求 Raydium 的 pairs 接口
	resp, err := http.Get("https://api.raydium.io/pairs")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 JSON
	var pairs []Pair
	err = json.Unmarshal(body, &pairs)
	if err != nil {
		panic(err)
	}
	// 遍历查找目标交易对，例如 SOL/USDC
	//for _, pair := range pairs {

	//}
}
