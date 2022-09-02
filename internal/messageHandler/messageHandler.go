package messageHandler

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/config"
	"gitlab.com/g6834/team17/task-service/internal/ports"
	"time"
)

type MessageHandler struct {
	task ports.Task
	log  *zerolog.Logger
}

type MessageHandlerWithRateLimiter struct {
	task        ports.Task
	log         *zerolog.Logger
	ticker      *time.Ticker
	maxRequests int64
}

func New(task ports.Task, logger *zerolog.Logger, cfg *config.Config) ports.MessageHandler {
	if cfg.MsgHandler.RatePeriodMicroseconds > 0 {
		return &MessageHandlerWithRateLimiter{
			task:        task,
			log:         logger,
			ticker:      time.NewTicker(time.Duration(cfg.MsgHandler.RatePeriodMicroseconds) * time.Second),
			maxRequests: cfg.MsgHandler.RequestsPerPeriod,
		}
	}
	return &MessageHandler{
		task: task,
		log:  logger,
	}
}

func (mh *MessageHandler) HandleMessages(ctx context.Context) error {
	for {
		tasks, err := mh.task.List(ctx)
		if err != nil {
			mh.log.Err(err)
		}

		for _, t := range tasks {
			err := mh.task.Send(ctx, t.ID)
			if err != nil {
				mh.log.Err(err)
			}
		}
	}
}

func (mh *MessageHandlerWithRateLimiter) HandleMessages(ctx context.Context) error {
	for {
		select {
		case <-mh.ticker.C:
			requestCounter := 0
			for requestCounter < int(mh.maxRequests) {
				tasks, err := mh.task.List(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return nil
					}
					mh.log.Err(err)
					continue
				}
				for _, t := range tasks {
					err := mh.task.Send(ctx, t.ID)
					if err != nil {
						mh.log.Err(err)
					}
				}
				mh.log.Printf("successfully sent msg")
				requestCounter++
			}
		case <-ctx.Done():
			return nil
		}
	}
}
