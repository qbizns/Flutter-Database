package einvoicing

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// EInvoicingService handles e-invoicing business logic
type EInvoicingService struct {
	docRepo    EInvoicingDocumentRepository
	seqRepo    DocumentSequenceRepository
	logger     *logging.Logger
}

// NewEInvoicingService creates a new e-invoicing service
func NewEInvoicingService(
	docRepo EInvoicingDocumentRepository,
	seqRepo DocumentSequenceRepository,
	logger *logging.Logger,
) *EInvoicingService {
	return &EInvoicingService{
		docRepo: docRepo,
		seqRepo: seqRepo,
		logger:  logger,
	}
}

// Document CRUD Operations

// ListDocuments retrieves e-invoicing documents with filters
func (s *EInvoicingService) ListDocuments(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentFilters) ([]EInvoicingDocument, error) {
	docs, err := s.docRepo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list e-invoicing documents", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return docs, nil
}

// CountDocuments counts e-invoicing documents
func (s *EInvoicingService) CountDocuments(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentFilters) (int64, error) {
	count, err := s.docRepo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count e-invoicing documents", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

// CreateDocument creates a new e-invoicing document
func (s *EInvoicingService) CreateDocument(ctx context.Context, doc *EInvoicingDocument) error {
	if doc.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("Organization ID is required")
	}
	if doc.SourceTable == "" || doc.SourceID == uuid.Nil {
		return apperrors.ValidationFailed("Source table and ID are required")
	}
	if doc.Authority == "" {
		return apperrors.ValidationFailed("Authority is required (ZATCA, ETA, OTHER)")
	}

	doc.ID = uuid.New()
	doc.DocumentUUID = uuid.New()
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	if doc.Status == "" {
		doc.Status = "draft"
	}
	doc.RetryCount = 0

	if err := s.docRepo.Create(ctx, doc); err != nil {
		s.logger.Error("failed to create e-invoicing document", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	// Log creation event
	status := doc.Status
	_ = s.logEvent(ctx, doc.OrganizationID, doc.ID, "created", nil, &status, "Document created")
	return nil
}

// GetDocument retrieves an e-invoicing document
func (s *EInvoicingService) GetDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EInvoicingDocument, error) {
	doc, err := s.docRepo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get e-invoicing document", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}
	if doc == nil {
		return nil, apperrors.NotFound("e-invoicing document")
	}
	return doc, nil
}

// GetDocumentByNumber retrieves an e-invoicing document by document number
func (s *EInvoicingService) GetDocumentByNumber(ctx context.Context, orgID uuid.UUID, authority string, documentNumber string) (*EInvoicingDocument, error) {
	doc, err := s.docRepo.GetByDocumentNumber(ctx, orgID, authority, documentNumber)
	if err != nil {
		s.logger.Error("failed to get e-invoicing document by number", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return doc, nil
}

// UpdateDocument updates an e-invoicing document
func (s *EInvoicingService) UpdateDocument(ctx context.Context, doc *EInvoicingDocument) error {
	doc.UpdatedAt = time.Now()

	if err := s.docRepo.Update(ctx, doc); err != nil {
		s.logger.Error("failed to update e-invoicing document", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// UpdateDocumentStatus updates document status and creates event
func (s *EInvoicingService) UpdateDocumentStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, newStatus string, errorCode *string, errorMessage *string) error {
	doc, err := s.GetDocument(ctx, orgID, id)
	if err != nil {
		return err
	}

	oldStatus := doc.Status
	doc.Status = newStatus
	doc.ErrorCode = errorCode
	doc.ErrorMessage = errorMessage

	if newStatus == "submitted" && doc.SubmittedAt == nil {
		now := time.Now()
		doc.SubmittedAt = &now
	}
	if newStatus == "accepted" || newStatus == "rejected" {
		now := time.Now()
		doc.ResponseAt = &now
	}

	if err := s.UpdateDocument(ctx, doc); err != nil {
		return err
	}

	// Log status change event
	_ = s.logEvent(ctx, orgID, id, "status_check", &oldStatus, &newStatus, "Status updated to "+newStatus)
	return nil
}

// DeleteDocument deletes an e-invoicing document
func (s *EInvoicingService) DeleteDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.docRepo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete e-invoicing document", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// Event Management

// ListEvents retrieves events for an e-invoicing document
func (s *EInvoicingService) ListEvents(ctx context.Context, orgID uuid.UUID, filters EInvoicingDocumentEventFilters) ([]EInvoicingDocumentEvent, error) {
	events, err := s.docRepo.ListEvents(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list e-invoicing events", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return events, nil
}

// GetEvent retrieves a specific event
func (s *EInvoicingService) GetEvent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*EInvoicingDocumentEvent, error) {
	event, err := s.docRepo.GetEvent(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get event", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return event, nil
}

// logEvent creates an event record
func (s *EInvoicingService) logEvent(ctx context.Context, orgID uuid.UUID, docID uuid.UUID, eventType string, previousStatus, newStatus *string, description string) error {
	event := &EInvoicingDocumentEvent{
		ID:                   uuid.New(),
		OrganizationID:       orgID,
		EInvoicingDocumentID: docID,
		EventType:            eventType,
		EventTimestamp:       time.Now(),
		PreviousStatus:       previousStatus,
		NewStatus:            newStatus,
		EventDescription:     &description,
		EventData:            JSONB{},
		TriggeredBy:          "system",
	}

	if err := s.docRepo.CreateEvent(ctx, event); err != nil {
		s.logger.Error("failed to log event", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// Document Sequence Operations

// CreateSequence creates a new document sequence
func (s *EInvoicingService) CreateSequence(ctx context.Context, seq *DocumentSequence) error {
	if seq.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("Organization ID is required")
	}
	if seq.DocumentType == "" {
		return apperrors.ValidationFailed("Document type is required")
	}

	seq.ID = uuid.New()
	seq.CreatedAt = time.Now()
	seq.UpdatedAt = time.Now()
	if seq.NextNumber == 0 {
		seq.NextNumber = 1
	}
	if seq.Padding == 0 {
		seq.Padding = 6
	}
	if seq.IncrementBy == 0 {
		seq.IncrementBy = 1
	}
	if seq.ResetFrequency == "" {
		seq.ResetFrequency = "never"
	}
	seq.IsActive = true

	if err := s.seqRepo.Create(ctx, seq); err != nil {
		s.logger.Error("failed to create document sequence", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// GetSequence retrieves a document sequence
func (s *EInvoicingService) GetSequence(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DocumentSequence, error) {
	seq, err := s.seqRepo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get document sequence", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if seq == nil {
		return nil, apperrors.NotFound("document sequence")
	}
	return seq, nil
}

// GetSequenceByDocumentType retrieves a sequence by document type
func (s *EInvoicingService) GetSequenceByDocumentType(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (*DocumentSequence, error) {
	seq, err := s.seqRepo.GetByDocumentType(ctx, orgID, documentType, locationID)
	if err != nil {
		s.logger.Error("failed to get sequence by document type", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return seq, nil
}

// ListSequences retrieves document sequences with filters
func (s *EInvoicingService) ListSequences(ctx context.Context, orgID uuid.UUID, filters DocumentSequenceFilters) ([]DocumentSequence, error) {
	seqs, err := s.seqRepo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list document sequences", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return seqs, nil
}

// UpdateSequence updates a document sequence
func (s *EInvoicingService) UpdateSequence(ctx context.Context, seq *DocumentSequence) error {
	seq.UpdatedAt = time.Now()

	if err := s.seqRepo.Update(ctx, seq); err != nil {
		s.logger.Error("failed to update document sequence", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// DeleteSequence deletes a document sequence
func (s *EInvoicingService) DeleteSequence(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := s.seqRepo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete document sequence", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// GetNextNumber gets the next document number from the sequence
func (s *EInvoicingService) GetNextNumber(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (string, error) {
	num, err := s.seqRepo.GetNextNumber(ctx, orgID, documentType, locationID)
	if err != nil {
		s.logger.Error("failed to get next document number", zap.Error(err))
		return "", apperrors.DatabaseError(err)
	}
	return num, nil
}

// PreviewNumber previews the next document number without incrementing
func (s *EInvoicingService) PreviewNumber(ctx context.Context, orgID uuid.UUID, documentType string, locationID *uuid.UUID) (string, error) {
	num, err := s.seqRepo.PreviewNumber(ctx, orgID, documentType, locationID)
	if err != nil {
		s.logger.Error("failed to preview document number", zap.Error(err))
		return "", apperrors.DatabaseError(err)
	}
	return num, nil
}

// ResetSequence resets a document sequence
func (s *EInvoicingService) ResetSequence(ctx context.Context, orgID uuid.UUID, documentType string, resetTo int64, locationID *uuid.UUID) error {
	if err := s.seqRepo.ResetSequence(ctx, orgID, documentType, resetTo, locationID); err != nil {
		s.logger.Error("failed to reset document sequence", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}
