package bn

import (
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex"
	"testing"
)

func setUp() BnClient {

	// Do not commit configuration files to Git.  配置请勿提交git
	conf := config.Config{
		ExchangeConfig: config.ExchangeConfig{
			Bn: config.CexExchangeConfig{
				ApiKey:       config.BnApiKey,
				ApiUrl:       config.BnApiUrl,
				ApiSecretKey: config.BnApiSecretKey,
				TimeOut:      1000,
			},
		},
	}

	bnClient, err := NewClient(conf.ExchangeConfig.Bn)
	if err != nil {
		panic(err)
	}

	return bnClient
}

func TestClient_GetAllCoinsInfoService(t *testing.T) {
	bnClient := setUp()
	rsp, err := bnClient.Tickers(cex.Spot, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(rsp)
}
