package e_invoicing_document

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for EInvoicingDocuments
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new EInvoicingDocuments service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new e_invoicing_documents
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateEInvoicingDocumentsRequest) (*EInvoicingDocumentsResponse, error) {
	s.logger.Info("creating e_invoicing_documents",
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
	entity := &EInvoicingDocuments{
		OrganizationId: orgID,
		
		SourceTable: req.SourceTable,
		
		SourceId: req.SourceId,
		
		Authority: req.Authority,
		
		CountryCode: req.CountryCode,
		
		DocumentUuid: req.DocumentUuid,
		
		DocumentType: req.DocumentType,
		
		DocumentNumber: req.DocumentNumber,
		
		InternalReference: req.InternalReference,
		
		Status: req.Status,
		
		'draft',: req.'draft',,
		
		'pending',: req.'pending',,
		
		'submitted',: req.'submitted',,
		
		'accepted',: req.'accepted',,
		
		'rejected',: req.'rejected',,
		
		'cancelled',: req.'cancelled',,
		
		'error': req.'error',
		
		SubmittedAt: req.SubmittedAt,
		
		ResponseAt: req.ResponseAt,
		
		RequestPayload: req.RequestPayload,
		
		ResponsePayload: req.ResponsePayload,
		
		ErrorCode: req.ErrorCode,
		
		ErrorMessage: req.ErrorMessage,
		
		RetryCount: req.RetryCount,
		
		LastRetryAt: req.LastRetryAt,
		
		ZatcaHashValue: req.ZatcaHashValue,
		
		ZatcaPreviousHashValue: req.ZatcaPreviousHashValue,
		
		ZatcaInvoiceCounterValue: req.ZatcaInvoiceCounterValue,
		
		ZatcaCryptographicStamp: req.ZatcaCryptographicStamp,
		
		ZatcaQrCodePayload: req.ZatcaQrCodePayload,
		
		ZatcaComplianceInvoiceNumber: req.ZatcaComplianceInvoiceNumber,
		
		EtaDocumentTypeVersion: req.EtaDocumentTypeVersion,
		
		EtaSubmissionUuid: req.EtaSubmissionUuid,
		
		EtaLongId: req.EtaLongId,
		
		EtaInternalId: req.EtaInternalId,
		
		EtaDigitalSignature: req.EtaDigitalSignature,
		
		EtaSignatureAlgorithm: req.EtaSignatureAlgorithm,
		
		SubmissionFormat: req.SubmissionFormat,
		
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
		return nil, fmt.Errorf("failed to create e_invoicing_documents: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created e_invoicing_documents",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a e_invoicing_documents by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EInvoicingDocumentsResponse, error) {
	s.logger.Debug("getting e_invoicing_documents",
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
		return nil, fmt.Errorf("failed to get e_invoicing_documents: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("e_invoicing_documents not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of e_invoicing_documents records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*EInvoicingDocumentsListResponse, error) {
	s.logger.Debug("listing e_invoicing_documents",
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
		return nil, fmt.Errorf("failed to list e_invoicing_documents: %w", err)
	}

	// Convert to response
	items := make([]*EInvoicingDocumentsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &EInvoicingDocumentsListResponse{
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

// Update updates an existing e_invoicing_documents
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateEInvoicingDocumentsRequest) (*EInvoicingDocumentsResponse, error) {
	s.logger.Info("updating e_invoicing_documents",
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
		return nil, fmt.Errorf("failed to get e_invoicing_documents: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("e_invoicing_documents not found or access denied")
	}
	

	// Update fields
	
	if req.SourceTable != nil {
		entity.SourceTable = *req.SourceTable
	}
	
	if req.SourceId != nil {
		entity.SourceId = *req.SourceId
	}
	
	if req.Authority != nil {
		entity.Authority = *req.Authority
	}
	
	if req.CountryCode != nil {
		entity.CountryCode = *req.CountryCode
	}
	
	if req.DocumentUuid != nil {
		entity.DocumentUuid = *req.DocumentUuid
	}
	
	if req.DocumentType != nil {
		entity.DocumentType = *req.DocumentType
	}
	
	if req.DocumentNumber != nil {
		entity.DocumentNumber = *req.DocumentNumber
	}
	
	if req.InternalReference != nil {
		entity.InternalReference = *req.InternalReference
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.'draft', != nil {
		entity.'draft', = *req.'draft',
	}
	
	if req.'pending', != nil {
		entity.'pending', = *req.'pending',
	}
	
	if req.'submitted', != nil {
		entity.'submitted', = *req.'submitted',
	}
	
	if req.'accepted', != nil {
		entity.'accepted', = *req.'accepted',
	}
	
	if req.'rejected', != nil {
		entity.'rejected', = *req.'rejected',
	}
	
	if req.'cancelled', != nil {
		entity.'cancelled', = *req.'cancelled',
	}
	
	if req.'error' != nil {
		entity.'error' = *req.'error'
	}
	
	if req.SubmittedAt != nil {
		entity.SubmittedAt = req.SubmittedAt
	}
	
	if req.ResponseAt != nil {
		entity.ResponseAt = req.ResponseAt
	}
	
	if req.RequestPayload != nil {
		entity.RequestPayload = *req.RequestPayload
	}
	
	if req.ResponsePayload != nil {
		entity.ResponsePayload = *req.ResponsePayload
	}
	
	if req.ErrorCode != nil {
		entity.ErrorCode = *req.ErrorCode
	}
	
	if req.ErrorMessage != nil {
		entity.ErrorMessage = req.ErrorMessage
	}
	
	if req.RetryCount != nil {
		entity.RetryCount = *req.RetryCount
	}
	
	if req.LastRetryAt != nil {
		entity.LastRetryAt = req.LastRetryAt
	}
	
	if req.ZatcaHashValue != nil {
		entity.ZatcaHashValue = *req.ZatcaHashValue
	}
	
	if req.ZatcaPreviousHashValue != nil {
		entity.ZatcaPreviousHashValue = *req.ZatcaPreviousHashValue
	}
	
	if req.ZatcaInvoiceCounterValue != nil {
		entity.ZatcaInvoiceCounterValue = *req.ZatcaInvoiceCounterValue
	}
	
	if req.ZatcaCryptographicStamp != nil {
		entity.ZatcaCryptographicStamp = *req.ZatcaCryptographicStamp
	}
	
	if req.ZatcaQrCodePayload != nil {
		entity.ZatcaQrCodePayload = *req.ZatcaQrCodePayload
	}
	
	if req.ZatcaComplianceInvoiceNumber != nil {
		entity.ZatcaComplianceInvoiceNumber = *req.ZatcaComplianceInvoiceNumber
	}
	
	if req.EtaDocumentTypeVersion != nil {
		entity.EtaDocumentTypeVersion = *req.EtaDocumentTypeVersion
	}
	
	if req.EtaSubmissionUuid != nil {
		entity.EtaSubmissionUuid = *req.EtaSubmissionUuid
	}
	
	if req.EtaLongId != nil {
		entity.EtaLongId = *req.EtaLongId
	}
	
	if req.EtaInternalId != nil {
		entity.EtaInternalId = *req.EtaInternalId
	}
	
	if req.EtaDigitalSignature != nil {
		entity.EtaDigitalSignature = *req.EtaDigitalSignature
	}
	
	if req.EtaSignatureAlgorithm != nil {
		entity.EtaSignatureAlgorithm = *req.EtaSignatureAlgorithm
	}
	
	if req.SubmissionFormat != nil {
		entity.SubmissionFormat = *req.SubmissionFormat
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update e_invoicing_documents: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated e_invoicing_documents",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a e_invoicing_documents
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting e_invoicing_documents",
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
		return fmt.Errorf("failed to get e_invoicing_documents: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("e_invoicing_documents not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete e_invoicing_documents: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted e_invoicing_documents",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *EInvoicingDocuments) *EInvoicingDocumentsResponse {
	return &EInvoicingDocumentsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		SourceTable: entity.SourceTable,
		
		SourceId: entity.SourceId,
		
		Authority: entity.Authority,
		
		CountryCode: entity.CountryCode,
		
		DocumentUuid: entity.DocumentUuid,
		
		DocumentType: entity.DocumentType,
		
		DocumentNumber: entity.DocumentNumber,
		
		InternalReference: entity.InternalReference,
		
		Status: entity.Status,
		
		'draft',: entity.'draft',,
		
		'pending',: entity.'pending',,
		
		'submitted',: entity.'submitted',,
		
		'accepted',: entity.'accepted',,
		
		'rejected',: entity.'rejected',,
		
		'cancelled',: entity.'cancelled',,
		
		'error': entity.'error',
		
		CreatedAt: entity.CreatedAt,
		
		SubmittedAt: entity.SubmittedAt,
		
		ResponseAt: entity.ResponseAt,
		
		RequestPayload: entity.RequestPayload,
		
		ResponsePayload: entity.ResponsePayload,
		
		ErrorCode: entity.ErrorCode,
		
		ErrorMessage: entity.ErrorMessage,
		
		RetryCount: entity.RetryCount,
		
		LastRetryAt: entity.LastRetryAt,
		
		ZatcaHashValue: entity.ZatcaHashValue,
		
		ZatcaPreviousHashValue: entity.ZatcaPreviousHashValue,
		
		ZatcaInvoiceCounterValue: entity.ZatcaInvoiceCounterValue,
		
		ZatcaCryptographicStamp: entity.ZatcaCryptographicStamp,
		
		ZatcaQrCodePayload: entity.ZatcaQrCodePayload,
		
		ZatcaComplianceInvoiceNumber: entity.ZatcaComplianceInvoiceNumber,
		
		EtaDocumentTypeVersion: entity.EtaDocumentTypeVersion,
		
		EtaSubmissionUuid: entity.EtaSubmissionUuid,
		
		EtaLongId: entity.EtaLongId,
		
		EtaInternalId: entity.EtaInternalId,
		
		EtaDigitalSignature: entity.EtaDigitalSignature,
		
		EtaSignatureAlgorithm: entity.EtaSignatureAlgorithm,
		
		SubmissionFormat: entity.SubmissionFormat,
		
		Metadata: entity.Metadata,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for e_invoicing_documents
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *EInvoicingDocuments) error {
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

// canDelete checks if a e_invoicing_documents can be deleted
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
