package model

type WsLoginBaseReq struct {
	Op   string   `json:"op"`
	Args []string `json:"args"` // apiKey,timestamp,singature
}

type WsLoginReq struct {
	ApiKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}
