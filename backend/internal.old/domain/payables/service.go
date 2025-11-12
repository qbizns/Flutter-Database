package payables

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for payables data access
type Repository interface {
	// Vendor Bills
	CreateVendorBill(ctx context.Context, bill *VendorBill) error
	GetVendorBillByID(ctx context.Context, organizationID, billID uuid.UUID) (*VendorBill, error)
	GetVendorBillByNumber(ctx context.Context, organizationID uuid.UUID, billNumber string) (*VendorBill, error)
	ListVendorBills(ctx context.Context, organizationID uuid.UUID, opts *VendorBillFilterOptions) ([]VendorBill, int64, error)
	UpdateVendorBill(ctx context.Context, bill *VendorBill) error
	DeleteVendorBill(ctx context.Context, organizationID, billID uuid.UUID) error

	// Vendor Bill Lines
	CreateVendorBillLine(ctx context.Context, line *VendorBillLine) error
	GetVendorBillLineByID(ctx context.Context, organizationID, lineID uuid.UUID) (*VendorBillLine, error)
	ListVendorBillLines(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorBillLine, error)
	UpdateVendorBillLine(ctx context.Context, line *VendorBillLine) error
	DeleteVendorBillLine(ctx context.Context, organizationID, lineID uuid.UUID) error
	DeleteVendorBillLines(ctx context.Context, organizationID, billID uuid.UUID) error

	// Vendor Payments
	CreateVendorPayment(ctx context.Context, payment *VendorPayment) error
	GetVendorPaymentByID(ctx context.Context, organizationID, paymentID uuid.UUID) (*VendorPayment, error)
	GetVendorPaymentByNumber(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*VendorPayment, error)
	ListVendorPayments(ctx context.Context, organizationID uuid.UUID, opts *VendorPaymentFilterOptions) ([]VendorPayment, int64, error)
	UpdateVendorPayment(ctx context.Context, payment *VendorPayment) error
	DeleteVendorPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error

	// Vendor Payment Applications
	CreateVendorPaymentApplication(ctx context.Context, app *VendorPaymentApplication) error
	GetVendorPaymentApplicationByID(ctx context.Context, organizationID, appID uuid.UUID) (*VendorPaymentApplication, error)
	ListVendorPaymentApplications(ctx context.Context, organizationID, paymentID uuid.UUID) ([]VendorPaymentApplication, error)
	ListPaymentApplicationsByBill(ctx context.Context, organizationID, billID uuid.UUID) ([]VendorPaymentApplication, error)
	DeleteVendorPaymentApplication(ctx context.Context, organizationID, appID uuid.UUID) error

	// Reporting
	GetVendorAgingReport(ctx context.Context, organizationID uuid.UUID) ([]VendorAgingReport, error)
	GetBillsByDueDate(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]VendorBill, error)
	GetSupplierBalance(ctx context.Context, organizationID, supplierID uuid.UUID) (float64, error)
}

// Service provides business logic for Accounts Payable
type Service struct {
	repo Repository
}

// NewService creates a new payables service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateVendorBill creates a new vendor bill with line items
func (s *Service) CreateVendorBill(ctx context.Context, req *CreateVendorBillRequest, orgID, userID uuid.UUID) (*VendorBillWithLines, error) {
	if req.TotalAmount < 0 {
		return nil, errors.New("total amount must be non-negative")
	}

	if req.DueDate.Before(req.BillDate) {
		return nil, errors.New("due date must be after bill date")
	}

	bill := &VendorBill{
		ID:               uuid.New(),
		OrganizationID:   orgID,
		BillNumber:       req.BillNumber,
		VendorBillNumber: req.VendorBillNumber,
		SupplierID:       req.SupplierID,
		BillDate:         req.BillDate,
		DueDate:          req.DueDate,
		PaymentTerms:     req.PaymentTerms,
		AccountingPeriodID: req.AccountingPeriodID,
		Subtotal:         req.Subtotal,
		TaxAmount:        req.TaxAmount,
		TotalAmount:      req.TotalAmount,
		BalanceDue:       req.TotalAmount,
		Status:           VendorBillStatusUnpaid,
		Description:      req.Description,
		Notes:            req.Notes,
		Memo:             req.Memo,
		Metadata:         req.Metadata,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		CreatedBy:        &userID,
		UpdatedBy:        &userID,
	}

	if err := s.repo.CreateVendorBill(ctx, bill); err != nil {
		return nil, fmt.Errorf("failed to create bill: %w", err)
	}

	lines := make([]VendorBillLine, 0, len(req.Lines))
	for _, lineReq := range req.Lines {
		line := &VendorBillLine{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			VendorBillID:     bill.ID,
			LineNumber:       lineReq.LineNumber,
			ExpenseAccountID: lineReq.ExpenseAccountID,
			Description:      lineReq.Description,
			Quantity:         lineReq.Quantity,
			UnitPrice:        lineReq.UnitPrice,
			Amount:           lineReq.Amount,
			LocationID:       lineReq.LocationID,
			Department:       lineReq.Department,
			ProjectCode:      lineReq.ProjectCode,
			TaxCode:          lineReq.TaxCode,
			TaxAmount:        lineReq.TaxAmount,
			ProductID:        lineReq.ProductID,
			Metadata:         lineReq.Metadata,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := s.repo.CreateVendorBillLine(ctx, line); err != nil {
			return nil, fmt.Errorf("failed to create bill line: %w", err)
		}
		lines = append(lines, *line)
	}

	return &VendorBillWithLines{Bill: bill, Lines: lines}, nil
}

// GetVendorBill retrieves a vendor bill with its line items
func (s *Service) GetVendorBill(ctx context.Context, organizationID, billID uuid.UUID) (*VendorBillWithLines, error) {
	bill, err := s.repo.GetVendorBillByID(ctx, organizationID, billID)
	if err != nil {
		return nil, err
	}

	lines, err := s.repo.ListVendorBillLines(ctx, organizationID, billID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bill lines: %w", err)
	}

	return &VendorBillWithLines{Bill: bill, Lines: lines}, nil
}

// ListVendorBills retrieves a list of vendor bills with filtering
func (s *Service) ListVendorBills(ctx context.Context, organizationID uuid.UUID, opts *VendorBillFilterOptions) ([]VendorBill, int64, error) {
	if opts == nil {
		opts = &VendorBillFilterOptions{Limit: 10, Offset: 0}
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	return s.repo.ListVendorBills(ctx, organizationID, opts)
}

// UpdateVendorBill updates a vendor bill and its line items
func (s *Service) UpdateVendorBill(ctx context.Context, billID uuid.UUID, req *UpdateVendorBillRequest, orgID, userID uuid.UUID) (*VendorBillWithLines, error) {
	bill, err := s.repo.GetVendorBillByID(ctx, orgID, billID)
	if err != nil {
		return nil, err
	}

	// Update simple fields
	if req.VendorBillNumber != nil {
		bill.VendorBillNumber = req.VendorBillNumber
	}
	if req.BillDate != nil {
		bill.BillDate = *req.BillDate
	}
	if req.DueDate != nil {
		if req.DueDate.Before(bill.BillDate) {
			return nil, errors.New("due date must be after bill date")
		}
		bill.DueDate = *req.DueDate
	}
	if req.PaymentTerms != nil {
		bill.PaymentTerms = req.PaymentTerms
	}
	if req.Subtotal != nil {
		bill.Subtotal = *req.Subtotal
	}
	if req.TaxAmount != nil {
		bill.TaxAmount = *req.TaxAmount
	}
	if req.TotalAmount != nil {
		bill.TotalAmount = *req.TotalAmount
		bill.BalanceDue = *req.TotalAmount - bill.PaidAmount
	}
	if req.Description != nil {
		bill.Description = req.Description
	}
	if req.Notes != nil {
		bill.Notes = req.Notes
	}
	if req.Memo != nil {
		bill.Memo = req.Memo
	}
	if req.Status != nil {
		bill.Status = *req.Status
	}
	if req.Metadata != nil {
		bill.Metadata = req.Metadata
	}

	bill.UpdatedAt = time.Now()
	bill.UpdatedBy = &userID

	if err := s.repo.UpdateVendorBill(ctx, bill); err != nil {
		return nil, fmt.Errorf("failed to update bill: %w", err)
	}

	// Update line items
	if req.Lines != nil && len(req.Lines) > 0 {
		if err := s.repo.DeleteVendorBillLines(ctx, orgID, billID); err != nil {
			return nil, fmt.Errorf("failed to delete old lines: %w", err)
		}

		for _, lineReq := range req.Lines {
			line := &VendorBillLine{
				ID:               uuid.New(),
				OrganizationID:   orgID,
				VendorBillID:     billID,
				LineNumber:       *lineReq.LineNumber,
				ExpenseAccountID: uuid.Nil, // This would need to come from req
				Description:      *lineReq.Description,
				Quantity:         *lineReq.Quantity,
				UnitPrice:        *lineReq.UnitPrice,
				Amount:           *lineReq.Amount,
				LocationID:       lineReq.LocationID,
				Department:       lineReq.Department,
				ProjectCode:      lineReq.ProjectCode,
				TaxCode:          lineReq.TaxCode,
				TaxAmount:        *lineReq.TaxAmount,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			if err := s.repo.CreateVendorBillLine(ctx, line); err != nil {
				return nil, fmt.Errorf("failed to create bill line: %w", err)
			}
		}
	}

	return s.GetVendorBill(ctx, orgID, billID)
}

// DeleteVendorBill soft deletes a vendor bill
func (s *Service) DeleteVendorBill(ctx context.Context, organizationID, billID uuid.UUID) error {
	bill, err := s.repo.GetVendorBillByID(ctx, organizationID, billID)
	if err != nil {
		return err
	}

	if bill.IsPosted {
		return errors.New("cannot delete a posted bill")
	}

	return s.repo.DeleteVendorBill(ctx, organizationID, billID)
}

// CreateVendorPayment creates a vendor payment with applications to bills
func (s *Service) CreateVendorPayment(ctx context.Context, req *CreateVendorPaymentRequest, orgID, userID uuid.UUID) (*VendorPaymentWithApplications, error) {
	if req.PaymentAmount <= 0 {
		return nil, errors.New("payment amount must be greater than zero")
	}

	payment := &VendorPayment{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		PaymentNumber:      req.PaymentNumber,
		SupplierID:         req.SupplierID,
		PaymentDate:        req.PaymentDate,
		PaymentMethod:      req.PaymentMethod,
		ReferenceNumber:    req.ReferenceNumber,
		PaymentAmount:      req.PaymentAmount,
		BankAccountID:      req.BankAccountID,
		AccountingPeriodID: req.AccountingPeriodID,
		Memo:               req.Memo,
		Notes:              req.Notes,
		Metadata:           req.Metadata,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          &userID,
		UpdatedBy:          &userID,
	}

	if err := s.repo.CreateVendorPayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	applications := make([]VendorPaymentApplication, 0, len(req.Applications))
	totalApplied := 0.0

	for _, appReq := range req.Applications {
		// Verify bill exists and belongs to the same supplier
		bill, err := s.repo.GetVendorBillByID(ctx, orgID, appReq.VendorBillID)
		if err != nil {
			return nil, fmt.Errorf("invalid bill ID: %w", err)
		}

		if bill.SupplierID != req.SupplierID {
			return nil, errors.New("bill does not belong to the same supplier")
		}

		if appReq.AppliedAmount <= 0 {
			return nil, errors.New("applied amount must be greater than zero")
		}

		if appReq.AppliedAmount > bill.BalanceDue {
			return nil, errors.New("applied amount exceeds bill balance")
		}

		app := &VendorPaymentApplication{
			ID:              uuid.New(),
			OrganizationID:  orgID,
			VendorPaymentID: payment.ID,
			VendorBillID:    appReq.VendorBillID,
			AppliedAmount:   appReq.AppliedAmount,
			CreatedAt:       time.Now(),
		}

		if err := s.repo.CreateVendorPaymentApplication(ctx, app); err != nil {
			return nil, fmt.Errorf("failed to create payment application: %w", err)
		}

		// Update bill balance
		bill.PaidAmount += appReq.AppliedAmount
		bill.BalanceDue = bill.TotalAmount - bill.PaidAmount
		if bill.BalanceDue <= 0 {
			bill.Status = VendorBillStatusPaid
		} else if bill.PaidAmount > 0 {
			bill.Status = VendorBillStatusPartial
		}

		if err := s.repo.UpdateVendorBill(ctx, bill); err != nil {
			return nil, fmt.Errorf("failed to update bill: %w", err)
		}

		applications = append(applications, *app)
		totalApplied += appReq.AppliedAmount
	}

	if totalApplied != payment.PaymentAmount {
		// This is a warning, not necessarily an error - payment might be unapplied
		// In production, you might want to track unapplied cash
	}

	return &VendorPaymentWithApplications{Payment: payment, Applications: applications}, nil
}

// GetVendorPayment retrieves a vendor payment with its applications
func (s *Service) GetVendorPayment(ctx context.Context, organizationID, paymentID uuid.UUID) (*VendorPaymentWithApplications, error) {
	payment, err := s.repo.GetVendorPaymentByID(ctx, organizationID, paymentID)
	if err != nil {
		return nil, err
	}

	applications, err := s.repo.ListVendorPaymentApplications(ctx, organizationID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment applications: %w", err)
	}

	return &VendorPaymentWithApplications{Payment: payment, Applications: applications}, nil
}

// ListVendorPayments retrieves a list of vendor payments with filtering
func (s *Service) ListVendorPayments(ctx context.Context, organizationID uuid.UUID, opts *VendorPaymentFilterOptions) ([]VendorPayment, int64, error) {
	if opts == nil {
		opts = &VendorPaymentFilterOptions{Limit: 10, Offset: 0}
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	return s.repo.ListVendorPayments(ctx, organizationID, opts)
}

// UpdateVendorPayment updates a vendor payment
func (s *Service) UpdateVendorPayment(ctx context.Context, paymentID uuid.UUID, req *UpdateVendorPaymentRequest, orgID, userID uuid.UUID) (*VendorPaymentWithApplications, error) {
	payment, err := s.repo.GetVendorPaymentByID(ctx, orgID, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.IsPosted {
		return nil, errors.New("cannot update a posted payment")
	}

	// Update fields
	if req.PaymentDate != nil {
		payment.PaymentDate = *req.PaymentDate
	}
	if req.PaymentMethod != nil {
		payment.PaymentMethod = *req.PaymentMethod
	}
	if req.ReferenceNumber != nil {
		payment.ReferenceNumber = req.ReferenceNumber
	}
	if req.PaymentAmount != nil && *req.PaymentAmount > 0 {
		payment.PaymentAmount = *req.PaymentAmount
	}
	if req.BankAccountID != nil {
		payment.BankAccountID = req.BankAccountID
	}
	if req.Memo != nil {
		payment.Memo = req.Memo
	}
	if req.Notes != nil {
		payment.Notes = req.Notes
	}
	if req.Metadata != nil {
		payment.Metadata = req.Metadata
	}

	payment.UpdatedAt = time.Now()
	payment.UpdatedBy = &userID

	if err := s.repo.UpdateVendorPayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return s.GetVendorPayment(ctx, orgID, paymentID)
}

// DeleteVendorPayment soft deletes a vendor payment
func (s *Service) DeleteVendorPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error {
	payment, err := s.repo.GetVendorPaymentByID(ctx, organizationID, paymentID)
	if err != nil {
		return err
	}

	if payment.IsPosted {
		return errors.New("cannot delete a posted payment")
	}

	return s.repo.DeleteVendorPayment(ctx, organizationID, paymentID)
}

// GetVendorAgingReport generates an aging report for AP
func (s *Service) GetVendorAgingReport(ctx context.Context, organizationID uuid.UUID) ([]VendorAgingReport, error) {
	return s.repo.GetVendorAgingReport(ctx, organizationID)
}

// GetSupplierBalance retrieves the total balance owed to a supplier
func (s *Service) GetSupplierBalance(ctx context.Context, organizationID, supplierID uuid.UUID) (float64, error) {
	return s.repo.GetSupplierBalance(ctx, organizationID, supplierID)
}

// GetOverdueBills retrieves all overdue bills
func (s *Service) GetOverdueBills(ctx context.Context, organizationID uuid.UUID) ([]VendorBill, error) {
	today := time.Now()
	return s.repo.GetBillsByDueDate(ctx, organizationID, today)
}

// CheckBillExists checks if a bill number already exists
func (s *Service) CheckBillExists(ctx context.Context, organizationID uuid.UUID, billNumber string) (bool, error) {
	bill, err := s.repo.GetVendorBillByNumber(ctx, organizationID, billNumber)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return bill != nil, nil
}

// CheckPaymentExists checks if a payment number already exists
func (s *Service) CheckPaymentExists(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (bool, error) {
	payment, err := s.repo.GetVendorPaymentByNumber(ctx, organizationID, paymentNumber)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return payment != nil, nil
}

// RemovePaymentApplication removes the application of a payment to a bill
func (s *Service) RemovePaymentApplication(ctx context.Context, orgID, appID uuid.UUID, userID uuid.UUID) error {
	app, err := s.repo.GetVendorPaymentApplicationByID(ctx, orgID, appID)
	if err != nil {
		return err
	}

	// Get the bill and reverse the payment
	bill, err := s.repo.GetVendorBillByID(ctx, orgID, app.VendorBillID)
	if err != nil {
		return err
	}

	bill.PaidAmount -= app.AppliedAmount
	bill.BalanceDue = bill.TotalAmount - bill.PaidAmount

	if bill.BalanceDue >= bill.TotalAmount {
		bill.Status = VendorBillStatusUnpaid
	} else if bill.PaidAmount > 0 {
		bill.Status = VendorBillStatusPartial
	}

	bill.UpdatedAt = time.Now()
	bill.UpdatedBy = &userID

	if err := s.repo.UpdateVendorBill(ctx, bill); err != nil {
		return fmt.Errorf("failed to update bill: %w", err)
	}

	return s.repo.DeleteVendorPaymentApplication(ctx, orgID, appID)
}
