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

type BitGetTask struct {
	spotPriceMap    *maps.PriceMap
	featurePriceMap *maps.PriceMap
	resourceCtx     context.Context

	resourceCancel context.CancelFunc
	tasks          tasks.Group
	ticker         *time.Ticker
}

func NewBitGetTask(cfg *context.Context, shutdown context.CancelCauseFunc, duration time.Duration, spotPriceMap *maps.PriceMap, featurePriceMap *maps.PriceMap) (*BitGetTask, error) {
	resCtx, resCancel := context.WithCancel(context.Background())
	return &BitGetTask{
		resourceCtx:    resCtx,
		resourceCancel: resCancel,
		tasks: tasks.Group{HandleCrit: func(err error) {
			shutdown(fmt.Errorf("bitget ws error "))
		}},
		ticker:          time.NewTicker(duration),
		spotPriceMap:    spotPriceMap,
		featurePriceMap: featurePriceMap,
	}, nil
}

func (t *BitGetTask) Start() error {
	log.Info("bitget task started")
	t.tasks.Go(func() error {
		for {

			select {

			case <-t.ticker.C:
				// todo  bitget ws data handler, spot and feature

			case <-t.resourceCtx.Done():
				log.Info("stop bitget task in work")
				return nil

			}

		}
	})
	return nil
}

func (t *BitGetTask) Close() error {
	var result error
	t.resourceCancel()
	t.ticker.Stop()
	log.Info("bitget task stopped")
	if err := t.tasks.Wait(); err != nil {
		result = errors.Join(result, fmt.Errorf("bitget ws task wait error"))
	}
	log.Info("bitget task stopped success")
	return nil
}
