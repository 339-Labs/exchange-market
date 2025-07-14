package bybit

import (
	"context"
	"errors"
	"github.com/339-Labs/exchange-market/common"
	symbol2 "github.com/339-Labs/exchange-market/database/symbol"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

func (c *Client) InitSpotSymbol() error {

	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/v5/market/instruments-info?category=spot", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		rspMap := common.BytesToMap(resp.Body)
		retMsg, ok1 := rspMap["retMsg"].(string)
		retCode, ok2 := rspMap["code"].(string)

		if ok1 && ok2 && retCode == "0" && retMsg == "success" {
			result, _ := rspMap["result"].(map[string]interface{})
			if list, ok := result["list"].([]map[string]interface{}); ok && len(list) > 0 {
				symbols := make([]symbol2.MarketSymbol, 0, len(list))
				for _, vv := range list {

					symbol, _ := vv["symbol"].(string)
					baseCoin, _ := vv["baseCoin"].(string)
					quoteCoin, _ := vv["quoteCoin"].(string)

					var marketSymbol = symbol2.MarketSymbol{
						Symbol:        symbol,
						UnifiedSymbol: baseCoin + "/" + quoteCoin,
						InstType:      "Spot",
						Exchange:      "ByBit",
						ChainId:       "999999",
						Base:          baseCoin,
						Quote:         quoteCoin,
						Timestamp:     uint64(time.Now().UnixMilli()),
					}
					symbols = append(symbols, marketSymbol)

				}
				log.Info("====", len(symbols))
			}
		}

	}
	return errors.New("bybit: init spot symbol error")
}

func (c *Client) InitFeatureSymbol() error {
	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/fapi/v1/ticker/price", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		arrMap := common.BytesToArrMap(resp.Body)

		symbols := make([]symbol2.MarketSymbol, len(arrMap), len(arrMap))
		for _, v := range arrMap {
			symbol, _ := v["symbol"].(string)
			baseCoin, _ := v["baseCoin"].(string)
			quoteCoin, _ := v["quoteCoin"].(string)
			contractType, _ := v[" contractType"].(string)
			status, _ := v["status"].(string)

			if "LinearPerpetual" == contractType && "Trading" == status {
				var marketSymbol = symbol2.MarketSymbol{
					Symbol:        symbol,
					UnifiedSymbol: baseCoin + "/" + quoteCoin,
					InstType:      "Feature",
					Exchange:      "Bn",
					ChainId:       "999999",
					Base:          baseCoin,
					Quote:         quoteCoin,
					Timestamp:     uint64(time.Now().UnixMilli()),
				}
				symbols = append(symbols, marketSymbol)
			}

		}
		log.Info("====", len(symbols))
	}
	return errors.New("bybit: init feature symbol error")
}
