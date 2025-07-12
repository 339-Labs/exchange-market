package main

import (
	"github.com/339-Labs/exchange-market/exchange/cex/bitget"
	"github.com/ethereum/go-ethereum/log"
	"time"

	"github.com/339-Labs/exchange-market/config"
)

func main() {
	//log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stdout, log.LevelInfo, true)))
	// 示例1: 使用okx特定的WebSocket客户端
	bnExample()
}

// okxExample 使用okx WebSocket客户端的示例
func bnExample() {

	// 配置
	config := &config.CexExchangeConfig{
		WsUrl:        config.BiGetWsUrl,
		ApiKey:       "your-api-key",
		ApiSecretKey: "your-api-secret",
	}

	//bn.ExecuteWsSpot(config, nil)
	bitget.ExecuteWs(config, nil, nil)

	log.Info("执行")

	// 等待一段时间再关闭
	time.Sleep(120 * time.Second)

}
