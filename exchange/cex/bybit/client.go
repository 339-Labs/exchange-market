package bybit

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/339-Labs/exchange-market/config"
	"github.com/bybit-exchange/bybit.go.api"
	"log"
)

const CexName = "ByBit"

type Client struct {
	bybitClient *bybit_connector.Client
}

func NewClient(config config.CexExchangeConfig) (BybitClient, error) {

	client := bybit_connector.NewBybitHttpClient(config.ApiKey, config.ApiSecretKey)
	client.Debug = true

	return &Client{
		bybitClient: client,
	}, nil
}

type BybitClient interface {
	MarketInstrumentsInfo(request *MarketInsuranceRequest) (rsp *bybit_connector.ServerResponse, err error)
}

type MarketInsuranceRequest struct {
	Category string
	Symbol   string
	Status   string
	BaseCoin string
	Limit    int
	Cursor   string
}

func (client *Client) MarketInstrumentsInfo(request *MarketInsuranceRequest) (rsp *bybit_connector.ServerResponse, err error) {
	ctx := context.Background()
	params, err := structToMap(request)
	if err != nil {
		return nil, err
	}
	rse, err := client.bybitClient.NewUtaBybitServiceWithParams(params).GetMarketInsurance(ctx)
	if err != nil {
		log.Fatalln("")
		return nil, err
	}

	if !isSuccess(rse.RetMsg) {
		log.Println(rse.RetMsg)
		return nil, errors.New(rse.RetMsg)
	}

	return rse, nil
}

func structToMap[T any](input T) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func isSuccess(retMsg string) bool {
	if retMsg == "OK" || retMsg == "success" || retMsg == "SUCCESS" || retMsg == "" {
		return true
	}
	return false
}
