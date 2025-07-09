package model

type SubscribeReq struct {
	InstId  string `json:"instId"`
	Channel string `json:"channel"`
}
