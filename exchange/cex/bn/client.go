package bn

import (
	"context"
	"github.com/339-Labs/exchange-market/config"
	"github.com/binance/binance-connector-go"
	"log"
)

const CexName = "BN"

type Client struct {
	bnClient *binance_connector.Client
}

func NewClient(config config.CexExchangeConfig) (BnClient, error) {
	client := binance_connector.NewClient(config.ApiKey, config.ApiSecretKey)
	client.Debug = true
	client.TimeOffset = -config.TimeOut
	return &Client{
		bnClient: client,
	}, nil
}

type BnClient interface {
	GetAllCoinsInfoService() ([]*binance_connector.CoinInfo, error)
}

func (client *Client) GetAllCoinsInfoService() ([]*binance_connector.CoinInfo, error) {
	ctx := context.Background()
	rsp, err := client.bnClient.NewGetAllCoinsInfoService().Do(ctx)
	if err != nil {
		log.Fatalln("NewGetAllCoinsInfoService fail : ", err)
		return nil, err
	}
	//binance_connector.PrettyPrint(rsp)
	return rsp, nil
}
