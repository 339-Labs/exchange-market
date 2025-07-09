package model

type BinanceTicker struct {
	EventType          string `json:"e"` // 事件类型
	EventTime          int64  `json:"E"` // 事件时间
	Symbol             string `json:"s"` // 交易对
	PriceChange        string `json:"p"` // 24小时价格变化
	PriceChangePercent string `json:"P"` // 24小时价格变化百分比
	WeightedAvgPrice   string `json:"w"` // 加权平均价
	PrevClosePrice     string `json:"x"` // 昨日收盘价
	LastPrice          string `json:"c"` // 最新价格
	LastQty            string `json:"Q"` // 最新成交量
	BidPrice           string `json:"b"` // 买一价
	BidQty             string `json:"B"` // 买一量
	AskPrice           string `json:"a"` // 卖一价
	AskQty             string `json:"A"` // 卖一量
	OpenPrice          string `json:"o"` // 开盘价
	HighPrice          string `json:"h"` // 最高价
	LowPrice           string `json:"l"` // 最低价
	Volume             string `json:"v"` // 成交量
	QuoteVolume        string `json:"q"` // 成交额
	OpenTime           int64  `json:"O"` // 开盘时间
	CloseTime          int64  `json:"C"` // 收盘时间
	FirstTradeId       int64  `json:"F"` // 第一笔交易ID
	LastTradeId        int64  `json:"L"` // 最后一笔交易ID
	Count              int64  `json:"n"` // 成交笔数
}
