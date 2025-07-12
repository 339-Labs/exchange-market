package bn

import (
	"encoding/json"
	"fmt"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/exchange/cex/bn/constants"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type BnExClient struct {
	BnWebSocketClient *BnWebSocketClient
	config            *config.CexExchangeConfig
	spotPriceMap      *maps.PriceMap
	featurePriceMap   *maps.PriceMap
	markPriceMap      *maps.PriceMap
}

func NewBnExClient(config *config.CexExchangeConfig, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap, markPriceMap *maps.PriceMap) (*BnExClient, error) {
	// 创建WebSocket客户端
	client := NewBnWebSocketClient(config, false)

	return &BnExClient{
		BnWebSocketClient: client,
		config:            config,
		spotPriceMap:      spotPriceMap,
		featurePriceMap:   featurePriceMap,
		markPriceMap:      markPriceMap,
	}, nil
}

func (bn *BnExClient) ExecuteWsSpot() {

	// 设置全局监听器
	bn.BnWebSocketClient.SetListeners(
		func(message string) {
			fmt.Printf("Global message: %s", message)
		},
		func(error string) {
			log.Error("Error: %s", error)
		},
	)

	// 连接
	if err := bn.BnWebSocketClient.Start(); err != nil {
		log.Error("Failed to connect:", err)
	}

	// 等待连接完成
	time.Sleep(2 * time.Second)

	// 订阅所有交易对精简 24小时价格变动统计
	err := bn.BnWebSocketClient.SubscribeMiniTickerAll(func(message string) {
		bn.handlerDataType(message, constants.Spot)
	})
	if err != nil {
		log.Error("Failed to subscribe ticker: %v", err)
	}

}

func (bn *BnExClient) ExecuteWsFeature() {

	// 创建WebSocket客户端
	cf := &config.CexExchangeConfig{
		WsUrl:        bn.config.WsUrlFeature,
		ApiKey:       bn.config.ApiKey,
		ApiSecretKey: bn.config.ApiSecretKey,
	}
	client := NewBnWebSocketClient(cf, false)

	// 设置全局监听器
	client.SetListeners(
		func(message string) {
			fmt.Printf("Global message: %s", message)
		},
		func(error string) {
			log.Error("Error: %s", error)
		},
	)

	// 连接
	if err := client.Start(); err != nil {
		log.Error("Failed to connect:", err)
	}

	// 等待连接完成
	time.Sleep(2 * time.Second)

	// 全部交易对精简信息 订阅24小时价格变动统计
	err := client.SubscribeAllTicker(func(message string) {
		bn.handlerDataType(message, constants.Feature)
	})
	if err != nil {
		log.Info("Failed to subscribe ticker: %v", err)
	}

	// 全部交易对标记价格 订阅24小b价格变动统计
	err = client.SubscribeMarkPriceAll(func(message string) {
		bn.handlerDataType(message, constants.Feature)
	})
	if err != nil {
		log.Info("Failed to subscribe ticker: %v", err)
	}

}

func (bn *BnExClient) handlerDataType(message string, t string) {
	var v interface{}
	err := json.Unmarshal([]byte(message), &v)
	if err != nil {
		panic(err)
	}
	switch vv := v.(type) {
	case map[string]interface{}:
		bn.handlerData(vv, t)
	case []interface{}:
		for _, data := range vv {
			if item, ok := data.(map[string]interface{}); ok {
				bn.handlerData(item, t)
			}
		}
	default:
		fmt.Println("未知类型 s%", vv)
	}
}

func (bn *BnExClient) handlerData(data map[string]interface{}, t string) {

	e, _ := data["e"]

	if t == constants.Spot && constants.EventTicker == e.(string) {
		bn.handlerSpot(data)
	}
	if t == constants.Spot && constants.EventMiniTicker == e.(string) {
		bn.handlerSpot(data)
	}
	if t == constants.Feature && constants.EventTicker == e.(string) {
		bn.handlerFeature(data)
	}
	if t == constants.Feature && constants.EventMiniTicker == e.(string) {
		bn.handlerFeature(data)
	}
	if t == constants.Feature && constants.EventMarkPrice == e.(string) {
		bn.handlerFeatureMark(data)
	}
}

func (bn *BnExClient) handlerSpot(spot map[string]interface{}) {
	log.Info("spot ------ ,instId: %s , lastPr: %s", spot["s"], spot["c"])

	bn.spotPriceMap.Write(spot["s"].(string), &maps.PriceData{
		Symbol:    spot["s"].(string),
		Price:     spot["c"].(string),
		Timestamp: spot["E"].(string),
	})
}

func (bn *BnExClient) handlerFeature(feature map[string]interface{}) {
	log.Info("feature ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["s"], feature["c"])

	bn.featurePriceMap.Write(feature["s"].(string), &maps.PriceData{
		Symbol:    feature["s"].(string),
		Price:     feature["c"].(string),
		Timestamp: feature["E"].(string),
	})

}

func (bn *BnExClient) handlerFeatureMark(feature map[string]interface{}) {
	log.Info("mark feature  ------ ,instId: %s , lastPr: %s , fundingRate: %s", feature["s"], feature["p"], feature["r"])

	bn.markPriceMap.Write(feature["s"].(string), &maps.PriceData{
		Symbol:      feature["s"].(string),
		FundingRate: feature["r"].(string),
		MarkPrice:   feature["p"].(string),
		Timestamp:   feature["E"].(string),
	})
}
