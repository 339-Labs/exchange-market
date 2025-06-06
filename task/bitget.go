package task

import (
	"context"
	"github.com/339-Labs/exchange-market/common/tasks"
	"time"
)

type BitGetTask struct {
	timeout  time.Duration
	interval time.Duration
	ctx      *context.Context
	cancel   context.CancelFunc
	task     tasks.Group
}
