package okx

import (
	"context"
	"errors"
	"github.com/339-Labs/exchange-market/common"
	symbol2 "github.com/339-Labs/exchange-market/database/symbol"
	"github.com/ethereum/go-ethereum/log"
	"strings"
	"time"
)

func (c *Client) InitSpotSymbol() error {

	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/api/v5/public/instruments?instType=SPOT", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		rspMap := common.BytesToMap(resp.Body)
		if retCode, ok := rspMap["code"].(string); ok && retCode == "0" {

			if list, ok1 := rspMap["data"].([]map[string]interface{}); ok1 {
				symbols := make([]symbol2.MarketSymbol, 0, len(list))
				for _, vv := range list {

					symbol, _ := vv["instId"].(string)
					baseCoin, _ := vv["baseCcy"].(string)
					quoteCoin, _ := vv["quoteCcy"].(string)

					var marketSymbol = symbol2.MarketSymbol{
						Symbol:        symbol,
						UnifiedSymbol: baseCoin + "/" + quoteCoin,
						InstType:      "Spot",
						Exchange:      "Okx",
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
	return errors.New("okx: init spot symbol error")
}

func (c *Client) InitFeatureSymbol() error {
	header := make(map[string]string, 0)
	resp, err := c.resty.GET(context.Background(), "/api/v5/public/instruments?instType=SWAP", header)
	if err != nil {
		return err
	}

	if resp.IsSuccess() && resp.StatusCode == 200 {
		rspMap := common.BytesToMap(resp.Body)
		if retCode, ok := rspMap["code"].(string); ok && retCode == "0" {

			if list, ok1 := rspMap["data"].([]map[string]interface{}); ok1 {
				symbols := make([]symbol2.MarketSymbol, 0, len(list))
				for _, vv := range list {

					symbol, _ := vv["instId"].(string)

					instFamily, _ := vv["instFamily"].(string)
					if base, quote, ok := splitSymbol(instFamily); ok {

						var marketSymbol = symbol2.MarketSymbol{
							Symbol:        symbol,
							UnifiedSymbol: base + "/" + quote,
							InstType:      "Feature",
							Exchange:      "Okx",
							ChainId:       "999999",
							Base:          base,
							Quote:         quote,
							Timestamp:     uint64(time.Now().UnixMilli()),
						}
						symbols = append(symbols, marketSymbol)

					}
				}
				log.Info("====", len(symbols))

			}

		}
	}
	return errors.New("okx: init feature symbol error")
}

func splitSymbol(symbol string) (base, quote string, ok bool) {
	parts := strings.Split(symbol, "-")
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
