package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/339-Labs/exchange-market/common/maps"
	"github.com/339-Labs/exchange-market/common/tasks"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

type ByBitTask struct {
	spotPriceMap    *maps.PriceMap
	featurePriceMap *maps.PriceMap
	resourceCtx     context.Context

	resourceCancel context.CancelFunc
	tasks          tasks.Group
	ticker         *time.Ticker
}

func NewByBitTask(cfg *context.Context, shutdown context.CancelCauseFunc, duration time.Duration, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap) (*ByBitTask, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &ByBitTask{
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("bybit ws error "))
		}},
		ticker:          time.NewTicker(duration),
		spotPriceMap:    spotPriceMap,
		featurePriceMap: featurePriceMap,
	}, nil
}

func (t *ByBitTask) Start() error {
	log.Info("bybit task started")
	t.tasks.Go(func() error {
		for {

			select {

			case <-t.ticker.C:
				// todo  bybit ws data handler, spot and feature

			case <-t.resourceCtx.Done():
				log.Info("stop bybit task in work")
				return nil

			}

		}
	})
	return nil
}

func (t *ByBitTask) Close() error {
	var result error
	t.resourceCancel()
	t.ticker.Stop()
	log.Info("bybit task stopped")
	if err := t.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("bybit ws task wait error"))
	}
	log.Info("bybit task stopped success")
	return nil
}
