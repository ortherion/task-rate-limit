package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gitlab.com/g6834/team17/task-service/internal/application"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	go application.Start(ctx)
	<-ctx.Done()
	application.Stop()
}
