package bybit

import (
	"github.com/339-Labs/exchange-market/config"
	"testing"
)

func setUp() BybitClient {
	// Do not commit configuration files to Git.  配置请勿提交git
	conf := config.Config{
		ExchangeConfig: config.ExchangeConfig{
			Bn: config.CexExchangeConfig{
				ApiKey:       config.ByBitApiKey,
				ApiUrl:       config.ByBitApiUrl,
				ApiSecretKey: config.ByBitApiSecretKey,
				TimeOut:      1000,
			},
		},
	}

	bybitClient, err := NewClient(conf.ExchangeConfig.ByBit)
	if err != nil {
		panic(err)
	}
	return bybitClient

}

func TestClient_MarketInstrumentsInfo(t *testing.T) {
	client := setUp()

	request := MarketInsuranceRequest{
		Category: Spot,
	}

	client.MarketInstrumentsInfo(&request)
}
