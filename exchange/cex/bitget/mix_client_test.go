package bitget

import (
	"github.com/339-Labs/exchange-market/config"
	"testing"
)

func SetUp() BitGetClient {

	// Do not commit configuration files to Git.  配置请勿提交git
	conf := config.Config{
		ExchangeConfig: config.ExchangeConfig{
			Bn: config.CexExchangeConfig{
				ApiKey:       config.BiGetApiKey,
				ApiUrl:       config.BiGetApiUrl,
				ApiSecretKey: config.BiGetApiSecretKey,
				WsUrl:        "",
				Passphrase:   config.BitGetPassphrase,
				TimeOut:      1000,
			},
		},
	}

	bnClient, err := NewClient(conf.ExchangeConfig.BitGet)
	if err != nil {
		panic(err)
	}

	return bnClient
}

func TestClient_AllTickerss(t *testing.T) {
	bitgetClient := SetUp()
	rsp, err := bitgetClient.AllTickers(USDT)
	if err != nil {
		t.Error(err)
	}
	t.Log(rsp)
}
