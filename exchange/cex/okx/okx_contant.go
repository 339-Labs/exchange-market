package okx

/*
* 产品类型
SPOT：币币
SWAP：永续合约
FUTURES：交割合约
OPTION：期权
*/
type InstType string

const (
	Spot    InstType = "SPOT"
	Linear  InstType = "SWAP"
	Inverse InstType = "FUTURES"
	Option  InstType = "OPTION"
)
