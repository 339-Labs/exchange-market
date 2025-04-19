package bitget

import (
	"github.com/339-Labs/exchange-market/config"
	config2 "github.com/339-Labs/v3-bitget-api-sdk-go/config"
	"github.com/339-Labs/v3-bitget-api-sdk-go/pkg/client"
	v2 "github.com/339-Labs/v3-bitget-api-sdk-go/pkg/client/v2"
)

const CexName = "BitGet"

type Client struct {
	bitgetApiClient     *client.BitgetApiClient
	v2MixAccountClient  *v2.MixAccountClient
	v2MixMarketClient   *v2.MixMarketClient
	v2MixOrderClient    *v2.MixOrderClient
	v2SpotAccountClient *v2.SpotAccountClient
	v2SpotMarketClient  *v2.SpotMarketClient
	v2SpotOrderClient   *v2.SpotOrderClient
	v2SpotWalletApi     *v2.SpotWalletApi
}

func NewClient(config config.CexExchangeConfig) (BitGetClient, error) {

	bitgetConfig := config2.NewBitgetConfig(config.ApiKey, config.ApiSecretKey, config.Passphrase, int(config.TimeOut), "")

	client := new(client.BitgetApiClient).Init(bitgetConfig)
	mixAccountClient := new(v2.MixAccountClient).Init(bitgetConfig)
	mixMarketClient := new(v2.MixMarketClient).Init(bitgetConfig)
	mixOrderClient := new(v2.MixOrderClient).Init(bitgetConfig)
	spotAccountClient := new(v2.SpotAccountClient).Init(bitgetConfig)
	spotMarketClient := new(v2.SpotMarketClient).Init(bitgetConfig)
	spotOrderClient := new(v2.SpotOrderClient).Init(bitgetConfig)
	spotWalletApi := new(v2.SpotWalletApi).Init(bitgetConfig)

	return &Client{
		bitgetApiClient:     client,
		v2MixAccountClient:  mixAccountClient,
		v2MixMarketClient:   mixMarketClient,
		v2MixOrderClient:    mixOrderClient,
		v2SpotAccountClient: spotAccountClient,
		v2SpotMarketClient:  spotMarketClient,
		v2SpotOrderClient:   spotOrderClient,
		v2SpotWalletApi:     spotWalletApi,
	}, nil
}

type BitGetClient interface {
	/**                              现货                                 */

	// 获取币种信息  如不填写，默认返回全部币种信息
	CoinsInfo(coin string) (string, error)
	// 如不填写，默认返回全部交易对信息
	SpotSymbols(symbol string) (string, error)
	//获取行情信息 如不填写，默认返回全部交易对信息
	SpotTickers(symbol string) (string, error)

	/**                              合约                                 */

	// 获取全部交易对行情
	AllTickers(productType ProductType) (string, error)
}
