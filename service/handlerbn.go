package service

import (
	"context"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/config"
	"github.com/339-Labs/exchange-market/database"
	"github.com/339-Labs/exchange-market/exchange/cex/bn"
	"github.com/339-Labs/exchange-market/redis"
	"github.com/339-Labs/exchange-market/worker"
	"github.com/ethereum/go-ethereum/log"
	"sync/atomic"
	"time"
)

type HandlerBN struct {
	BnExClient  *bn.BnExClient
	BinanceTask *worker.BinanceTask

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewHandlerBN(config *config.Config, db *database.DB, redis *redis.RedisClient, shutdown context.CancelCauseFunc) (*HandlerBN, error) {
	spotPriceMap := maps.NewPriceMap(10)
	featurePriceMap := maps.NewPriceMap(10)
	markPriceMap := maps.NewPriceMap(10)
	bnExClient, _ := bn.NewBnExClient(&config.ExchangeConfig.Bn, spotPriceMap, featurePriceMap, markPriceMap)
	bnTask, _ := worker.NewBinanceTask(shutdown, time.Second*1, spotPriceMap, featurePriceMap, markPriceMap)

	return &HandlerBN{
		BnExClient:  bnExClient,
		BinanceTask: bnTask,
		shutdown:    shutdown,
	}, nil
}

func (h *HandlerBN) Start(ctx context.Context) error {
	h.BnExClient.ExecuteWsSpot()
	h.BnExClient.ExecuteWsFeature()
	h.BinanceTask.Start()
	return nil
}

func (h *HandlerBN) Stop(ctx context.Context) error {
	h.BinanceTask.Close()
	h.BnExClient.BnWebSocketClient.Stop()
	log.Info("stop notify success")
	return nil
}

func (h *HandlerBN) Stopped() bool {
	return h.stopped.Load()
}
