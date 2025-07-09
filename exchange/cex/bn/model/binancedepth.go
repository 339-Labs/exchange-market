package model

type BinanceDepth struct {
	EventType     string     `json:"e"` // 事件类型
	EventTime     int64      `json:"E"` // 事件时间
	Symbol        string     `json:"s"` // 交易对
	FirstUpdateId int64      `json:"U"` // 第一个更新ID
	FinalUpdateId int64      `json:"u"` // 最终更新ID
	Bids          [][]string `json:"b"` // 买盘更新
	Asks          [][]string `json:"a"` // 卖盘更新
}
