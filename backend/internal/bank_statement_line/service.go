package bank_statement_line

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for BankStatementLines
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new BankStatementLines service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new bank_statement_lines
func (s *Service) Create(ctx context.Context, req *CreateBankStatementLinesRequest) (*BankStatementLinesResponse, error) {
	s.logger.Info("creating bank_statement_lines",
		
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Convert DTO to entity
	entity := &BankStatementLines{
		
		
		BankStatementId: req.BankStatementId,
		
		LineNumber: req.LineNumber,
		
		TransactionDate: req.TransactionDate,
		
		ValueDate: req.ValueDate,
		
		Amount: req.Amount,
		
		CurrencyCode: req.CurrencyCode,
		
		Description: req.Description,
		
		Reference: req.Reference,
		
		CounterpartyName: req.CounterpartyName,
		
		CounterpartyAccount: req.CounterpartyAccount,
		
		BankReference: req.BankReference,
		
		Status: req.Status,
		
		Status: req.Status,
		
		Notes: req.Notes,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create bank_statement_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created bank_statement_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a bank_statement_lines by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*BankStatementLinesResponse, error) {
	s.logger.Debug("getting bank_statement_lines",
		zap.String("id", id.String()),
		
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank_statement_lines: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of bank_statement_lines records
func (s *Service) List(ctx context.Context, page, limit int) (*BankStatementLinesListResponse, error) {
	s.logger.Debug("listing bank_statement_lines",
		
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Get from database
	entities, total, err := s.repo.List(ctx, tx, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list bank_statement_lines: %w", err)
	}

	// Convert to response
	items := make([]*BankStatementLinesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &BankStatementLinesListResponse{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing bank_statement_lines
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateBankStatementLinesRequest) (*BankStatementLinesResponse, error) {
	s.logger.Info("updating bank_statement_lines",
		zap.String("id", id.String()),
		
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank_statement_lines: %w", err)
	}

	

	// Update fields
	
	if req.BankStatementId != nil {
		entity.BankStatementId = req.BankStatementId
	}
	
	if req.LineNumber != nil {
		entity.LineNumber = req.LineNumber
	}
	
	if req.TransactionDate != nil {
		entity.TransactionDate = req.TransactionDate
	}
	
	if req.ValueDate != nil {
		entity.ValueDate = req.ValueDate
	}
	
	if req.Amount != nil {
		entity.Amount = req.Amount
	}
	
	if req.CurrencyCode != nil {
		entity.CurrencyCode = req.CurrencyCode
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Reference != nil {
		entity.Reference = req.Reference
	}
	
	if req.CounterpartyName != nil {
		entity.CounterpartyName = req.CounterpartyName
	}
	
	if req.CounterpartyAccount != nil {
		entity.CounterpartyAccount = req.CounterpartyAccount
	}
	
	if req.BankReference != nil {
		entity.BankReference = req.BankReference
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update bank_statement_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated bank_statement_lines",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a bank_statement_lines
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting bank_statement_lines",
		zap.String("id", id.String()),
		
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete bank_statement_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted bank_statement_lines",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *BankStatementLines) *BankStatementLinesResponse {
	return &BankStatementLinesResponse{
		
		Id: entity.Id,
		
		BankStatementId: entity.BankStatementId,
		
		LineNumber: entity.LineNumber,
		
		TransactionDate: entity.TransactionDate,
		
		ValueDate: entity.ValueDate,
		
		Amount: entity.Amount,
		
		CurrencyCode: entity.CurrencyCode,
		
		Description: entity.Description,
		
		Reference: entity.Reference,
		
		CounterpartyName: entity.CounterpartyName,
		
		CounterpartyAccount: entity.CounterpartyAccount,
		
		BankReference: entity.BankReference,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		Notes: entity.Notes,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for bank_statement_lines
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *BankStatementLines) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a bank_statement_lines can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
