package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Worker processes background jobs
type Worker struct {
	server   *asynq.Server
	mux      *asynq.ServeMux
	handlers *Handlers
	logger   *logging.Logger
}

// WorkerConfig holds worker configuration
type WorkerConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int // Number of concurrent workers
	Queues        map[string]int // Queue priorities
}

// NewWorker creates a new job worker
func NewWorker(cfg WorkerConfig, handlers *Handlers, logger *logging.Logger) *Worker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		},
		asynq.Config{
			Concurrency: cfg.Concurrency,
			Queues:      cfg.Queues,
			// Retry configuration
			RetryDelayFunc: func(n int, err error, task *asynq.Task) time.Duration {
				// Exponential backoff: 1min, 5min, 30min
				return time.Duration(1<<uint(n)) * time.Minute
			},
			// Error handling
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				logger.Error("job execution failed",
					zap.String("type", task.Type()),
					zap.Error(err),
					zap.Int("retry_count", task.ResultWriter().(*asynq.ResultWriter).Retried()),
				)
			}),
			// Logging
			Logger: &asynqLogger{logger: logger},
		},
	)

	mux := asynq.NewServeMux()

	worker := &Worker{
		server:   server,
		mux:      mux,
		handlers: handlers,
		logger:   logger,
	}

	worker.registerHandlers()
	return worker
}

// registerHandlers registers all job handlers
func (w *Worker) registerHandlers() {
	// Register each job type with its handler
	w.mux.HandleFunc(TypeEmailNotification, w.wrapHandler(w.handlers.HandleEmailNotification))
	w.mux.HandleFunc(TypeReportGeneration, w.wrapHandler(w.handlers.HandleReportGeneration))
	w.mux.HandleFunc(TypeDataExport, w.wrapHandler(w.handlers.HandleDataExport))
	w.mux.HandleFunc(TypeDatabaseBackup, w.wrapHandler(w.handlers.HandleDatabaseBackup))
	w.mux.HandleFunc(TypeInvoiceGeneration, w.wrapHandler(w.handlers.HandleInvoiceGeneration))
	w.mux.HandleFunc(TypePostingBatch, w.wrapHandler(w.handlers.HandlePostingBatch))
	w.mux.HandleFunc(TypeAuditLogCleanup, w.wrapHandler(w.handlers.HandleAuditLogCleanup))
}

// wrapHandler wraps a handler function with metrics and logging
func (w *Worker) wrapHandler(handler func(context.Context, *asynq.Task) error) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		startTime := time.Now()

		w.logger.Info("processing job",
			zap.String("type", task.Type()),
			zap.String("payload_size", fmt.Sprintf("%d bytes", len(task.Payload()))),
		)

		// Execute the handler
		err := handler(ctx, task)

		duration := time.Since(startTime)
		status := "success"
		if err != nil {
			status = "error"
			w.logger.Error("job failed",
				zap.String("type", task.Type()),
				zap.Error(err),
				zap.Duration("duration", duration),
			)
		} else {
			w.logger.Info("job completed",
				zap.String("type", task.Type()),
				zap.Duration("duration", duration),
			)
		}

		// Record metrics
		metrics.RecordJobExecution(task.Type(), status, duration)

		return err
	}
}

// Start starts the worker server
func (w *Worker) Start() error {
	w.logger.Info("starting job worker",
		zap.Int("concurrency", w.server.(*asynq.Server).Config().Concurrency),
	)

	if err := w.server.Run(w.mux); err != nil {
		w.logger.Error("worker failed to start", zap.Error(err))
		return err
	}

	return nil
}

// Stop gracefully stops the worker server
func (w *Worker) Stop(ctx context.Context) error {
	w.logger.Info("stopping job worker")

	w.server.Stop()
	w.server.Shutdown()

	w.logger.Info("job worker stopped")
	return nil
}

// GetServerInfo returns server information
func (w *Worker) GetServerInfo() *asynq.ServerInfo {
	// In production, you'd query the server for actual stats
	return &asynq.ServerInfo{
		// This would contain real stats from Redis
	}
}

// asynqLogger adapts our logger to asynq's logger interface
type asynqLogger struct {
	logger *logging.Logger
}

func (l *asynqLogger) Debug(args ...interface{}) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *asynqLogger) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *asynqLogger) Warn(args ...interface{}) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *asynqLogger) Error(args ...interface{}) {
	l.logger.Error(fmt.Sprint(args...))
}

func (l *asynqLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(fmt.Sprint(args...))
}
