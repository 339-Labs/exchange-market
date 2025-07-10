package main

import (
	"github.com/339-Labs/exchange-market/exchange/cex/bybit"
	"log"
	"time"

	"github.com/339-Labs/exchange-market/config"
)

func main() {
	// 示例1: 使用okx特定的WebSocket客户端
	bnExample()
}

// okxExample 使用okx WebSocket客户端的示例
func bnExample() {

	// 配置
	config := &config.CexExchangeConfig{
		WsUrl:        config.ByBitWsUrl,
		ApiKey:       "your-api-key",
		ApiSecretKey: "your-api-secret",
	}

	bybit.ExecuteSpotWs(config, nil) // true表示需要登录

	log.Println("执行")

	// 等待一段时间再关闭
	time.Sleep(120 * time.Second)

}
