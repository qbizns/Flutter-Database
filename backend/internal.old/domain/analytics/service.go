package analytics

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// Repository defines data access operations for analytics
type Repository interface {
	// Analytic Plans
	CreateAnalyticPlan(ctx context.Context, plan *AnalyticPlan) error
	GetAnalyticPlan(ctx context.Context, id uuid.UUID) (*AnalyticPlan, error)
	GetAnalyticPlanByCode(ctx context.Context, organizationID uuid.UUID, code string) (*AnalyticPlan, error)
	ListAnalyticPlans(ctx context.Context, query AnalyticQuery) ([]*AnalyticPlan, int64, error)
	UpdateAnalyticPlan(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteAnalyticPlan(ctx context.Context, id uuid.UUID) error

	// Analytic Accounts
	CreateAnalyticAccount(ctx context.Context, account *AnalyticAccount) error
	GetAnalyticAccount(ctx context.Context, id uuid.UUID) (*AnalyticAccount, error)
	GetAnalyticAccountByCode(ctx context.Context, organizationID uuid.UUID, code string) (*AnalyticAccount, error)
	ListAnalyticAccounts(ctx context.Context, query AnalyticQuery) ([]*AnalyticAccount, int64, error)
	UpdateAnalyticAccount(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteAnalyticAccount(ctx context.Context, id uuid.UUID) error
	GetAnalyticAccountHierarchy(ctx context.Context, organizationID uuid.UUID) ([]*AnalyticAccount, error)

	// Deferred Revenue
	CreateDeferredRevenueContract(ctx context.Context, contract *DeferredRevenueContract) error
	GetDeferredRevenueContract(ctx context.Context, id uuid.UUID) (*DeferredRevenueContract, error)
	ListDeferredRevenueContracts(ctx context.Context, query DeferralQuery) ([]*DeferredRevenueContract, int64, error)
	UpdateDeferredRevenueContract(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteDeferredRevenueContract(ctx context.Context, id uuid.UUID) error

	// Deferred Revenue Schedules
	CreateDeferredRevenueSchedule(ctx context.Context, schedule *DeferredRevenueSchedule) error
	GetDeferredRevenueSchedule(ctx context.Context, id uuid.UUID) (*DeferredRevenueSchedule, error)
	ListDeferredRevenueSchedules(ctx context.Context, contractID uuid.UUID, limit, offset int) ([]*DeferredRevenueSchedule, int64, error)
	UpdateDeferredRevenueSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteDeferredRevenueSchedule(ctx context.Context, id uuid.UUID) error
	GetPendingDeferredRevenueSchedules(ctx context.Context, organizationID uuid.UUID, beforeDate time.Time) ([]*DeferredRevenueSchedule, error)

	// Deferred Expense
	CreateDeferredExpenseContract(ctx context.Context, contract *DeferredExpenseContract) error
	GetDeferredExpenseContract(ctx context.Context, id uuid.UUID) (*DeferredExpenseContract, error)
	ListDeferredExpenseContracts(ctx context.Context, query DeferralQuery) ([]*DeferredExpenseContract, int64, error)
	UpdateDeferredExpenseContract(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteDeferredExpenseContract(ctx context.Context, id uuid.UUID) error

	// Deferred Expense Schedules
	CreateDeferredExpenseSchedule(ctx context.Context, schedule *DeferredExpenseSchedule) error
	GetDeferredExpenseSchedule(ctx context.Context, id uuid.UUID) (*DeferredExpenseSchedule, error)
	ListDeferredExpenseSchedules(ctx context.Context, contractID uuid.UUID, limit, offset int) ([]*DeferredExpenseSchedule, int64, error)
	UpdateDeferredExpenseSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteDeferredExpenseSchedule(ctx context.Context, id uuid.UUID) error
	GetPendingDeferredExpenseSchedules(ctx context.Context, organizationID uuid.UUID, beforeDate time.Time) ([]*DeferredExpenseSchedule, error)

	// Budgets
	CreateBudget(ctx context.Context, budget *Budget) error
	GetBudget(ctx context.Context, id uuid.UUID) (*Budget, error)
	GetBudgetByCode(ctx context.Context, organizationID uuid.UUID, code string) (*Budget, error)
	ListBudgets(ctx context.Context, query BudgetQuery) ([]*Budget, int64, error)
	UpdateBudget(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteBudget(ctx context.Context, id uuid.UUID) error

	// Budget Lines
	CreateBudgetLine(ctx context.Context, line *BudgetLine) error
	GetBudgetLine(ctx context.Context, id uuid.UUID) (*BudgetLine, error)
	ListBudgetLines(ctx context.Context, budgetID uuid.UUID, limit, offset int) ([]*BudgetLine, int64, error)
	UpdateBudgetLine(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteBudgetLine(ctx context.Context, id uuid.UUID) error
}

// Service provides business logic for analytics and budgeting
type Service struct {
	repo Repository
}

// NewService creates a new analytics service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ==================== ANALYTIC PLANS ====================

// CreateAnalyticPlan creates a new analytic plan
func (s *Service) CreateAnalyticPlan(ctx context.Context, req *CreateAnalyticPlanRequest, organizationID uuid.UUID) (*AnalyticPlan, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if req.PlanCode == "" {
		return nil, ErrAnalyticPlanCodeRequired
	}

	if req.PlanName == "" {
		return nil, ErrAnalyticPlanNameRequired
	}

	// Check for duplicate
	existing, _ := s.repo.GetAnalyticPlanByCode(ctx, organizationID, req.PlanCode)
	if existing != nil {
		return nil, ErrDuplicateAnalyticPlanCode
	}

	plan := &AnalyticPlan{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		PlanCode:       req.PlanCode,
		PlanName:       req.PlanName,
		IsActive:       true,
		Description:    req.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateAnalyticPlan(ctx, plan); err != nil {
		return nil, fmt.Errorf("failed to create analytic plan: %w", err)
	}

	return plan, nil
}

// GetAnalyticPlan retrieves an analytic plan
func (s *Service) GetAnalyticPlan(ctx context.Context, id uuid.UUID) (*AnalyticPlan, error) {
	plan, err := s.repo.GetAnalyticPlan(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic plan: %w", err)
	}
	if plan == nil {
		return nil, ErrAnalyticPlanNotFound
	}
	return plan, nil
}

// ListAnalyticPlans retrieves analytic plans
func (s *Service) ListAnalyticPlans(ctx context.Context, query AnalyticQuery) ([]*AnalyticPlan, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	plans, total, err := s.repo.ListAnalyticPlans(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analytic plans: %w", err)
	}

	return plans, total, nil
}

// UpdateAnalyticPlan updates an analytic plan
func (s *Service) UpdateAnalyticPlan(ctx context.Context, id uuid.UUID, req *UpdateAnalyticPlanRequest) error {
	if _, err := s.GetAnalyticPlan(ctx, id); err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.PlanName != nil {
		updates["plan_name"] = *req.PlanName
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateAnalyticPlan(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update analytic plan: %w", err)
	}

	return nil
}

// ==================== ANALYTIC ACCOUNTS ====================

// CreateAnalyticAccount creates a new analytic account
func (s *Service) CreateAnalyticAccount(ctx context.Context, req *CreateAnalyticAccountRequest, organizationID uuid.UUID) (*AnalyticAccount, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if req.AccountCode == "" {
		return nil, ErrAnalyticAccountCodeRequired
	}

	if req.AccountName == "" {
		return nil, ErrAnalyticAccountNameRequired
	}

	// Check for duplicate
	existing, _ := s.repo.GetAnalyticAccountByCode(ctx, organizationID, req.AccountCode)
	if existing != nil {
		return nil, ErrDuplicateAnalyticAccountCode
	}

	// Check for circular reference
	if req.ParentAccountID != nil {
		if err := s.checkCircularReference(ctx, *req.ParentAccountID, uuid.Nil); err != nil {
			return nil, err
		}
	}

	accountLevel := 1
	if req.ParentAccountID != nil {
		parent, err := s.GetAnalyticAccount(ctx, *req.ParentAccountID)
		if err != nil {
			return nil, err
		}
		accountLevel = parent.AccountLevel + 1
	}

	account := &AnalyticAccount{
		ID:              uuid.New(),
		OrganizationID:  organizationID,
		AnalyticPlanID:  req.AnalyticPlanID,
		AccountCode:     req.AccountCode,
		AccountName:     req.AccountName,
		ParentAccountID: req.ParentAccountID,
		AccountLevel:    accountLevel,
		IsActive:        true,
		Description:     req.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateAnalyticAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create analytic account: %w", err)
	}

	return account, nil
}

// GetAnalyticAccount retrieves an analytic account
func (s *Service) GetAnalyticAccount(ctx context.Context, id uuid.UUID) (*AnalyticAccount, error) {
	account, err := s.repo.GetAnalyticAccount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytic account: %w", err)
	}
	if account == nil {
		return nil, ErrAnalyticAccountNotFound
	}
	return account, nil
}

// ListAnalyticAccounts retrieves analytic accounts
func (s *Service) ListAnalyticAccounts(ctx context.Context, query AnalyticQuery) ([]*AnalyticAccount, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	accounts, total, err := s.repo.ListAnalyticAccounts(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analytic accounts: %w", err)
	}

	return accounts, total, nil
}

// GetAnalyticAccountHierarchy retrieves the full hierarchy
func (s *Service) GetAnalyticAccountHierarchy(ctx context.Context, organizationID uuid.UUID) ([]*AnalyticAccount, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	accounts, err := s.repo.GetAnalyticAccountHierarchy(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hierarchy: %w", err)
	}

	return accounts, nil
}

// UpdateAnalyticAccount updates an analytic account
func (s *Service) UpdateAnalyticAccount(ctx context.Context, id uuid.UUID, req *UpdateAnalyticAccountRequest) error {
	if _, err := s.GetAnalyticAccount(ctx, id); err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.AccountName != nil {
		updates["account_name"] = *req.AccountName
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateAnalyticAccount(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update analytic account: %w", err)
	}

	return nil
}

// checkCircularReference prevents circular hierarchies
func (s *Service) checkCircularReference(ctx context.Context, parentID uuid.UUID, childID uuid.UUID) error {
	parent, err := s.GetAnalyticAccount(ctx, parentID)
	if err != nil {
		return err
	}

	if parent.ParentAccountID == nil {
		return nil
	}

	if *parent.ParentAccountID == childID {
		return ErrCircularAnalyticReference
	}

	return s.checkCircularReference(ctx, *parent.ParentAccountID, childID)
}

// ==================== DEFERRED REVENUE ====================

// CreateDeferredRevenueContract creates a deferred revenue contract
func (s *Service) CreateDeferredRevenueContract(ctx context.Context, req *CreateDeferredRevenueRequest, organizationID uuid.UUID) (*DeferredRevenueContract, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if req.TotalDeferredAmount <= 0 {
		return nil, ErrInvalidDeferralAmount
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDeferralDates
	}

	contract := &DeferredRevenueContract{
		ID:                  uuid.New(),
		OrganizationID:      organizationID,
		CustomerInvoiceID:   req.CustomerInvoiceID,
		InvoiceLineID:       req.InvoiceLineID,
		ContractName:        req.ContractName,
		TotalDeferredAmount: req.TotalDeferredAmount,
		StartDate:           req.StartDate,
		EndDate:             req.EndDate,
		RecognitionMethod:   req.RecognitionMethod,
		DeferredAccountID:   req.DeferredAccountID,
		RevenueAccountID:    req.RevenueAccountID,
		Status:              DeferralStatusDraft,
		RecognizedAmount:    0,
		Notes:               req.Notes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.CreateDeferredRevenueContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to create deferred revenue contract: %w", err)
	}

	return contract, nil
}

// GetDeferredRevenueContract retrieves a deferred revenue contract
func (s *Service) GetDeferredRevenueContract(ctx context.Context, id uuid.UUID) (*DeferredRevenueContract, error) {
	contract, err := s.repo.GetDeferredRevenueContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deferred revenue contract: %w", err)
	}
	if contract == nil {
		return nil, ErrDeferredRevenueNotFound
	}
	return contract, nil
}

// ListDeferredRevenueContracts retrieves deferred revenue contracts
func (s *Service) ListDeferredRevenueContracts(ctx context.Context, query DeferralQuery) ([]*DeferredRevenueContract, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	contracts, total, err := s.repo.ListDeferredRevenueContracts(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deferred revenue contracts: %w", err)
	}

	return contracts, total, nil
}

// UpdateDeferredRevenueContract updates a deferred revenue contract
func (s *Service) UpdateDeferredRevenueContract(ctx context.Context, id uuid.UUID, req *UpdateDeferredRevenueRequest) error {
	contract, err := s.GetDeferredRevenueContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.Status == DeferralStatusCompleted {
		return ErrDeferralAlreadyCompleted
	}

	updates := make(map[string]interface{})
	if req.ContractName != nil {
		updates["contract_name"] = *req.ContractName
	}
	if req.EndDate != nil {
		if req.EndDate.Before(contract.StartDate) {
			return ErrInvalidDeferralDates
		}
		updates["end_date"] = *req.EndDate
	}
	if req.RecognitionMethod != nil {
		updates["recognition_method"] = *req.RecognitionMethod
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateDeferredRevenueContract(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update deferred revenue contract: %w", err)
	}

	return nil
}

// GenerateDeferredRevenueSchedule generates recognition schedule
func (s *Service) GenerateDeferredRevenueSchedule(ctx context.Context, contractID uuid.UUID) error {
	contract, err := s.GetDeferredRevenueContract(ctx, contractID)
	if err != nil {
		return err
	}

	if contract.Status != DeferralStatusDraft {
		return fmt.Errorf("can only generate schedule for draft contracts")
	}

	// Get existing schedules
	existing, _, err := s.repo.ListDeferredRevenueSchedules(ctx, contractID, 1000, 0)
	if err == nil && len(existing) > 0 {
		// Delete existing schedules
		for _, sch := range existing {
			s.repo.DeleteDeferredRevenueSchedule(ctx, sch.ID)
		}
	}

	// Generate based on recognition method
	var schedules []*DeferredRevenueSchedule
	switch contract.RecognitionMethod {
	case RecognitionMethodStraightLine:
		schedules = s.generateStraightLineSchedule(contract)
	case RecognitionMethodMilestone:
		// Milestone-based would be manually entered
		return nil
	case RecognitionMethodCustom:
		// Custom would be manually entered
		return nil
	}

	// Create schedules
	for _, schedule := range schedules {
		if err := s.repo.CreateDeferredRevenueSchedule(ctx, schedule); err != nil {
			return fmt.Errorf("failed to create schedule: %w", err)
		}
	}

	return nil
}

func (s *Service) generateStraightLineSchedule(contract *DeferredRevenueContract) []*DeferredRevenueSchedule {
	// Calculate monthly recognition
	daysDiff := int(contract.EndDate.Sub(contract.StartDate).Hours() / 24)
	months := (daysDiff + 14) / 30 // Approximate month count
	if months <= 0 {
		months = 1
	}

	monthlyAmount := contract.TotalDeferredAmount / float64(months)
	var schedules []*DeferredRevenueSchedule

	currentDate := contract.StartDate
	for i := 1; i <= months; i++ {
		// Set to end of month or contract end
		nextDate := currentDate.AddDate(0, 1, 0)
		if nextDate.After(contract.EndDate) {
			nextDate = contract.EndDate
		}

		amount := monthlyAmount
		if i == months {
			// Last entry gets remainder
			amount = contract.TotalDeferredAmount - (monthlyAmount * float64(i-1))
		}

		schedule := &DeferredRevenueSchedule{
			ID:                uuid.New(),
			ContractID:        contract.ID,
			LineNumber:        i,
			RecognitionDate:   currentDate.AddDate(0, 0, 1), // Next day
			RecognitionAmount: amount,
			Status:            ScheduleStatusPending,
			CreatedAt:         time.Now(),
		}

		schedules = append(schedules, schedule)
		currentDate = nextDate
	}

	return schedules
}

// ==================== DEFERRED EXPENSE ====================

// CreateDeferredExpenseContract creates a deferred expense contract
func (s *Service) CreateDeferredExpenseContract(ctx context.Context, req *CreateDeferredExpenseRequest, organizationID uuid.UUID) (*DeferredExpenseContract, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if req.TotalDeferredAmount <= 0 {
		return nil, ErrInvalidDeferralAmount
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDeferralDates
	}

	contract := &DeferredExpenseContract{
		ID:                  uuid.New(),
		OrganizationID:      organizationID,
		VendorBillID:        req.VendorBillID,
		BillLineID:          req.BillLineID,
		ContractName:        req.ContractName,
		TotalDeferredAmount: req.TotalDeferredAmount,
		StartDate:           req.StartDate,
		EndDate:             req.EndDate,
		RecognitionMethod:   req.RecognitionMethod,
		DeferredAccountID:   req.DeferredAccountID,
		ExpenseAccountID:    req.ExpenseAccountID,
		Status:              DeferralStatusDraft,
		RecognizedAmount:    0,
		Notes:               req.Notes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.CreateDeferredExpenseContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to create deferred expense contract: %w", err)
	}

	return contract, nil
}

// GetDeferredExpenseContract retrieves a deferred expense contract
func (s *Service) GetDeferredExpenseContract(ctx context.Context, id uuid.UUID) (*DeferredExpenseContract, error) {
	contract, err := s.repo.GetDeferredExpenseContract(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deferred expense contract: %w", err)
	}
	if contract == nil {
		return nil, ErrDeferredExpenseNotFound
	}
	return contract, nil
}

// ListDeferredExpenseContracts retrieves deferred expense contracts
func (s *Service) ListDeferredExpenseContracts(ctx context.Context, query DeferralQuery) ([]*DeferredExpenseContract, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	contracts, total, err := s.repo.ListDeferredExpenseContracts(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list deferred expense contracts: %w", err)
	}

	return contracts, total, nil
}

// UpdateDeferredExpenseContract updates a deferred expense contract
func (s *Service) UpdateDeferredExpenseContract(ctx context.Context, id uuid.UUID, req *UpdateDeferredExpenseRequest) error {
	contract, err := s.GetDeferredExpenseContract(ctx, id)
	if err != nil {
		return err
	}

	if contract.Status == DeferralStatusCompleted {
		return ErrDeferralAlreadyCompleted
	}

	updates := make(map[string]interface{})
	if req.ContractName != nil {
		updates["contract_name"] = *req.ContractName
	}
	if req.EndDate != nil {
		if req.EndDate.Before(contract.StartDate) {
			return ErrInvalidDeferralDates
		}
		updates["end_date"] = *req.EndDate
	}
	if req.RecognitionMethod != nil {
		updates["recognition_method"] = *req.RecognitionMethod
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateDeferredExpenseContract(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update deferred expense contract: %w", err)
	}

	return nil
}

// ==================== BUDGETS ====================

// CreateBudget creates a new budget
func (s *Service) CreateBudget(ctx context.Context, req *CreateBudgetRequest, organizationID uuid.UUID) (*Budget, error) {
	if organizationID == uuid.Nil {
		return nil, ErrOrganizationIDRequired
	}

	if req.BudgetCode == "" {
		return nil, ErrBudgetCodeRequired
	}

	if req.BudgetName == "" {
		return nil, ErrBudgetNameRequired
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, ErrInvalidDeferralDates
	}

	// Check for duplicate
	existing, _ := s.repo.GetBudgetByCode(ctx, organizationID, req.BudgetCode)
	if existing != nil {
		return nil, ErrDuplicateBudgetCode
	}

	budget := &Budget{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		BudgetCode:     req.BudgetCode,
		BudgetName:     req.BudgetName,
		FiscalYearID:   req.FiscalYearID,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		BudgetType:     req.BudgetType,
		Status:         BudgetStatusDraft,
		Notes:          req.Notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateBudget(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return budget, nil
}

// GetBudget retrieves a budget
func (s *Service) GetBudget(ctx context.Context, id uuid.UUID) (*Budget, error) {
	budget, err := s.repo.GetBudget(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	if budget == nil {
		return nil, ErrBudgetNotFound
	}
	return budget, nil
}

// ListBudgets retrieves budgets
func (s *Service) ListBudgets(ctx context.Context, query BudgetQuery) ([]*Budget, int64, error) {
	if query.OrganizationID == uuid.Nil {
		return nil, 0, ErrOrganizationIDRequired
	}

	budgets, total, err := s.repo.ListBudgets(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list budgets: %w", err)
	}

	return budgets, total, nil
}

// UpdateBudget updates a budget
func (s *Service) UpdateBudget(ctx context.Context, id uuid.UUID, req *UpdateBudgetRequest) error {
	budget, err := s.GetBudget(ctx, id)
	if err != nil {
		return err
	}

	if budget.Status != BudgetStatusDraft {
		return ErrBudgetCannotBeModified
	}

	updates := make(map[string]interface{})
	if req.BudgetName != nil {
		updates["budget_name"] = *req.BudgetName
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateBudget(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}

	return nil
}

// ApproveBudget approves a budget
func (s *Service) ApproveBudget(ctx context.Context, id uuid.UUID) error {
	budget, err := s.GetBudget(ctx, id)
	if err != nil {
		return err
	}

	if budget.Status != BudgetStatusDraft {
		return ErrBudgetAlreadyApproved
	}

	return s.repo.UpdateBudget(ctx, id, map[string]interface{}{
		"status":     BudgetStatusApproved,
		"updated_at": time.Now(),
	})
}

// ActivateBudget activates a budget
func (s *Service) ActivateBudget(ctx context.Context, id uuid.UUID) error {
	budget, err := s.GetBudget(ctx, id)
	if err != nil {
		return err
	}

	if budget.Status != BudgetStatusApproved && budget.Status != BudgetStatusDraft {
		return ErrBudgetAlreadyActive
	}

	return s.repo.UpdateBudget(ctx, id, map[string]interface{}{
		"status":     BudgetStatusActive,
		"updated_at": time.Now(),
	})
}

// CreateBudgetLine creates a budget line
func (s *Service) CreateBudgetLine(ctx context.Context, req *CreateBudgetLineRequest) (*BudgetLine, error) {
	if req.BudgetID == uuid.Nil {
		return nil, fmt.Errorf("budget_id is required")
	}

	if req.AccountID == nil && req.AnalyticAccountID == nil {
		return nil, ErrBudgetLineNeedsDimension
	}

	line := &BudgetLine{
		ID:                 uuid.New(),
		BudgetID:           req.BudgetID,
		AccountID:          req.AccountID,
		AnalyticAccountID:  req.AnalyticAccountID,
		AccountingPeriodID: req.AccountingPeriodID,
		PeriodStartDate:    req.PeriodStartDate,
		PeriodEndDate:      req.PeriodEndDate,
		PlannedAmount:      req.PlannedAmount,
		Notes:              req.Notes,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.repo.CreateBudgetLine(ctx, line); err != nil {
		return nil, fmt.Errorf("failed to create budget line: %w", err)
	}

	return line, nil
}

// GetBudgetLine retrieves a budget line
func (s *Service) GetBudgetLine(ctx context.Context, id uuid.UUID) (*BudgetLine, error) {
	line, err := s.repo.GetBudgetLine(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget line: %w", err)
	}
	if line == nil {
		return nil, ErrBudgetLineNotFound
	}
	return line, nil
}

// ListBudgetLines retrieves budget lines
func (s *Service) ListBudgetLines(ctx context.Context, budgetID uuid.UUID, limit, offset int) ([]*BudgetLine, int64, error) {
	lines, total, err := s.repo.ListBudgetLines(ctx, budgetID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list budget lines: %w", err)
	}
	return lines, total, nil
}

// UpdateBudgetLine updates a budget line
func (s *Service) UpdateBudgetLine(ctx context.Context, id uuid.UUID, req *UpdateBudgetLineRequest) error {
	if _, err := s.GetBudgetLine(ctx, id); err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.PlannedAmount != nil {
		updates["planned_amount"] = *req.PlannedAmount
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	updates["updated_at"] = time.Now()

	if err := s.repo.UpdateBudgetLine(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update budget line: %w", err)
	}

	return nil
}

// GetPendingDeferrals retrieves pending deferrals for posting
func (s *Service) GetPendingDeferrals(ctx context.Context, organizationID uuid.UUID, beforeDate time.Time) (
	[]DeferralScheduleItem, error) {

	revSchedules, err := s.repo.GetPendingDeferredRevenueSchedules(ctx, organizationID, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue schedules: %w", err)
	}

	expSchedules, err := s.repo.GetPendingDeferredExpenseSchedules(ctx, organizationID, beforeDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense schedules: %w", err)
	}

	var items []DeferralScheduleItem

	for _, rev := range revSchedules {
		items = append(items, DeferralScheduleItem{
			ID:                rev.ID,
			ContractID:        rev.ContractID,
			LineNumber:        rev.LineNumber,
			RecognitionDate:   rev.RecognitionDate,
			RecognitionAmount: rev.RecognitionAmount,
			Status:            rev.Status,
			IsPosted:          rev.Status == ScheduleStatusPosted,
		})
	}

	for _, exp := range expSchedules {
		items = append(items, DeferralScheduleItem{
			ID:                exp.ID,
			ContractID:        exp.ContractID,
			LineNumber:        exp.LineNumber,
			RecognitionDate:   exp.RecognitionDate,
			RecognitionAmount: exp.RecognitionAmount,
			Status:            exp.Status,
			IsPosted:          exp.Status == ScheduleStatusPosted,
		})
	}

	return items, nil
}

// CalculateBudgetVariance calculates variance for a budget
func (s *Service) CalculateBudgetVariance(planned, actual float64) float64 {
	if planned == 0 {
		return 0
	}
	return ((actual - planned) / math.Abs(planned)) * 100
}
