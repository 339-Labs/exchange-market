package cex

/*
* 产品类型
SPOT：币币
SWAP：永续合约
MARGIN：币币杠杠
FUTURES：交割合约
OPTION：期权
*/
type InstType string

const (
	Spot    InstType = "SPOT"
	MARGIN  InstType = "MARGIN"
	SWAP    InstType = "SWAP"
	Inverse InstType = "FUTURES"
	Option  InstType = "OPTION"
)
