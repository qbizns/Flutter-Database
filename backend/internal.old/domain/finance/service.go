package finance

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// FinanceRepository defines the interface for finance database operations
type FinanceRepository interface {
	// Currencies
	CreateCurrency(ctx context.Context, currency *Currency) error
	GetCurrency(ctx context.Context, id uuid.UUID) (*Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*Currency, error)
	ListCurrencies(ctx context.Context, filter map[string]interface{}) ([]*Currency, error)
	UpdateCurrency(ctx context.Context, currency *Currency) error
	DeleteCurrency(ctx context.Context, id uuid.UUID) error

	// Currency Rates
	CreateCurrencyRate(ctx context.Context, rate *CurrencyRate) error
	GetCurrencyRate(ctx context.Context, id uuid.UUID) (*CurrencyRate, error)
	GetCurrencyRateByDate(ctx context.Context, organizationID uuid.UUID, currencyCode string, rateDate time.Time) (*CurrencyRate, error)
	ListCurrencyRates(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*CurrencyRate, error)
	UpdateCurrencyRate(ctx context.Context, rate *CurrencyRate) error
	DeleteCurrencyRate(ctx context.Context, id uuid.UUID) error
	GetLatestCurrencyRate(ctx context.Context, organizationID uuid.UUID, currencyCode string) (*CurrencyRate, error)

	// Payment Terms
	CreatePaymentTerm(ctx context.Context, term *PaymentTerm) error
	GetPaymentTerm(ctx context.Context, id, organizationID uuid.UUID) (*PaymentTerm, error)
	GetPaymentTermByCode(ctx context.Context, organizationID uuid.UUID, code string) (*PaymentTerm, error)
	ListPaymentTerms(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*PaymentTerm, error)
	UpdatePaymentTerm(ctx context.Context, term *PaymentTerm) error
	DeletePaymentTerm(ctx context.Context, id, organizationID uuid.UUID) error

	// Payment Term Lines
	CreatePaymentTermLine(ctx context.Context, line *PaymentTermLine) error
	GetPaymentTermLine(ctx context.Context, id uuid.UUID) (*PaymentTermLine, error)
	ListPaymentTermLines(ctx context.Context, paymentTermID uuid.UUID) ([]*PaymentTermLine, error)
	DeletePaymentTermLine(ctx context.Context, id uuid.UUID) error

	// Invoice Payment Schedules
	CreateInvoicePaymentSchedules(ctx context.Context, schedules []*InvoicePaymentSchedule) error
	GetInvoicePaymentSchedule(ctx context.Context, id uuid.UUID) (*InvoicePaymentSchedule, error)
	ListInvoicePaymentSchedules(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*InvoicePaymentSchedule, error)
	ListInvoicePaymentSchedulesBySource(ctx context.Context, organizationID uuid.UUID, sourceType SourceType, sourceID uuid.UUID) ([]*InvoicePaymentSchedule, error)
	UpdateInvoicePaymentSchedule(ctx context.Context, schedule *InvoicePaymentSchedule) error
	DeleteInvoicePaymentSchedule(ctx context.Context, id uuid.UUID) error
	GetInvoicePaymentScheduleBySourceAndLine(ctx context.Context, organizationID uuid.UUID, sourceType SourceType, sourceID uuid.UUID, lineNumber int) (*InvoicePaymentSchedule, error)
}

// Service handles finance business logic
type Service struct {
	repo FinanceRepository
}

// NewService creates a new finance service
func NewService(repo FinanceRepository) *Service {
	return &Service{repo: repo}
}

// ========================
// CURRENCY MANAGEMENT
// ========================

// CreateCurrency creates a new currency
func (s *Service) CreateCurrency(ctx context.Context, req *CreateCurrencyRequest) (*Currency, error) {
	if len(req.CurrencyCode) != 3 {
		return nil, ErrInvalidCurrencyCode
	}

	currency := &Currency{
		ID:             uuid.New(),
		CurrencyCode:   req.CurrencyCode,
		CurrencyName:   req.CurrencyName,
		CurrencySymbol: req.CurrencySymbol,
		DecimalPlaces:  req.DecimalPlaces,
		IsActive:       req.IsActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateCurrency(ctx, currency); err != nil {
		return nil, err
	}

	return currency, nil
}

// GetCurrency retrieves a currency by ID
func (s *Service) GetCurrency(ctx context.Context, id uuid.UUID) (*Currency, error) {
	return s.repo.GetCurrency(ctx, id)
}

// GetCurrencyByCode retrieves a currency by code
func (s *Service) GetCurrencyByCode(ctx context.Context, code string) (*Currency, error) {
	return s.repo.GetCurrencyByCode(ctx, code)
}

// UpdateCurrency updates a currency
func (s *Service) UpdateCurrency(ctx context.Context, id uuid.UUID, req *UpdateCurrencyRequest) (*Currency, error) {
	currency, err := s.repo.GetCurrency(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CurrencyName != nil {
		currency.CurrencyName = *req.CurrencyName
	}
	if req.CurrencySymbol != nil {
		currency.CurrencySymbol = req.CurrencySymbol
	}
	if req.DecimalPlaces != nil {
		currency.DecimalPlaces = *req.DecimalPlaces
	}
	if req.IsActive != nil {
		currency.IsActive = *req.IsActive
	}

	currency.UpdatedAt = time.Now()

	if err := s.repo.UpdateCurrency(ctx, currency); err != nil {
		return nil, err
	}

	return currency, nil
}

// ListCurrencies lists currencies with optional filters
func (s *Service) ListCurrencies(ctx context.Context, filter map[string]interface{}) ([]*Currency, error) {
	return s.repo.ListCurrencies(ctx, filter)
}

// ========================
// CURRENCY RATE MANAGEMENT
// ========================

// CreateCurrencyRate creates a new exchange rate
func (s *Service) CreateCurrencyRate(ctx context.Context, organizationID uuid.UUID, req *CreateCurrencyRateRequest, createdBy uuid.UUID) (*CurrencyRate, error) {
	// Validate currency exists
	_, err := s.repo.GetCurrencyByCode(ctx, req.CurrencyCode)
	if err != nil {
		return nil, err
	}

	// Validate rate date is not in the future
	if req.RateDate.After(time.Now().Add(24 * time.Hour)) {
		return nil, ErrInvalidConversionDate
	}

	// Validate rate format
	if _, err := ParseDecimal(req.Rate); err != nil {
		return nil, err
	}

	rate := &CurrencyRate{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		CurrencyCode:   req.CurrencyCode,
		RateDate:       req.RateDate,
		Rate:           req.Rate,
		Source:         req.Source,
		CreatedBy:      &createdBy,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.CreateCurrencyRate(ctx, rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// GetCurrencyRate retrieves an exchange rate by ID
func (s *Service) GetCurrencyRate(ctx context.Context, id uuid.UUID) (*CurrencyRate, error) {
	return s.repo.GetCurrencyRate(ctx, id)
}

// GetCurrencyRateByDate retrieves the exchange rate for a specific date
func (s *Service) GetCurrencyRateByDate(ctx context.Context, organizationID uuid.UUID, currencyCode string, rateDate time.Time) (*CurrencyRate, error) {
	return s.repo.GetCurrencyRateByDate(ctx, organizationID, currencyCode, rateDate)
}

// GetLatestCurrencyRate retrieves the latest exchange rate
func (s *Service) GetLatestCurrencyRate(ctx context.Context, organizationID uuid.UUID, currencyCode string) (*CurrencyRate, error) {
	return s.repo.GetLatestCurrencyRate(ctx, organizationID, currencyCode)
}

// UpdateCurrencyRate updates an exchange rate
func (s *Service) UpdateCurrencyRate(ctx context.Context, id uuid.UUID, req *UpdateCurrencyRateRequest) (*CurrencyRate, error) {
	rate, err := s.repo.GetCurrencyRate(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Rate != nil {
		if _, err := ParseDecimal(*req.Rate); err != nil {
			return nil, err
		}
		rate.Rate = *req.Rate
	}
	if req.Source != nil {
		rate.Source = *req.Source
	}

	if err := s.repo.UpdateCurrencyRate(ctx, rate); err != nil {
		return nil, err
	}

	return rate, nil
}

// ListCurrencyRates lists exchange rates with filters
func (s *Service) ListCurrencyRates(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*CurrencyRate, error) {
	return s.repo.ListCurrencyRates(ctx, organizationID, filter)
}

// ========================
// CURRENCY CONVERSION
// ========================

// ConvertCurrency converts an amount from one currency to another
func (s *Service) ConvertCurrency(ctx context.Context, organizationID uuid.UUID, req *CurrencyConversionRequest) (*CurrencyConversionResult, error) {
	// Validate currencies exist
	fromCurr, err := s.repo.GetCurrencyByCode(ctx, req.FromCurrency)
	if err != nil {
		return nil, fmt.Errorf("from currency not found: %w", err)
	}

	toCurr, err := s.repo.GetCurrencyByCode(ctx, req.ToCurrency)
	if err != nil {
		return nil, fmt.Errorf("to currency not found: %w", err)
	}

	// If same currency, return 1:1
	if req.FromCurrency == req.ToCurrency {
		return &CurrencyConversionResult{
			Amount:          req.Amount,
			FromCurrency:    req.FromCurrency,
			ToCurrency:      req.ToCurrency,
			ConversionDate:  req.ConversionDate,
			FromRate:        "1",
			ToRate:          "1",
			ConvertedAmount: req.Amount,
		}, nil
	}

	// Get exchange rates
	fromRate, err := s.repo.GetCurrencyRateByDate(ctx, organizationID, req.FromCurrency, req.ConversionDate)
	if err != nil {
		return nil, fmt.Errorf("from currency rate not found: %w", err)
	}

	toRate, err := s.repo.GetCurrencyRateByDate(ctx, organizationID, req.ToCurrency, req.ConversionDate)
	if err != nil {
		return nil, fmt.Errorf("to currency rate not found: %w", err)
	}

	// Calculate conversion
	amountDecimal, _ := ParseDecimal(req.Amount)
	fromRateDecimal, _ := ParseDecimal(fromRate.Rate)
	toRateDecimal, _ := ParseDecimal(toRate.Rate)

	// Convert: amount * (toRate / fromRate)
	ratio := new(big.Float).Quo(toRateDecimal, fromRateDecimal)
	convertedAmount := new(big.Float).Mul(amountDecimal, ratio)

	result := &CurrencyConversionResult{
		Amount:          req.Amount,
		FromCurrency:    req.FromCurrency,
		ToCurrency:      req.ToCurrency,
		ConversionDate:  req.ConversionDate,
		FromRate:        fromRate.Rate,
		ToRate:          toRate.Rate,
		ConvertedAmount: FormatDecimal(convertedAmount, toCurr.DecimalPlaces),
	}

	return result, nil
}

// ========================
// PAYMENT TERM MANAGEMENT
// ========================

// CreatePaymentTerm creates a new payment term with lines
func (s *Service) CreatePaymentTerm(ctx context.Context, organizationID uuid.UUID, req *CreatePaymentTermRequest, createdBy uuid.UUID) (*PaymentTerm, error) {
	// Validate at least one line
	if len(req.Lines) == 0 {
		return nil, ErrPaymentTermLineValidation
	}

	now := time.Now()
	term := &PaymentTerm{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		TermCode:       req.TermCode,
		TermName:       req.TermName,
		Note:           req.Note,
		IsActive:       req.IsActive,
		CreatedBy:      &createdBy,
		UpdatedBy:      &createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Create lines
	lines := make([]*PaymentTermLine, 0, len(req.Lines))
	for _, lineReq := range req.Lines {
		line := &PaymentTermLine{
			ID:            uuid.New(),
			PaymentTermID: term.ID,
			Sequence:      lineReq.Sequence,
			ValueType:     lineReq.ValueType,
			ValueAmount:   lineReq.ValueAmount,
			DaysAfter:     lineReq.DaysAfter,
			EndOfMonth:    lineReq.EndOfMonth,
			DayOfMonth:    lineReq.DayOfMonth,
			CreatedAt:     now,
		}

		// Validate line
		if line.ValueType == PaymentTermValueTypeBalance && lineReq.ValueAmount != nil {
			return nil, fmt.Errorf("balance type cannot have value amount")
		}
		if line.ValueType != PaymentTermValueTypeBalance && lineReq.ValueAmount == nil {
			return nil, fmt.Errorf("non-balance type requires value amount")
		}

		lines = append(lines, line)
	}

	if err := s.repo.CreatePaymentTerm(ctx, term); err != nil {
		return nil, err
	}

	// Create lines in separate operation
	for _, line := range lines {
		if err := s.repo.CreatePaymentTermLine(ctx, line); err != nil {
			return nil, err
		}
	}

	term.Lines = lines
	return term, nil
}

// GetPaymentTerm retrieves a payment term with lines
func (s *Service) GetPaymentTerm(ctx context.Context, id, organizationID uuid.UUID) (*PaymentTerm, error) {
	term, err := s.repo.GetPaymentTerm(ctx, id, organizationID)
	if err != nil {
		return nil, err
	}

	// Load lines
	lines, err := s.repo.ListPaymentTermLines(ctx, term.ID)
	if err != nil {
		return nil, err
	}
	term.Lines = lines

	return term, nil
}

// GetPaymentTermByCode retrieves a payment term by code
func (s *Service) GetPaymentTermByCode(ctx context.Context, organizationID uuid.UUID, code string) (*PaymentTerm, error) {
	term, err := s.repo.GetPaymentTermByCode(ctx, organizationID, code)
	if err != nil {
		return nil, err
	}

	// Load lines
	lines, err := s.repo.ListPaymentTermLines(ctx, term.ID)
	if err != nil {
		return nil, err
	}
	term.Lines = lines

	return term, nil
}

// UpdatePaymentTerm updates a payment term
func (s *Service) UpdatePaymentTerm(ctx context.Context, id, organizationID uuid.UUID, req *UpdatePaymentTermRequest, updatedBy uuid.UUID) (*PaymentTerm, error) {
	term, err := s.GetPaymentTerm(ctx, id, organizationID)
	if err != nil {
		return nil, err
	}

	if req.TermName != nil {
		term.TermName = *req.TermName
	}
	if req.Note != nil {
		term.Note = req.Note
	}
	if req.IsActive != nil {
		term.IsActive = *req.IsActive
	}

	term.UpdatedAt = time.Now()
	term.UpdatedBy = &updatedBy

	if err := s.repo.UpdatePaymentTerm(ctx, term); err != nil {
		return nil, err
	}

	// Update lines if provided
	if len(req.Lines) > 0 {
		// Delete existing lines
		existingLines, _ := s.repo.ListPaymentTermLines(ctx, term.ID)
		for _, line := range existingLines {
			s.repo.DeletePaymentTermLine(ctx, line.ID)
		}

		// Create new lines
		newLines := make([]*PaymentTermLine, 0, len(req.Lines))
		for _, lineReq := range req.Lines {
			line := &PaymentTermLine{
				ID:            uuid.New(),
				PaymentTermID: term.ID,
				Sequence:      lineReq.Sequence,
				ValueType:     lineReq.ValueType,
				ValueAmount:   lineReq.ValueAmount,
				DaysAfter:     lineReq.DaysAfter,
				EndOfMonth:    lineReq.EndOfMonth,
				DayOfMonth:    lineReq.DayOfMonth,
				CreatedAt:     time.Now(),
			}
			if err := s.repo.CreatePaymentTermLine(ctx, line); err != nil {
				return nil, err
			}
			newLines = append(newLines, line)
		}
		term.Lines = newLines
	}

	return term, nil
}

// ListPaymentTerms lists payment terms
func (s *Service) ListPaymentTerms(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*PaymentTerm, error) {
	terms, err := s.repo.ListPaymentTerms(ctx, organizationID, filter)
	if err != nil {
		return nil, err
	}

	// Load lines for each term
	for _, term := range terms {
		lines, err := s.repo.ListPaymentTermLines(ctx, term.ID)
		if err != nil {
			continue
		}
		term.Lines = lines
	}

	return terms, nil
}

// ========================
// INVOICE PAYMENT SCHEDULES
// ========================

// CalculatePaymentSchedule calculates payment schedule from a payment term and invoice amount
func (s *Service) CalculatePaymentSchedule(ctx context.Context, organizationID uuid.UUID, paymentTermID uuid.UUID, invoiceAmount string, invoiceDate time.Time) (*PaymentScheduleCalculation, error) {
	term, err := s.GetPaymentTerm(ctx, paymentTermID, organizationID)
	if err != nil {
		return nil, err
	}

	if len(term.Lines) == 0 {
		return nil, ErrPaymentTermLineValidation
	}

	// Sort lines by sequence
	payments := make([]*CalculatedPayment, 0)
	invoiceAmountDecimal, _ := ParseDecimal(invoiceAmount)
	balanceAmount := new(big.Float).Set(invoiceAmountDecimal)

	for _, line := range term.Lines {
		var dueDate time.Time

		// Calculate due date
		if line.EndOfMonth {
			// Last day of month + days_after
			dueDate = invoiceDate.AddDate(0, 1, 0)
			dueDate = time.Date(dueDate.Year(), dueDate.Month(), 1, 0, 0, 0, 0, dueDate.Location()).AddDate(0, 0, -1)
			dueDate = dueDate.AddDate(0, 0, line.DaysAfter)
		} else if line.DayOfMonth != nil {
			// Specific day of month + days_after
			nextMonth := invoiceDate.AddDate(0, 1, 0)
			dueDate = time.Date(nextMonth.Year(), nextMonth.Month(), *line.DayOfMonth, 0, 0, 0, 0, nextMonth.Location())
			dueDate = dueDate.AddDate(0, 0, line.DaysAfter)
		} else {
			// Days after invoice date
			dueDate = invoiceDate.AddDate(0, 0, line.DaysAfter)
		}

		// Calculate amount
		var amountDue *big.Float
		if line.ValueType == PaymentTermValueTypeBalance {
			amountDue = balanceAmount
		} else if line.ValueType == PaymentTermValueTypePercentage {
			percent, _ := ParseDecimal(*line.ValueAmount)
			percentFactor := new(big.Float).Quo(percent, new(big.Float).SetInt64(100))
			amountDue = new(big.Float).Mul(invoiceAmountDecimal, percentFactor)
		} else { // Fixed
			amountDue, _ = ParseDecimal(*line.ValueAmount)
		}

		// Update balance
		balanceAmount = new(big.Float).Sub(balanceAmount, amountDue)

		payments = append(payments, &CalculatedPayment{
			LineNumber:  line.Sequence,
			DueDate:     dueDate,
			AmountDue:   FormatDecimal(amountDue, 4),
			PaymentType: line.ValueType.String(),
		})
	}

	return &PaymentScheduleCalculation{
		Payments:     payments,
		TotalAmount:  invoiceAmount,
		RemainingDue: FormatDecimal(balanceAmount, 4),
	}, nil
}

// CreateInvoicePaymentSchedules creates payment schedules for an invoice
func (s *Service) CreateInvoicePaymentSchedules(ctx context.Context, organizationID uuid.UUID, req *CreateInvoicePaymentScheduleRequest) ([]*InvoicePaymentSchedule, error) {
	schedules := make([]*InvoicePaymentSchedule, 0, len(req.Lines))

	sourceType := SourceType(req.SourceType)
	now := time.Now()

	for _, lineReq := range req.Lines {
		schedule := &InvoicePaymentSchedule{
			ID:             uuid.New(),
			OrganizationID: organizationID,
			SourceType:     sourceType,
			SourceID:       req.SourceID,
			LineNumber:     lineReq.LineNumber,
			DueDate:        lineReq.DueDate,
			AmountDue:      lineReq.AmountDue,
			AmountPaid:     "0",
			Status:         InvoicePaymentScheduleStatusPending,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		schedules = append(schedules, schedule)
	}

	if err := s.repo.CreateInvoicePaymentSchedules(ctx, schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

// GetInvoicePaymentSchedule retrieves a payment schedule
func (s *Service) GetInvoicePaymentSchedule(ctx context.Context, id uuid.UUID) (*InvoicePaymentSchedule, error) {
	return s.repo.GetInvoicePaymentSchedule(ctx, id)
}

// ListInvoicePaymentSchedules lists payment schedules
func (s *Service) ListInvoicePaymentSchedules(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*InvoicePaymentSchedule, error) {
	return s.repo.ListInvoicePaymentSchedules(ctx, organizationID, filter)
}

// ListInvoicePaymentSchedulesBySource lists schedules for a specific invoice/bill
func (s *Service) ListInvoicePaymentSchedulesBySource(ctx context.Context, organizationID uuid.UUID, sourceType SourceType, sourceID uuid.UUID) ([]*InvoicePaymentSchedule, error) {
	return s.repo.ListInvoicePaymentSchedulesBySource(ctx, organizationID, sourceType, sourceID)
}

// UpdateInvoicePaymentSchedule updates a payment schedule (mark as paid, etc.)
func (s *Service) UpdateInvoicePaymentSchedule(ctx context.Context, id uuid.UUID, req *UpdateInvoicePaymentScheduleRequest) (*InvoicePaymentSchedule, error) {
	schedule, err := s.repo.GetInvoicePaymentSchedule(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.AmountPaid != nil {
		schedule.AmountPaid = *req.AmountPaid

		// Update status based on amounts
		amountPaid, _ := ParseDecimal(schedule.AmountPaid)
		amountDue, _ := ParseDecimal(schedule.AmountDue)

		cmp := amountPaid.Cmp(amountDue)
		if cmp >= 0 {
			schedule.Status = InvoicePaymentScheduleStatusPaid
		} else if cmp > 0 {
			schedule.Status = InvoicePaymentScheduleStatusPartial
		} else {
			schedule.Status = InvoicePaymentScheduleStatusPending
		}
	}

	if req.Status != nil {
		schedule.Status = InvoicePaymentScheduleStatus(*req.Status)
	}

	schedule.UpdatedAt = time.Now()

	if err := s.repo.UpdateInvoicePaymentSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}
