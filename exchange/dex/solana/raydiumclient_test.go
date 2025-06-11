package solana

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"
)

// 定义结构体用于解析 JSON 响应
type Pair struct {
	PairID          string  `json:"pair_id"`   // 交易对的唯一标识符
	Name            string  `json:"name"`      // 交易对名称(如 SOL/USDC)
	LpMint          string  `json:"lp_mint"`   // 流动性提供者(LP)代币的铸币地址
	Official        bool    `json:"official"`  // 布尔值，表示是否是官方交易对
	Liquidity       float64 `json:"liquidity"` // 流动性池的总流动性
	Market          string  `json:"market"`
	Volume          float64 `json:"volume_24h"`
	VolumeQuote     float64 `json:"volume_24h_quote"`
	Fee             float64 `json:"fee_24h"`
	FeeQuote        float64 `json:"fee_24h_quote"`
	VolumeD         float64 `json:"volume_7d"`
	VolumeDQuote    float64 `json:"volume_7d_quote"`
	FeeD            float64 `json:"fee_7d"`
	FeeDQuote       float64 `json:"fee_7d_quote"`
	Price           float64 `json:"price"`             // 当前价格(可能是基础代币相对于报价代币的价格)
	LpPrice         float64 `json:"lp_price"`          // LP 代币的价格
	AmmId           string  `json:"amm_id"`            // 自动化做市商(AMM)合约的地址
	TokenAmountCoin float64 `json:"token_amount_coin"` // 基础代币(如 SOL)的数量
	TokenAmountPc   float64 `json:"token_amount_pc"`   // 报价代币(如 USDC)的数量
	TokenAmountLp   float64 `json:"token_amount_lp"`   // LP 代币的总供应量
	Apy             float64 `json:"apy"`               // 年化收益率(基于手续费收入等计算)
}

type TokenMappingLp struct {
	Id                 string `json:"id"`           // 该流动性池的唯一标识（AMM 账户地址）PairID
	BaseMint           string `json:"baseMint"`     // 交易对中的基础代币（如 USDT）
	QuoteMint          string `json:"quoteMint"`    // 交易对中的报价代币（如 USDC）
	LpMint             string `json:"lpMint"`       // LP 代币的铸币地址（用户质押流动性后获得的代币）
	ProgramId          string `json:"programId"`    // 该流动性池所属的 AMM 程序 ID
	Authority          string `json:"authority"`    // 控制该流动性池的权限账户（通常为 PDA）
	OpenOrders         string `json:"openOrders"`   // AMM 在 Serum 市场的开放订单账户
	TargetOrders       string `json:"targetOrders"` // AMM 的未完成订单账户
	BaseVault          string `json:"baseVault"`    // 存储基础代币（baseMint）的托管账户
	QuoteVault         string `json:"quoteVault"`   // 存储报价代币（quoteMint）的托管账户
	Version            int    `json:"version"`
	BaseDecimals       int    `json:"baseDecimals"`
	QuoteDecimals      int    `json:"quoteDecimals"`
	LpDecimals         int    `json:"lpDecimals"`
	WithdrawQueue      string `json:"withdrawQueue"`
	LpVault            string `json:"lpVault"` // 存储 LP 代币的托管账户（示例中未使用）
	MarketVersion      int    `json:"marketVersion"`
	MarketProgramId    string `json:"marketProgramId"`
	MarketId           string `json:"marketId"` // 关联的 Serum 市场 ID
	MarketAuthority    string `json:"marketAuthority"`
	MarketBaseVault    string `json:"marketBaseVault"`
	MarketQuoteVault   string `json:"marketQuoteVault"`
	MarketBids         string `json:"marketBids"`
	MarketAsks         string `json:"marketAsks"`
	MarketEventQueue   string `json:"marketEventQueue"`
	ModelDataAccount   string `json:"modelDataAccount"`
	LookupTableAccount string `json:"lookupTableAccount"`
}

func Test_allLp(t *testing.T) {
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

	log.Println(pairs[0])
	log.Println(pairs[1])

}

// {"id":"2EXiumdi14E9b8Fy62QcA5Uh6WdHS2b38wtSxp72Mibj","baseMint":"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB","quoteMint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","lpMint":"As3EGgLtUVpdNpE6WCKauyNRrCCwcQ57trWQ3wyRXDa6","programId":"5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h","authority":"3uaZBfHPfmpAHW7dsimC1SnyR61X4bJqQZKWmRSCXJxv","openOrders":"4zbGjjRx8bmZjynJg2KnkJ54VAk1crcrYsGMy79EXK1P","targetOrders":"AYf5abBGrwjz2n2gGP4YG91hJer22zakrizrRhddTehS","baseVault":"5XkWQL9FJL4qEvL8c3zCzzWnMGzerM3jbGuuyRprsEgG","quoteVault":"jfrmNrBtxnX1FH36ATeiaXnpA4ppQcKtv7EfrgMsgLJ","version":5,"baseDecimals":6,"quoteDecimals":6,"lpDecimals":6,"withdrawQueue":"11111111111111111111111111111111","lpVault":"11111111111111111111111111111111","marketVersion":3,"marketProgramId":"9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin","marketId":"77quYg4MGneUdjgXCunt9GgM1usmrxKY31twEy3WHwcS","marketAuthority":"FGBvMAu88q9d1Csz7ZECB5a2gbWwp6qicNxN2Mo7QhWG","marketBaseVault":"H61Y7xVnbWVXrQQx3EojTEqf3ogKVY5GfGjEn5ewyX7B","marketQuoteVault":"9FLih4qwFMjdqRAGmHeCxa64CgjP1GtcgKJgHHgz44ar","marketBids":"37m9QdvxmKRdjm3KKV2AjTiGcXMfWHQpVFnmhtb289yo","marketAsks":"AQKXXC29ybqL8DLeAVNt3ebpwMv8Sb4csberrP6Hz6o5","marketEventQueue":"9MgPMkdEHFX7DZaitSh6Crya3kCCr1As6JC75bm3mjuC","modelDataAccount":"CDSr3ssLcRB6XYPJwAfFt18MZvEZp4LjHcvzBVZ45duo","lookupTableAccount":"11111111111111111111111111111111"}

func Test_LpMappingToken(t *testing.T) {

	// 请求 Raydium 的 映射 接口
	resp, err := http.Get("https://api.raydium.io/v2/sdk/liquidity/mainnet.json")
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
	var mappings []TokenMappingLp
	err = json.Unmarshal(body, &mappings)
	if err != nil {
		panic(err)
	}

	log.Println(mappings[0])
	log.Println(mappings[1])
}
