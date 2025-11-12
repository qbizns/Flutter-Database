package vendor_bill

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for VendorBills
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new VendorBills service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new vendor_bills
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateVendorBillsRequest) (*VendorBillsResponse, error) {
	s.logger.Info("creating vendor_bills",
		zap.String("organization_id", orgID.String()),
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

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Convert DTO to entity
	entity := &VendorBills{
		OrganizationID: orgID,
		
		BillNumber: req.BillNumber,
		
		VendorBillNumber: req.VendorBillNumber,
		
		SupplierId: req.SupplierId,
		
		BillDate: req.BillDate,
		
		DueDate: req.DueDate,
		
		PaymentTerms: req.PaymentTerms,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		Subtotal: req.Subtotal,
		
		TaxAmount: req.TaxAmount,
		
		TotalAmount: req.TotalAmount,
		
		PaidAmount: req.PaidAmount,
		
		BalanceDue: req.BalanceDue,
		
		Status: req.Status,
		
		JournalEntryId: req.JournalEntryId,
		
		IsPosted: req.IsPosted,
		
		PurchaseOrderId: req.PurchaseOrderId,
		
		Description: req.Description,
		
		Notes: req.Notes,
		
		Memo: req.Memo,
		
		Attachments: req.Attachments,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create vendor_bills: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created vendor_bills",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a vendor_bills by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*VendorBillsResponse, error) {
	s.logger.Debug("getting vendor_bills",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor_bills: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("vendor_bills not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of vendor_bills records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*VendorBillsListResponse, error) {
	s.logger.Debug("listing vendor_bills",
		zap.String("organization_id", orgID.String()),
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

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	// Get from database (organization-scoped)
	entities, total, err := s.repo.ListByOrganization(ctx, tx, orgID, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list vendor_bills: %w", err)
	}

	// Convert to response
	items := make([]*VendorBillsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &VendorBillsListResponse{
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

// Update updates an existing vendor_bills
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateVendorBillsRequest) (*VendorBillsResponse, error) {
	s.logger.Info("updating vendor_bills",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
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

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor_bills: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("vendor_bills not found or access denied")
	}
	

	// Update fields
	
	if req.BillNumber != nil {
		entity.BillNumber = *req.BillNumber
	}
	
	if req.VendorBillNumber != nil {
		entity.VendorBillNumber = *req.VendorBillNumber
	}
	
	if req.SupplierId != nil {
		entity.SupplierId = *req.SupplierId
	}
	
	if req.BillDate != nil {
		entity.BillDate = *req.BillDate
	}
	
	if req.DueDate != nil {
		entity.DueDate = *req.DueDate
	}
	
	if req.PaymentTerms != nil {
		entity.PaymentTerms = *req.PaymentTerms
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = *req.AccountingPeriodId
	}
	
	if req.Subtotal != nil {
		entity.Subtotal = *req.Subtotal
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = *req.TaxAmount
	}
	
	if req.TotalAmount != nil {
		entity.TotalAmount = *req.TotalAmount
	}
	
	if req.PaidAmount != nil {
		entity.PaidAmount = *req.PaidAmount
	}
	
	if req.BalanceDue != nil {
		entity.BalanceDue = *req.BalanceDue
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.JournalEntryId != nil {
		entity.JournalEntryId = *req.JournalEntryId
	}
	
	if req.IsPosted != nil {
		entity.IsPosted = *req.IsPosted
	}
	
	if req.PurchaseOrderId != nil {
		entity.PurchaseOrderId = *req.PurchaseOrderId
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.Memo != nil {
		entity.Memo = *req.Memo
	}
	
	if req.Attachments != nil {
		entity.Attachments = *req.Attachments
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = *req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update vendor_bills: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated vendor_bills",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a vendor_bills
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting vendor_bills",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return err
	}

	// Verify ownership
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get vendor_bills: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("vendor_bills not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete vendor_bills: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted vendor_bills",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *VendorBills) *VendorBillsResponse {
	return &VendorBillsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		BillNumber: entity.BillNumber,
		
		VendorBillNumber: entity.VendorBillNumber,
		
		SupplierId: entity.SupplierId,
		
		BillDate: entity.BillDate,
		
		DueDate: entity.DueDate,
		
		PaymentTerms: entity.PaymentTerms,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		Subtotal: entity.Subtotal,
		
		TaxAmount: entity.TaxAmount,
		
		TotalAmount: entity.TotalAmount,
		
		PaidAmount: entity.PaidAmount,
		
		BalanceDue: entity.BalanceDue,
		
		Status: entity.Status,
		
		JournalEntryId: entity.JournalEntryId,
		
		IsPosted: entity.IsPosted,
		
		PurchaseOrderId: entity.PurchaseOrderId,
		
		Description: entity.Description,
		
		Notes: entity.Notes,
		
		Memo: entity.Memo,
		
		Attachments: entity.Attachments,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
	}
}


// setOrganizationContext sets the organization context for RLS
func (s *Service) setOrganizationContext(ctx context.Context, tx pgx.Tx, orgID uuid.UUID) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.current_organization_id = $1", orgID)
	if err != nil {
		return fmt.Errorf("failed to set organization context: %w", err)
	}
	return nil
}


// validateBusinessRules validates business rules for vendor_bills
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *VendorBills) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a vendor_bills can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
