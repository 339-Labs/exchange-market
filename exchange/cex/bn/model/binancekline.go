package model

// 币安行情数据结构
type BinanceKline struct {
	Symbol                   string `json:"s"` // 交易对
	OpenTime                 int64  `json:"t"` // 开盘时间
	CloseTime                int64  `json:"T"` // 收盘时间
	Symbol2                  string `json:"s"` // 交易对
	Interval                 string `json:"i"` // K线间隔
	FirstTradeID             int64  `json:"f"` // 第一笔交易ID
	LastTradeID              int64  `json:"L"` // 最后一笔交易ID
	OpenPrice                string `json:"o"` // 开盘价
	ClosePrice               string `json:"c"` // 收盘价
	HighPrice                string `json:"h"` // 最高价
	LowPrice                 string `json:"l"` // 最低价
	Volume                   string `json:"v"` // 成交量
	NumberOfTrades           int64  `json:"n"` // 成交笔数
	IsClosed                 bool   `json:"x"` // 是否为最终K线
	QuoteAssetVolume         string `json:"q"` // 成交额
	TakerBuyBaseAssetVolume  string `json:"V"` // 主动买入成交量
	TakerBuyQuoteAssetVolume string `json:"Q"` // 主动买入成交额
}

type BinanceKlineStream struct {
	EventType string       `json:"e"` // 事件类型
	EventTime int64        `json:"E"` // 事件时间
	Symbol    string       `json:"s"` // 交易对
	Kline     BinanceKline `json:"k"` // K线数据
}
