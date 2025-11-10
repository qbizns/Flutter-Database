package purchases

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

var poNumberRegex = regexp.MustCompile(`^[A-Z0-9\-_]{3,50}$`)

// Valid status values
const (
	StatusDraft      = "draft"
	StatusPending    = "pending"
	StatusApproved   = "approved"
	StatusOrdered    = "ordered"
	StatusPartial    = "partial"
	StatusReceived   = "received"
	StatusCancelled  = "cancelled"
)

// Valid payment status values
const (
	PaymentStatusPending = "pending"
	PaymentStatusPartial = "partial"
	PaymentStatusPaid    = "paid"
	PaymentStatusOverdue = "overdue"
)

// Service handles purchase order business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new purchase order service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// List retrieves purchase orders with filters
func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters PurchaseOrderFilters) ([]PurchaseOrder, error) {
	orders, err := s.repo.ListOrders(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list purchase orders", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return orders, nil
}

// Count returns total number of purchase orders matching filters
func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters PurchaseOrderFilters) (int64, error) {
	count, err := s.repo.CountOrders(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count purchase orders", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

// Create creates a new purchase order with line items
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, req *CreatePurchaseOrderRequest) (*PurchaseOrder, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check for duplicate PO number
	existing, err := s.repo.GetOrderByNumber(ctx, orgID, req.PONumber)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate PO number", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if existing != nil {
		return nil, apperrors.AlreadyExists("purchase_order", fmt.Sprintf("PO number %s already exists", req.PONumber))
	}

	// Create purchase order
	po := &PurchaseOrder{
		ID:             uuid.New(),
		OrganizationID: orgID,
		SupplierID:     req.SupplierID,
		PONumber:       req.PONumber,
		PODate:         req.PODate,
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		Status:         StatusDraft,
		PaymentStatus:  PaymentStatusPending,
		Notes:          req.Notes,
		TermsAndConditions: req.TermsAndConditions,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		CreatedBy:      userID,
	}

	if req.Status != "" && isValidStatus(req.Status) {
		po.Status = req.Status
	}

	// Create line items and calculate totals
	items := make([]PurchaseOrderItem, 0, len(req.Items))
	for _, itemReq := range req.Items {
		item := s.createItemFromRequest(orgID, po.ID, itemReq)
		items = append(items, item)

		// Add to totals
		po.Subtotal += item.LineTotal
		po.TaxAmount += item.TaxAmount
		po.DiscountAmount += item.DiscountAmount
	}

	po.Items = items
	po.TotalAmount = po.Subtotal + po.TaxAmount + po.ShippingCost - po.DiscountAmount

	// Persist purchase order
	if err := s.repo.CreateOrder(ctx, po); err != nil {
		s.logger.Error("failed to create purchase order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	// Create line items
	for _, item := range items {
		item := item // Create a copy
		if err := s.repo.CreateItem(ctx, &item); err != nil {
			s.logger.Error("failed to create purchase order item", zap.Error(err), zap.String("po_id", po.ID.String()))
			return nil, apperrors.DatabaseError(err)
		}
	}

	return po, nil
}

// Get retrieves a purchase order with its line items
func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PurchaseOrder, error) {
	po, err := s.repo.GetOrder(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get purchase order", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}
	if po == nil {
		return nil, apperrors.NotFound("purchase_order")
	}

	// Get line items
	items, err := s.repo.ListItems(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get purchase order items", zap.Error(err), zap.String("po_id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}
	po.Items = items

	return po, nil
}

// Update updates a purchase order (draft status only)
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, userID *uuid.UUID, req *UpdatePurchaseOrderRequest) (*PurchaseOrder, error) {
	// Get existing purchase order
	po, err := s.repo.GetOrder(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get purchase order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if po == nil {
		return nil, apperrors.NotFound("purchase_order")
	}

	// Only draft orders can be updated
	if po.Status != StatusDraft {
		return nil, apperrors.ValidationFailed("only draft purchase orders can be updated")
	}

	// Validate update request
	if err := s.validateUpdateRequest(req, po); err != nil {
		return nil, err
	}

	// Check for duplicate PO number if changed
	if req.PONumber != "" && req.PONumber != po.PONumber {
		existing, err := s.repo.GetOrderByNumber(ctx, orgID, req.PONumber)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate PO number", zap.Error(err))
			return nil, apperrors.DatabaseError(err)
		}
		if existing != nil {
			return nil, apperrors.AlreadyExists("purchase_order", fmt.Sprintf("PO number %s already exists", req.PONumber))
		}
		po.PONumber = req.PONumber
	}

	// Update basic fields
	if req.ExpectedDeliveryDate != nil {
		po.ExpectedDeliveryDate = req.ExpectedDeliveryDate
	}
	if req.Status != "" && isValidStatus(req.Status) {
		po.Status = req.Status
	}
	if req.Notes != "" {
		po.Notes = req.Notes
	}
	if req.TermsAndConditions != "" {
		po.TermsAndConditions = req.TermsAndConditions
	}

	po.UpdatedAt = time.Now()
	if userID != nil {
		po.UpdatedBy = userID
	}

	// Update items if provided
	if len(req.Items) > 0 {
		// Delete existing items
		if err := s.repo.DeleteItemsByOrderID(ctx, orgID, id); err != nil {
			s.logger.Error("failed to delete items for update", zap.Error(err))
			return nil, apperrors.DatabaseError(err)
		}

		// Recalculate totals
		po.Subtotal = 0
		po.TaxAmount = 0
		po.DiscountAmount = 0

		newItems := make([]PurchaseOrderItem, 0, len(req.Items))
		for _, itemReq := range req.Items {
			item := PurchaseOrderItem{
				ID:              uuid.New(),
				OrganizationID:  orgID,
				PurchaseOrderID: po.ID,
				ProductID:       nil,
				ProductName:     itemReq.ProductName,
				ProductSKU:      itemReq.ProductSKU,
				QuantityOrdered: itemReq.QuantityOrdered,
				UnitOfMeasure:   itemReq.UnitOfMeasure,
				UnitCost:        itemReq.UnitCost,
				DiscountPercent: itemReq.DiscountPercent,
				TaxPercent:      itemReq.TaxPercent,
				Notes:           itemReq.Notes,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			// Calculate amounts
			item.DiscountAmount = (item.UnitCost * item.QuantityOrdered) * (item.DiscountPercent / 100)
			subtotalAfterDiscount := (item.UnitCost * item.QuantityOrdered) - item.DiscountAmount
			item.TaxAmount = subtotalAfterDiscount * (item.TaxPercent / 100)
			item.LineTotal = subtotalAfterDiscount + item.TaxAmount

			po.Subtotal += item.LineTotal - item.TaxAmount
			po.TaxAmount += item.TaxAmount
			po.DiscountAmount += item.DiscountAmount

			newItems = append(newItems, item)
		}

		po.Items = newItems

		// Create new items
		for _, item := range newItems {
			item := item // Create a copy
			if err := s.repo.CreateItem(ctx, &item); err != nil {
				s.logger.Error("failed to create purchase order item on update", zap.Error(err))
				return nil, apperrors.DatabaseError(err)
			}
		}
	}

	po.TotalAmount = po.Subtotal + po.TaxAmount + po.ShippingCost - po.DiscountAmount

	// Update in database
	if err := s.repo.UpdateOrder(ctx, po); err != nil {
		s.logger.Error("failed to update purchase order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return po, nil
}

// Approve approves a purchase order
func (s *Service) Approve(ctx context.Context, orgID uuid.UUID, id uuid.UUID, userID uuid.UUID) (*PurchaseOrder, error) {
	po, err := s.repo.GetOrder(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get purchase order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if po == nil {
		return nil, apperrors.NotFound("purchase_order")
	}

	// Can approve from draft or pending status
	if po.Status != StatusDraft && po.Status != StatusPending {
		return nil, apperrors.ValidationFailed(fmt.Sprintf("cannot approve purchase order in %s status", po.Status))
	}

	po.Status = StatusApproved
	po.ApprovedBy = &userID
	now := time.Now()
	po.ApprovedAt = &now
	po.UpdatedAt = now
	po.UpdatedBy = &userID

	if err := s.repo.ApproveOrder(ctx, orgID, id, userID); err != nil {
		s.logger.Error("failed to approve purchase order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return po, nil
}

// UpdateStatus updates the purchase order status
func (s *Service) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, newStatus string, userID *uuid.UUID) (*PurchaseOrder, error) {
	if !isValidStatus(newStatus) {
		return nil, apperrors.ValidationFailed(fmt.Sprintf("invalid status: %s", newStatus))
	}

	po, err := s.repo.GetOrder(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get purchase order", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if po == nil {
		return nil, apperrors.NotFound("purchase_order")
	}

	// Validate status transition
	if !isValidStatusTransition(po.Status, newStatus) {
		return nil, apperrors.ValidationFailed(fmt.Sprintf("cannot transition from %s to %s", po.Status, newStatus))
	}

	po.Status = newStatus
	po.UpdatedAt = time.Now()
	if userID != nil {
		po.UpdatedBy = userID
	}

	if err := s.repo.UpdateOrderStatus(ctx, orgID, id, newStatus); err != nil {
		s.logger.Error("failed to update purchase order status", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return po, nil
}

// RecordItemReceipt records received quantity for an item
func (s *Service) RecordItemReceipt(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, itemID uuid.UUID, quantityReceived float64) (*PurchaseOrderItem, error) {
	item, err := s.repo.GetItem(ctx, orgID, itemID)
	if err != nil {
		s.logger.Error("failed to get purchase order item", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if item == nil {
		return nil, apperrors.NotFound("purchase_order_item")
	}

	if item.PurchaseOrderID != poID {
		return nil, apperrors.ValidationFailed("item does not belong to this purchase order")
	}

	if quantityReceived < 0 || quantityReceived > item.QuantityOrdered {
		return nil, apperrors.ValidationFailed(fmt.Sprintf("quantity received must be between 0 and %f", item.QuantityOrdered))
	}

	item.QuantityReceived = quantityReceived
	item.UpdatedAt = time.Now()

	if err := s.repo.UpdateItem(ctx, item); err != nil {
		s.logger.Error("failed to update receipt quantity", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	// Update PO status if all items received
	if err := s.updatePOStatusBasedOnReceipts(ctx, orgID, poID); err != nil {
		s.logger.Error("failed to update PO status", zap.Error(err))
		// Don't fail the entire operation, just log the error
	}

	return item, nil
}

// Delete soft-deletes a purchase order (draft only)
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	po, err := s.repo.GetOrder(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get purchase order", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if po == nil {
		return apperrors.NotFound("purchase_order")
	}

	// Only draft orders can be deleted
	if po.Status != StatusDraft {
		return apperrors.ValidationFailed("only draft purchase orders can be deleted")
	}

	if err := s.repo.DeleteOrder(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete purchase order", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// Helper functions

func (s *Service) validateCreateRequest(req *CreatePurchaseOrderRequest) error {
	if req.SupplierID == uuid.Nil {
		return apperrors.ValidationFailed("supplier_id is required")
	}
	if req.PONumber == "" {
		return apperrors.ValidationFailed("po_number is required")
	}
	if !poNumberRegex.MatchString(req.PONumber) {
		return apperrors.ValidationFailed("po_number must be 3-50 characters and contain only uppercase letters, numbers, hyphens, and underscores")
	}
	if req.PODate.IsZero() {
		return apperrors.ValidationFailed("po_date is required")
	}
	if len(req.Items) == 0 {
		return apperrors.ValidationFailed("at least one item is required")
	}
	for i, item := range req.Items {
		if err := s.validateOrderItem(item); err != nil {
			return apperrors.ValidationFailed(fmt.Sprintf("item %d: %v", i+1, err))
		}
	}
	return nil
}

func (s *Service) validateUpdateRequest(req *UpdatePurchaseOrderRequest, po *PurchaseOrder) error {
	if req.PONumber != "" && !poNumberRegex.MatchString(req.PONumber) {
		return apperrors.ValidationFailed("po_number must be 3-50 characters and contain only uppercase letters, numbers, hyphens, and underscores")
	}
	if req.Status != "" && !isValidStatus(req.Status) {
		return apperrors.ValidationFailed(fmt.Sprintf("invalid status: %s", req.Status))
	}
	for i, item := range req.Items {
		if err := s.validateUpdateOrderItem(item); err != nil {
			return apperrors.ValidationFailed(fmt.Sprintf("item %d: %v", i+1, err))
		}
	}
	return nil
}

func (s *Service) validateOrderItem(item CreatePurchaseOrderItemRequest) error {
	if item.ProductName == "" {
		return fmt.Errorf("product_name is required")
	}
	if item.QuantityOrdered <= 0 {
		return fmt.Errorf("quantity_ordered must be greater than 0")
	}
	if item.UnitCost < 0 {
		return fmt.Errorf("unit_cost must be non-negative")
	}
	if item.DiscountPercent < 0 || item.DiscountPercent > 100 {
		return fmt.Errorf("discount_percentage must be between 0 and 100")
	}
	if item.TaxPercent < 0 || item.TaxPercent > 100 {
		return fmt.Errorf("tax_percentage must be between 0 and 100")
	}
	return nil
}

func (s *Service) validateUpdateOrderItem(item UpdatePurchaseOrderItemRequest) error {
	if item.ProductName == "" {
		return fmt.Errorf("product_name is required")
	}
	if item.QuantityOrdered > 0 && item.QuantityOrdered < 0 {
		return fmt.Errorf("quantity_ordered must be positive")
	}
	if item.QuantityReceived < 0 {
		return fmt.Errorf("quantity_received must be non-negative")
	}
	if item.UnitCost < 0 {
		return fmt.Errorf("unit_cost must be non-negative")
	}
	return nil
}

func (s *Service) createItemFromRequest(orgID uuid.UUID, poID uuid.UUID, req CreatePurchaseOrderItemRequest) PurchaseOrderItem {
	item := PurchaseOrderItem{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		PurchaseOrderID: poID,
		ProductID:       req.ProductID,
		ProductName:     req.ProductName,
		ProductSKU:      req.ProductSKU,
		QuantityOrdered: req.QuantityOrdered,
		UnitOfMeasure:   req.UnitOfMeasure,
		UnitCost:        req.UnitCost,
		DiscountPercent: req.DiscountPercent,
		TaxPercent:      req.TaxPercent,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Calculate line amounts
	item.DiscountAmount = (item.UnitCost * item.QuantityOrdered) * (item.DiscountPercent / 100)
	subtotalAfterDiscount := (item.UnitCost * item.QuantityOrdered) - item.DiscountAmount
	item.TaxAmount = subtotalAfterDiscount * (item.TaxPercent / 100)
	item.LineTotal = subtotalAfterDiscount + item.TaxAmount

	return item
}

func (s *Service) updatePOStatusBasedOnReceipts(ctx context.Context, orgID uuid.UUID, poID uuid.UUID) error {
	po, err := s.repo.GetOrder(ctx, orgID, poID)
	if err != nil {
		return err
	}
	if po == nil {
		return apperrors.NotFound("purchase_order")
	}

	items, err := s.repo.ListItems(ctx, orgID, poID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	// Check if all items fully received
	allReceived := true
	anyReceived := false
	for _, item := range items {
		if item.QuantityReceived >= item.QuantityOrdered {
			continue
		}
		allReceived = false
		if item.QuantityReceived > 0 {
			anyReceived = true
		}
	}

	var newStatus string
	if allReceived {
		newStatus = StatusReceived
	} else if anyReceived {
		newStatus = StatusPartial
	} else {
		return nil // No change needed
	}

	if newStatus != po.Status {
		if err := s.repo.UpdateOrderStatus(ctx, orgID, poID, newStatus); err != nil {
			return err
		}
	}

	return nil
}

func isValidStatus(status string) bool {
	switch status {
	case StatusDraft, StatusPending, StatusApproved, StatusOrdered, StatusPartial, StatusReceived, StatusCancelled:
		return true
	default:
		return false
	}
}

func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		StatusDraft: {StatusPending, StatusCancelled},
		StatusPending: {StatusApproved, StatusCancelled},
		StatusApproved: {StatusOrdered, StatusCancelled},
		StatusOrdered: {StatusPartial, StatusReceived, StatusCancelled},
		StatusPartial: {StatusReceived, StatusCancelled},
		StatusReceived: {},
		StatusCancelled: {},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, v := range allowed {
		if v == to {
			return true
		}
	}
	return false
}
