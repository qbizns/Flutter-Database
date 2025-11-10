package receivables

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for receivables data access
type Repository interface {
	// Customer Invoices
	CreateCustomerInvoice(ctx context.Context, invoice *CustomerInvoice) error
	GetCustomerInvoiceByID(ctx context.Context, organizationID, invoiceID uuid.UUID) (*CustomerInvoice, error)
	GetCustomerInvoiceByNumber(ctx context.Context, organizationID uuid.UUID, invoiceNumber string) (*CustomerInvoice, error)
	ListCustomerInvoices(ctx context.Context, organizationID uuid.UUID, opts *CustomerInvoiceFilterOptions) ([]CustomerInvoice, int64, error)
	UpdateCustomerInvoice(ctx context.Context, invoice *CustomerInvoice) error
	DeleteCustomerInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) error

	// Customer Invoice Lines
	CreateCustomerInvoiceLine(ctx context.Context, line *CustomerInvoiceLine) error
	GetCustomerInvoiceLineByID(ctx context.Context, organizationID, lineID uuid.UUID) (*CustomerInvoiceLine, error)
	ListCustomerInvoiceLines(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerInvoiceLine, error)
	UpdateCustomerInvoiceLine(ctx context.Context, line *CustomerInvoiceLine) error
	DeleteCustomerInvoiceLine(ctx context.Context, organizationID, lineID uuid.UUID) error
	DeleteCustomerInvoiceLines(ctx context.Context, organizationID, invoiceID uuid.UUID) error

	// Customer Payments
	CreateCustomerPayment(ctx context.Context, payment *CustomerPayment) error
	GetCustomerPaymentByID(ctx context.Context, organizationID, paymentID uuid.UUID) (*CustomerPayment, error)
	GetCustomerPaymentByNumber(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*CustomerPayment, error)
	ListCustomerPayments(ctx context.Context, organizationID uuid.UUID, opts *CustomerPaymentFilterOptions) ([]CustomerPayment, int64, error)
	UpdateCustomerPayment(ctx context.Context, payment *CustomerPayment) error
	DeleteCustomerPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error

	// Customer Payment Applications
	CreateCustomerPaymentApplication(ctx context.Context, app *CustomerPaymentApplication) error
	GetCustomerPaymentApplicationByID(ctx context.Context, organizationID, appID uuid.UUID) (*CustomerPaymentApplication, error)
	ListCustomerPaymentApplications(ctx context.Context, organizationID, paymentID uuid.UUID) ([]CustomerPaymentApplication, error)
	ListPaymentApplicationsByInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]CustomerPaymentApplication, error)
	DeleteCustomerPaymentApplication(ctx context.Context, organizationID, appID uuid.UUID) error

	// Reporting
	GetCustomerAgingReport(ctx context.Context, organizationID uuid.UUID) ([]CustomerAgingReport, error)
	GetInvoicesByDueDate(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]CustomerInvoice, error)
	GetCustomerBalance(ctx context.Context, organizationID, customerID uuid.UUID) (float64, error)
}

// Service provides business logic for Accounts Receivable
type Service struct {
	repo Repository
}

// NewService creates a new receivables service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateCustomerInvoice creates a new customer invoice with line items
func (s *Service) CreateCustomerInvoice(ctx context.Context, req *CreateCustomerInvoiceRequest, orgID, userID uuid.UUID) (*CustomerInvoiceWithLines, error) {
	if req.TotalAmount < 0 {
		return nil, errors.New("total amount must be non-negative")
	}

	if req.DueDate.Before(req.InvoiceDate) {
		return nil, errors.New("due date must be after invoice date")
	}

	invoice := &CustomerInvoice{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		InvoiceNumber:      req.InvoiceNumber,
		CustomerID:         req.CustomerID,
		InvoiceDate:        req.InvoiceDate,
		DueDate:            req.DueDate,
		PaymentTerms:       req.PaymentTerms,
		AccountingPeriodID: req.AccountingPeriodID,
		Subtotal:           req.Subtotal,
		TaxAmount:          req.TaxAmount,
		DiscountAmount:     req.DiscountAmount,
		TotalAmount:        req.TotalAmount,
		BalanceDue:         req.TotalAmount,
		Status:             InvoiceStatusUnpaid,
		SaleID:             req.SaleID,
		Description:        req.Description,
		Notes:              req.Notes,
		Memo:               req.Memo,
		Metadata:           req.Metadata,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          &userID,
		UpdatedBy:          &userID,
	}

	if err := s.repo.CreateCustomerInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	lines := make([]CustomerInvoiceLine, 0, len(req.Lines))
	for _, lineReq := range req.Lines {
		line := &CustomerInvoiceLine{
			ID:                 uuid.New(),
			OrganizationID:     orgID,
			CustomerInvoiceID:  invoice.ID,
			LineNumber:         lineReq.LineNumber,
			RevenueAccountID:   lineReq.RevenueAccountID,
			Description:        lineReq.Description,
			Quantity:           lineReq.Quantity,
			UnitPrice:          lineReq.UnitPrice,
			Amount:             lineReq.Amount,
			LocationID:         lineReq.LocationID,
			Department:         lineReq.Department,
			ProjectCode:        lineReq.ProjectCode,
			TaxCode:            lineReq.TaxCode,
			TaxAmount:          lineReq.TaxAmount,
			ProductID:          lineReq.ProductID,
			Metadata:           lineReq.Metadata,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		if err := s.repo.CreateCustomerInvoiceLine(ctx, line); err != nil {
			return nil, fmt.Errorf("failed to create invoice line: %w", err)
		}
		lines = append(lines, *line)
	}

	return &CustomerInvoiceWithLines{Invoice: invoice, Lines: lines}, nil
}

// GetCustomerInvoice retrieves a customer invoice with its line items
func (s *Service) GetCustomerInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) (*CustomerInvoiceWithLines, error) {
	invoice, err := s.repo.GetCustomerInvoiceByID(ctx, organizationID, invoiceID)
	if err != nil {
		return nil, err
	}

	lines, err := s.repo.ListCustomerInvoiceLines(ctx, organizationID, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice lines: %w", err)
	}

	return &CustomerInvoiceWithLines{Invoice: invoice, Lines: lines}, nil
}

// ListCustomerInvoices retrieves a list of customer invoices with filtering
func (s *Service) ListCustomerInvoices(ctx context.Context, organizationID uuid.UUID, opts *CustomerInvoiceFilterOptions) ([]CustomerInvoice, int64, error) {
	if opts == nil {
		opts = &CustomerInvoiceFilterOptions{Limit: 10, Offset: 0}
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	return s.repo.ListCustomerInvoices(ctx, organizationID, opts)
}

// UpdateCustomerInvoice updates a customer invoice and its line items
func (s *Service) UpdateCustomerInvoice(ctx context.Context, invoiceID uuid.UUID, req *UpdateCustomerInvoiceRequest, orgID, userID uuid.UUID) (*CustomerInvoiceWithLines, error) {
	invoice, err := s.repo.GetCustomerInvoiceByID(ctx, orgID, invoiceID)
	if err != nil {
		return nil, err
	}

	// Update simple fields
	if req.InvoiceDate != nil {
		invoice.InvoiceDate = *req.InvoiceDate
	}
	if req.DueDate != nil {
		if req.DueDate.Before(invoice.InvoiceDate) {
			return nil, errors.New("due date must be after invoice date")
		}
		invoice.DueDate = *req.DueDate
	}
	if req.PaymentTerms != nil {
		invoice.PaymentTerms = req.PaymentTerms
	}
	if req.Subtotal != nil {
		invoice.Subtotal = *req.Subtotal
	}
	if req.TaxAmount != nil {
		invoice.TaxAmount = *req.TaxAmount
	}
	if req.DiscountAmount != nil {
		invoice.DiscountAmount = *req.DiscountAmount
	}
	if req.TotalAmount != nil {
		invoice.TotalAmount = *req.TotalAmount
		invoice.BalanceDue = *req.TotalAmount - invoice.PaidAmount
	}
	if req.Description != nil {
		invoice.Description = req.Description
	}
	if req.Notes != nil {
		invoice.Notes = req.Notes
	}
	if req.Memo != nil {
		invoice.Memo = req.Memo
	}
	if req.Status != nil {
		invoice.Status = *req.Status
	}
	if req.Metadata != nil {
		invoice.Metadata = req.Metadata
	}

	invoice.UpdatedAt = time.Now()
	invoice.UpdatedBy = &userID

	if err := s.repo.UpdateCustomerInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	// Update line items
	if req.Lines != nil && len(req.Lines) > 0 {
		if err := s.repo.DeleteCustomerInvoiceLines(ctx, orgID, invoiceID); err != nil {
			return nil, fmt.Errorf("failed to delete old lines: %w", err)
		}

		for _, lineReq := range req.Lines {
			line := &CustomerInvoiceLine{
				ID:                uuid.New(),
				OrganizationID:    orgID,
				CustomerInvoiceID: invoiceID,
				LineNumber:        *lineReq.LineNumber,
				RevenueAccountID:  uuid.Nil, // This would need to come from req
				Description:       *lineReq.Description,
				Quantity:          *lineReq.Quantity,
				UnitPrice:         *lineReq.UnitPrice,
				Amount:            *lineReq.Amount,
				LocationID:        lineReq.LocationID,
				Department:        lineReq.Department,
				ProjectCode:       lineReq.ProjectCode,
				TaxCode:           lineReq.TaxCode,
				TaxAmount:         *lineReq.TaxAmount,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			if err := s.repo.CreateCustomerInvoiceLine(ctx, line); err != nil {
				return nil, fmt.Errorf("failed to create invoice line: %w", err)
			}
		}
	}

	return s.GetCustomerInvoice(ctx, orgID, invoiceID)
}

// DeleteCustomerInvoice soft deletes a customer invoice
func (s *Service) DeleteCustomerInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) error {
	invoice, err := s.repo.GetCustomerInvoiceByID(ctx, organizationID, invoiceID)
	if err != nil {
		return err
	}

	if invoice.IsPosted {
		return errors.New("cannot delete a posted invoice")
	}

	return s.repo.DeleteCustomerInvoice(ctx, organizationID, invoiceID)
}

// CreateCustomerPayment creates a customer payment with applications to invoices
func (s *Service) CreateCustomerPayment(ctx context.Context, req *CreateCustomerPaymentRequest, orgID, userID uuid.UUID) (*CustomerPaymentWithApplications, error) {
	if req.PaymentAmount <= 0 {
		return nil, errors.New("payment amount must be greater than zero")
	}

	payment := &CustomerPayment{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		PaymentNumber:      req.PaymentNumber,
		CustomerID:         req.CustomerID,
		PaymentDate:        req.PaymentDate,
		PaymentMethod:      req.PaymentMethod,
		ReferenceNumber:    req.ReferenceNumber,
		PaymentAmount:      req.PaymentAmount,
		DepositAccountID:   req.DepositAccountID,
		AccountingPeriodID: req.AccountingPeriodID,
		Memo:               req.Memo,
		Notes:              req.Notes,
		Metadata:           req.Metadata,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          &userID,
		UpdatedBy:          &userID,
	}

	if err := s.repo.CreateCustomerPayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	applications := make([]CustomerPaymentApplication, 0, len(req.Applications))
	totalApplied := 0.0

	for _, appReq := range req.Applications {
		// Verify invoice exists and belongs to the same customer
		invoice, err := s.repo.GetCustomerInvoiceByID(ctx, orgID, appReq.CustomerInvoiceID)
		if err != nil {
			return nil, fmt.Errorf("invalid invoice ID: %w", err)
		}

		if invoice.CustomerID != req.CustomerID {
			return nil, errors.New("invoice does not belong to the same customer")
		}

		if appReq.AppliedAmount <= 0 {
			return nil, errors.New("applied amount must be greater than zero")
		}

		if appReq.AppliedAmount > invoice.BalanceDue {
			return nil, errors.New("applied amount exceeds invoice balance")
		}

		app := &CustomerPaymentApplication{
			ID:                uuid.New(),
			OrganizationID:    orgID,
			CustomerPaymentID: payment.ID,
			CustomerInvoiceID: appReq.CustomerInvoiceID,
			AppliedAmount:     appReq.AppliedAmount,
			CreatedAt:         time.Now(),
		}

		if err := s.repo.CreateCustomerPaymentApplication(ctx, app); err != nil {
			return nil, fmt.Errorf("failed to create payment application: %w", err)
		}

		// Update invoice balance
		invoice.PaidAmount += appReq.AppliedAmount
		invoice.BalanceDue = invoice.TotalAmount - invoice.PaidAmount
		if invoice.BalanceDue <= 0 {
			invoice.Status = InvoiceStatusPaid
		} else if invoice.PaidAmount > 0 {
			invoice.Status = InvoiceStatusPartial
		}

		if err := s.repo.UpdateCustomerInvoice(ctx, invoice); err != nil {
			return nil, fmt.Errorf("failed to update invoice: %w", err)
		}

		applications = append(applications, *app)
		totalApplied += appReq.AppliedAmount
	}

	if totalApplied != payment.PaymentAmount {
		// This is a warning, not necessarily an error - payment might be unapplied
		// In production, you might want to track unapplied cash
	}

	return &CustomerPaymentWithApplications{Payment: payment, Applications: applications}, nil
}

// GetCustomerPayment retrieves a customer payment with its applications
func (s *Service) GetCustomerPayment(ctx context.Context, organizationID, paymentID uuid.UUID) (*CustomerPaymentWithApplications, error) {
	payment, err := s.repo.GetCustomerPaymentByID(ctx, organizationID, paymentID)
	if err != nil {
		return nil, err
	}

	applications, err := s.repo.ListCustomerPaymentApplications(ctx, organizationID, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment applications: %w", err)
	}

	return &CustomerPaymentWithApplications{Payment: payment, Applications: applications}, nil
}

// ListCustomerPayments retrieves a list of customer payments with filtering
func (s *Service) ListCustomerPayments(ctx context.Context, organizationID uuid.UUID, opts *CustomerPaymentFilterOptions) ([]CustomerPayment, int64, error) {
	if opts == nil {
		opts = &CustomerPaymentFilterOptions{Limit: 10, Offset: 0}
	}

	if opts.Limit == 0 {
		opts.Limit = 10
	}

	return s.repo.ListCustomerPayments(ctx, organizationID, opts)
}

// UpdateCustomerPayment updates a customer payment
func (s *Service) UpdateCustomerPayment(ctx context.Context, paymentID uuid.UUID, req *UpdateCustomerPaymentRequest, orgID, userID uuid.UUID) (*CustomerPaymentWithApplications, error) {
	payment, err := s.repo.GetCustomerPaymentByID(ctx, orgID, paymentID)
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
	if req.DepositAccountID != nil {
		payment.DepositAccountID = req.DepositAccountID
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

	if err := s.repo.UpdateCustomerPayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return s.GetCustomerPayment(ctx, orgID, paymentID)
}

// DeleteCustomerPayment soft deletes a customer payment
func (s *Service) DeleteCustomerPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error {
	payment, err := s.repo.GetCustomerPaymentByID(ctx, organizationID, paymentID)
	if err != nil {
		return err
	}

	if payment.IsPosted {
		return errors.New("cannot delete a posted payment")
	}

	return s.repo.DeleteCustomerPayment(ctx, organizationID, paymentID)
}

// GetCustomerAgingReport generates an aging report for AR
func (s *Service) GetCustomerAgingReport(ctx context.Context, organizationID uuid.UUID) ([]CustomerAgingReport, error) {
	return s.repo.GetCustomerAgingReport(ctx, organizationID)
}

// GetCustomerBalance retrieves the total balance owed by a customer
func (s *Service) GetCustomerBalance(ctx context.Context, organizationID, customerID uuid.UUID) (float64, error) {
	return s.repo.GetCustomerBalance(ctx, organizationID, customerID)
}

// GetOverdueInvoices retrieves all overdue invoices
func (s *Service) GetOverdueInvoices(ctx context.Context, organizationID uuid.UUID) ([]CustomerInvoice, error) {
	today := time.Now()
	return s.repo.GetInvoicesByDueDate(ctx, organizationID, today)
}

// CheckInvoiceExists checks if an invoice number already exists
func (s *Service) CheckInvoiceExists(ctx context.Context, organizationID uuid.UUID, invoiceNumber string) (bool, error) {
	invoice, err := s.repo.GetCustomerInvoiceByNumber(ctx, organizationID, invoiceNumber)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return invoice != nil, nil
}

// CheckPaymentExists checks if a payment number already exists
func (s *Service) CheckPaymentExists(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (bool, error) {
	payment, err := s.repo.GetCustomerPaymentByNumber(ctx, organizationID, paymentNumber)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return payment != nil, nil
}

// RemovePaymentApplication removes the application of a payment to an invoice
func (s *Service) RemovePaymentApplication(ctx context.Context, orgID, appID uuid.UUID, userID uuid.UUID) error {
	app, err := s.repo.GetCustomerPaymentApplicationByID(ctx, orgID, appID)
	if err != nil {
		return err
	}

	// Get the invoice and reverse the payment
	invoice, err := s.repo.GetCustomerInvoiceByID(ctx, orgID, app.CustomerInvoiceID)
	if err != nil {
		return err
	}

	invoice.PaidAmount -= app.AppliedAmount
	invoice.BalanceDue = invoice.TotalAmount - invoice.PaidAmount

	if invoice.BalanceDue >= invoice.TotalAmount {
		invoice.Status = InvoiceStatusUnpaid
	} else if invoice.PaidAmount > 0 {
		invoice.Status = InvoiceStatusPartial
	}

	invoice.UpdatedAt = time.Now()
	invoice.UpdatedBy = &userID

	if err := s.repo.UpdateCustomerInvoice(ctx, invoice); err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return s.repo.DeleteCustomerPaymentApplication(ctx, orgID, appID)
}
