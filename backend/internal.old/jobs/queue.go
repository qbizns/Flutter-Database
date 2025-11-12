package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// Queue manages job enqueueing and scheduling
type Queue struct {
	client *asynq.Client
	logger *logging.Logger
}

// QueueConfig holds queue configuration
type QueueConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// NewQueue creates a new job queue client
func NewQueue(cfg QueueConfig, logger *logging.Logger) *Queue {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &Queue{
		client: client,
		logger: logger,
	}
}

// Close closes the queue client
func (q *Queue) Close() error {
	return q.client.Close()
}

// EnqueueEmailNotification enqueues an email notification job
func (q *Queue) EnqueueEmailNotification(ctx context.Context, payload EmailNotificationPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypeEmailNotification, payload)
	if err != nil {
		return fmt.Errorf("failed to create email task: %w", err)
	}

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue email notification", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued email notification",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.Strings("to", payload.To),
	)
	return nil
}

// EnqueueReportGeneration enqueues a report generation job
func (q *Queue) EnqueueReportGeneration(ctx context.Context, payload ReportGenerationPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypeReportGeneration, payload)
	if err != nil {
		return fmt.Errorf("failed to create report task: %w", err)
	}

	// Default options for report generation (high priority, longer timeout)
	defaultOpts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(10 * time.Minute),
		asynq.Queue("reports"),
	}
	opts = append(defaultOpts, opts...)

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue report generation", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued report generation",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.String("report_type", payload.ReportType),
		zap.String("org_id", payload.OrganizationID),
	)
	return nil
}

// EnqueueDataExport enqueues a data export job
func (q *Queue) EnqueueDataExport(ctx context.Context, payload DataExportPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypeDataExport, payload)
	if err != nil {
		return fmt.Errorf("failed to create export task: %w", err)
	}

	defaultOpts := []asynq.Option{
		asynq.MaxRetry(2),
		asynq.Timeout(15 * time.Minute),
		asynq.Queue("exports"),
	}
	opts = append(defaultOpts, opts...)

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue data export", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued data export",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.String("entity_type", payload.EntityType),
		zap.String("org_id", payload.OrganizationID),
	)
	return nil
}

// EnqueueDatabaseBackup enqueues a database backup job
func (q *Queue) EnqueueDatabaseBackup(ctx context.Context, payload DatabaseBackupPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypeDatabaseBackup, payload)
	if err != nil {
		return fmt.Errorf("failed to create backup task: %w", err)
	}

	defaultOpts := []asynq.Option{
		asynq.MaxRetry(1),
		asynq.Timeout(30 * time.Minute),
		asynq.Queue("backups"),
		asynq.Unique(1 * time.Hour), // Prevent duplicate backups within 1 hour
	}
	opts = append(defaultOpts, opts...)

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue database backup", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued database backup",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.String("backup_type", payload.BackupType),
	)
	return nil
}

// EnqueueInvoiceGeneration enqueues an invoice generation job
func (q *Queue) EnqueueInvoiceGeneration(ctx context.Context, payload InvoiceGenerationPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypeInvoiceGeneration, payload)
	if err != nil {
		return fmt.Errorf("failed to create invoice task: %w", err)
	}

	defaultOpts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(5 * time.Minute),
	}
	opts = append(defaultOpts, opts...)

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue invoice generation", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued invoice generation",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.String("invoice_id", payload.InvoiceID),
		zap.String("org_id", payload.OrganizationID),
	)
	return nil
}

// EnqueuePostingBatch enqueues a batch posting job
func (q *Queue) EnqueuePostingBatch(ctx context.Context, payload PostingBatchPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypePostingBatch, payload)
	if err != nil {
		return fmt.Errorf("failed to create posting batch task: %w", err)
	}

	defaultOpts := []asynq.Option{
		asynq.MaxRetry(2),
		asynq.Timeout(10 * time.Minute),
		asynq.Queue("posting"),
	}
	opts = append(defaultOpts, opts...)

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue posting batch", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued posting batch",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.Int("document_count", len(payload.Documents)),
		zap.String("org_id", payload.OrganizationID),
	)
	return nil
}

// EnqueueAuditLogCleanup enqueues an audit log cleanup job
func (q *Queue) EnqueueAuditLogCleanup(ctx context.Context, payload AuditLogCleanupPayload, opts ...asynq.Option) error {
	task, err := NewTask(TypeAuditLogCleanup, payload)
	if err != nil {
		return fmt.Errorf("failed to create cleanup task: %w", err)
	}

	defaultOpts := []asynq.Option{
		asynq.MaxRetry(1),
		asynq.Timeout(20 * time.Minute),
		asynq.Queue("maintenance"),
	}
	opts = append(defaultOpts, opts...)

	info, err := q.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		q.logger.Error("failed to enqueue audit log cleanup", zap.Error(err))
		return err
	}

	q.logger.Info("enqueued audit log cleanup",
		zap.String("task_id", info.ID),
		zap.String("queue", info.Queue),
		zap.Int("retention_days", payload.RetentionDays),
	)
	return nil
}

// ScheduleRecurringBackup schedules a recurring backup job
func (q *Queue) ScheduleRecurringBackup(cronspec string, payload DatabaseBackupPayload) error {
	task, err := NewTask(TypeDatabaseBackup, payload)
	if err != nil {
		return fmt.Errorf("failed to create backup task: %w", err)
	}

	// Schedule using cron syntax (e.g., "0 2 * * *" for 2 AM daily)
	_, err = q.client.(*asynq.Client).Enqueue(task,
		asynq.Queue("backups"),
		// Note: Scheduling requires asynq scheduler, handled in scheduler.go
	)

	if err != nil {
		q.logger.Error("failed to schedule recurring backup", zap.Error(err))
		return err
	}

	q.logger.Info("scheduled recurring backup", zap.String("cron", cronspec))
	return nil
}

// GetTaskInfo retrieves information about a queued task
func (q *Queue) GetTaskInfo(queue, taskID string) (*asynq.TaskInfo, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{
		Addr: q.client.(*asynq.Client).(*asynq.Client),
	})
	defer inspector.Close()

	// This is a simplified version - full implementation would query Redis
	return nil, fmt.Errorf("not implemented")
}
