package ports

import "context"

type MessageHandler interface {
	HandleMessages(ctx context.Context) error
}
