package application

import (
	"context"
	"errors"
	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team17/task-service/internal/adapters/grpc"
	"gitlab.com/g6834/team17/task-service/internal/adapters/presenters"
	"gitlab.com/g6834/team17/task-service/internal/adapters/sender"
	"gitlab.com/g6834/team17/task-service/internal/messageHandler"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	dbMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/adapters/http"
	"gitlab.com/g6834/team17/task-service/internal/adapters/postgres"
	"gitlab.com/g6834/team17/task-service/internal/config"
	"gitlab.com/g6834/team17/task-service/internal/domain/task"
	"gitlab.com/g6834/team17/task-service/internal/utils"
	"golang.org/x/sync/errgroup"
)

var (
	srv    *http.Server
	logger *zerolog.Logger
)

func Start(ctx context.Context) {
	/* CONFIG init */
	if err := config.ReadConfigYML("config.yaml"); err != nil {
		log.Fatal("cannot read config file", err)
	}
	cfg := config.New()

	/* LOGGER init */
	logger = utils.NewLogger(cfg)

	/* DATABASE init */
	db, err := postgres.New(cfg, logger)
	if err != nil {
		logger.Error().Err(err).Msg("cannot initialize database")
	}
	defer db.Close()
	if cfg.Database.UseMigrations == "true" {
		err := runMigrations(db, cfg)
		if err != nil {
			logger.Error().Err(err).Msg("cannot up migrate")
		}
	}

	/* SERVICES init */
	//TODO: add dsn
	/* Sentry init */
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:   "https://40376d00d7b8408f8bc64950ce173be9@sentry.k8s.golang-mts-teta.ru/49",
		Debug: true,
	}); err != nil {
		logger.Error().Err(err).Msg("cannot init sentry")
	}
	defer sentry.Flush(2 * time.Second)

	/* Jaeger init */
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Jaeger.Host), jaeger.WithAgentPort(cfg.Jaeger.Port)))
	if err != nil {
		logger.Error().Err(err).Msg("cannot init jaeger collector")
		sentry.CaptureException(err)
	}

	/* Tracer provider init */
	traceProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(cfg.Jaeger.Service))))

	otel.SetTracerProvider(traceProvider)

	/* Auth Service init */
	//authService, err := grpc.NewAuthClient(cfg, logger)
	//if err != nil {
	//	sentry.CaptureException(err)
	//	logger.Error().Err(err).Msg("cannot initialize auth service")
	//}
	//defer authService.Close()

	/* Auth Mock Service init */
	authService := grpc.NewAuthClientMock(cfg, logger)

	/* Message Sender init */
	s := &sender.StdOut{}

	/* Task Service init */
	taskService := task.New(db, authService, s)

	/* Presenters and helpers init */
	presenters := presenters.New(logger)

	/* Message handler init */
	msgHandler := messageHandler.New(taskService, logger, cfg)

	/* Http Server init */
	srv, err = http.New(cfg, logger, authService, taskService, presenters)
	if err != nil {
		sentry.CaptureException(err)
		logger.Error().Err(err).Msg("cannot initialize server")
	}

	/* Start app */
	var g errgroup.Group
	g.Go(func() error {
		return srv.Start()
	})
	g.Go(func() error {
		return msgHandler.HandleMessages(ctx)
	})

	logger.Info().Msg("app is started")
	err = g.Wait()
	if err != nil {
		sentry.CaptureException(err)
		logger.Fatal().Err(err).Msg("http server start failed")
	}
}

func Stop() {
	logger.Warn().Msg("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(2)*time.Second)
	defer cancel()

	err := srv.Stop(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error while stopping")
	}

	logger.Warn().Msg("app has stopped")
}

func runMigrations(pg *postgres.Database, cfg *config.Config) error {
	// Migrations block
	driver, err := dbMigrate.WithInstance(pg.DB().DB, &dbMigrate.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file:"+cfg.Database.Migrations, cfg.Database.Name, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
