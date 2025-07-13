package bn

import (
	"github.com/339-Labs/exchange-market/api/service"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
)

type Client struct {
	config config.Config
	db     database.DB
}

func NewClient(config config.Config, db database.DB) service.HandlerSymbolAdaptor {
	return &Client{
		config: config,
		db:     db,
	}
}
