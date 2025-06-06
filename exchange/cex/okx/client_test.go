package okx

import (
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex"
	"testing"
)

func setUp() OkxClient {
	// Do not commit configuration files to Git.  配置请勿提交git
	conf := config.Config{
		ExchangeConfig: config.ExchangeConfig{
			Bn: config.CexExchangeConfig{
				ApiKey:       config.OkxApiKey,
				ApiUrl:       config.OkxApiUrl,
				ApiSecretKey: config.OkxApiSecretKey,
				Passphrase:   config.OkxPassphrase,
				TimeOut:      1000,
			},
		},
	}

	okx, err := NewClient(conf.ExchangeConfig.Okx)
	if err != nil {
		panic(err)
	}
	return okx

}

func TestClient_Ticker(t *testing.T) {
	okx := setUp()
	ticker, err := okx.Ticker("BTC-USD-SWAP")
	if err != nil {
		t.Error(err)
	}
	t.Log(ticker)
}

func TestClient_Tickers(t *testing.T) {
	okx := setUp()
	tickers, err := okx.Tickers(cex.Spot)
	if err != nil {
		t.Error(err)
	}
	t.Log(tickers)
}

func TestClient_Mark(t *testing.T) {
	okx := setUp()
	rsp, err := okx.MarkPrice(cex.MARGIN, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(rsp)

}
