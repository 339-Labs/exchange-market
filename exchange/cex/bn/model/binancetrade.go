package model

type BinanceTrade struct {
	EventType     string `json:"e"` // 事件类型
	EventTime     int64  `json:"E"` // 事件时间
	Symbol        string `json:"s"` // 交易对
	TradeId       int64  `json:"t"` // 交易ID
	Price         string `json:"p"` // 价格
	Quantity      string `json:"q"` // 数量
	BuyerOrderId  int64  `json:"b"` // 买方订单ID
	SellerOrderId int64  `json:"a"` // 卖方订单ID
	TradeTime     int64  `json:"T"` // 交易时间
	IsBuyerMaker  bool   `json:"m"` // 买方是否为maker
}
