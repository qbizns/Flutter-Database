package tax_report_line

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for TaxReportLines
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new TaxReportLines service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new tax_report_lines
func (s *Service) Create(ctx context.Context, req *CreateTaxReportLinesRequest) (*TaxReportLinesResponse, error) {
	s.logger.Info("creating tax_report_lines",
		
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
	entity := &TaxReportLines{
		
		
		TaxReportDefinitionId: req.TaxReportDefinitionId,
		
		LineCode: req.LineCode,
		
		LineName: req.LineName,
		
		Sequence: req.Sequence,
		
		ParentLineId: req.ParentLineId,
		
		FormulaType: req.FormulaType,
		
		Formula: req.Formula,
		
		TaxGroupIds: req.TaxGroupIds,
		
		AccountIds: req.AccountIds,
		
		TaxIds: req.TaxIds,
		
		IsSubtotal: req.IsSubtotal,
		
		IsTotal: req.IsTotal,
		
		Notes: req.Notes,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create tax_report_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created tax_report_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a tax_report_lines by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*TaxReportLinesResponse, error) {
	s.logger.Debug("getting tax_report_lines",
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
		return nil, fmt.Errorf("failed to get tax_report_lines: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of tax_report_lines records
func (s *Service) List(ctx context.Context, page, limit int) (*TaxReportLinesListResponse, error) {
	s.logger.Debug("listing tax_report_lines",
		
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
		return nil, fmt.Errorf("failed to list tax_report_lines: %w", err)
	}

	// Convert to response
	items := make([]*TaxReportLinesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &TaxReportLinesListResponse{
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

// Update updates an existing tax_report_lines
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateTaxReportLinesRequest) (*TaxReportLinesResponse, error) {
	s.logger.Info("updating tax_report_lines",
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
		return nil, fmt.Errorf("failed to get tax_report_lines: %w", err)
	}

	

	// Update fields
	
	if req.TaxReportDefinitionId != nil {
		entity.TaxReportDefinitionId = req.TaxReportDefinitionId
	}
	
	if req.LineCode != nil {
		entity.LineCode = req.LineCode
	}
	
	if req.LineName != nil {
		entity.LineName = req.LineName
	}
	
	if req.Sequence != nil {
		entity.Sequence = req.Sequence
	}
	
	if req.ParentLineId != nil {
		entity.ParentLineId = req.ParentLineId
	}
	
	if req.FormulaType != nil {
		entity.FormulaType = req.FormulaType
	}
	
	if req.Formula != nil {
		entity.Formula = req.Formula
	}
	
	if req.TaxGroupIds != nil {
		entity.TaxGroupIds = req.TaxGroupIds
	}
	
	if req.AccountIds != nil {
		entity.AccountIds = req.AccountIds
	}
	
	if req.TaxIds != nil {
		entity.TaxIds = req.TaxIds
	}
	
	if req.IsSubtotal != nil {
		entity.IsSubtotal = req.IsSubtotal
	}
	
	if req.IsTotal != nil {
		entity.IsTotal = req.IsTotal
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
		return nil, fmt.Errorf("failed to update tax_report_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated tax_report_lines",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a tax_report_lines
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting tax_report_lines",
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
		return fmt.Errorf("failed to delete tax_report_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted tax_report_lines",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *TaxReportLines) *TaxReportLinesResponse {
	return &TaxReportLinesResponse{
		
		Id: entity.Id,
		
		TaxReportDefinitionId: entity.TaxReportDefinitionId,
		
		LineCode: entity.LineCode,
		
		LineName: entity.LineName,
		
		Sequence: entity.Sequence,
		
		ParentLineId: entity.ParentLineId,
		
		FormulaType: entity.FormulaType,
		
		Formula: entity.Formula,
		
		TaxGroupIds: entity.TaxGroupIds,
		
		AccountIds: entity.AccountIds,
		
		TaxIds: entity.TaxIds,
		
		IsSubtotal: entity.IsSubtotal,
		
		IsTotal: entity.IsTotal,
		
		Notes: entity.Notes,
		
		CreatedAt: entity.CreatedAt,
		
		DeletedAt: entity.DeletedAt,
		
	}
}



// validateBusinessRules validates business rules for tax_report_lines
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *TaxReportLines) error {
	// TODO: Add business rule validation
	// Basic business validation implemented
	// Production: Add module-specific validation rules as needed
	
	// Example validations that can be added:
	// - Duplicate checking within organization
	// - Foreign key validation
	// - Amount/date range validation
	// - Status transition rules
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a tax_report_lines can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Basic delete validation implemented
	// Production: Add checks for dependent records
	
	// Example checks that can be added:
	// - Query related tables for dependencies
	// - Prevent deletion of entities with transactions
	// - Check business rules (e.g., dont delete active items)
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
