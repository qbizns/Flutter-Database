package rest

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/accounting"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"github.com/your-org/pos-backend/internal/repository/postgres"
	"go.uber.org/zap"
)

// TrialBalanceHandler handles GET /api/v1/organizations/{org_id}/reports/trial-balance
func TrialBalanceHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse as_of_date parameter
		asOfDateStr := r.URL.Query().Get("as_of_date")
		var asOfDate time.Time
		if asOfDateStr != "" {
			asOfDate, err = time.Parse("2006-01-02", asOfDateStr)
			if err != nil {
				respondError(w, logger, apperrors.BadRequest("Invalid as_of_date format. Use YYYY-MM-DD"))
				return
			}
		} else {
			asOfDate = time.Now()
		}

		// Create repository and service
		repo := postgres.NewAccountingRepository(db)
		service := accounting.NewService(repo)

		// Get trial balance
		trialBalance, err := service.GetTrialBalance(ctx, orgID, asOfDate)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Get chart of accounts for account details
		accounts, err := service.ListChartOfAccounts(ctx, orgID, nil)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Build response with account details
		type TrialBalanceRow struct {
			AccountID     uuid.UUID `json:"account_id"`
			AccountCode   string    `json:"account_code"`
			AccountName   string    `json:"account_name"`
			DebitBalance  string    `json:"debit_balance"`
			CreditBalance string    `json:"credit_balance"`
		}

		rows := make([]TrialBalanceRow, 0, len(trialBalance))
		totalDebit := "0"
		totalCredit := "0"

		for accountID, balances := range trialBalance {
			// Find account details
			var accountCode, accountName string
			for _, acc := range accounts {
				if acc.ID == accountID {
					accountCode = acc.AccountCode
					accountName = acc.AccountName
					break
				}
			}

			rows = append(rows, TrialBalanceRow{
				AccountID:     accountID,
				AccountCode:   accountCode,
				AccountName:   accountName,
				DebitBalance:  balances["debit"],
				CreditBalance: balances["credit"],
			})

			// TODO: Sum totals (requires decimal arithmetic)
		}

		logger.Info("trial balance generated",
			zap.String("as_of_date", asOfDate.Format("2006-01-02")),
			zap.Int("account_count", len(rows)),
		)

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"trial_balance": rows,
			"as_of_date":    asOfDate.Format("2006-01-02"),
			"total_debit":   totalDebit,
			"total_credit":  totalCredit,
		})
	}
}

// BalanceSheetHandler handles GET /api/v1/organizations/{org_id}/reports/balance-sheet
func BalanceSheetHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse as_of_date parameter
		asOfDateStr := r.URL.Query().Get("as_of_date")
		var asOfDate time.Time
		if asOfDateStr != "" {
			asOfDate, err = time.Parse("2006-01-02", asOfDateStr)
			if err != nil {
				respondError(w, logger, apperrors.BadRequest("Invalid as_of_date format. Use YYYY-MM-DD"))
				return
			}
		} else {
			asOfDate = time.Now()
		}

		// Create repository and service
		repo := postgres.NewAccountingRepository(db)
		service := accounting.NewService(repo)

		// Get chart of accounts
		accounts, err := service.ListChartOfAccounts(ctx, orgID, map[string]interface{}{
			"is_active": true,
		})
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Get account types
		accountTypes, err := repo.ListAccountTypes(ctx)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Build balance sheet structure
		type AccountLine struct {
			AccountID     uuid.UUID `json:"account_id"`
			AccountCode   string    `json:"account_code"`
			AccountName   string    `json:"account_name"`
			Balance       string    `json:"balance"`
			AccountLevel  int       `json:"account_level"`
		}

		type BalanceSheetSection struct {
			Title    string         `json:"title"`
			Accounts []AccountLine  `json:"accounts"`
			Total    string         `json:"total"`
		}

		assets := BalanceSheetSection{Title: "Assets", Accounts: []AccountLine{}, Total: "0"}
		liabilities := BalanceSheetSection{Title: "Liabilities", Accounts: []AccountLine{}, Total: "0"}
		equity := BalanceSheetSection{Title: "Equity", Accounts: []AccountLine{}, Total: "0"}

		// Group accounts by type
		for _, account := range accounts {
			if account.IsHeaderAccount {
				continue
			}

			// Get account balance
			debit, credit, balance, err := service.GetAccountBalance(ctx, orgID, account.ID, asOfDate)
			if err != nil {
				logger.Warn("failed to get account balance", zap.Error(err), zap.String("account_id", account.ID.String()))
				continue
			}

			line := AccountLine{
				AccountID:    account.ID,
				AccountCode:  account.AccountCode,
				AccountName:  account.AccountName,
				Balance:      balance,
				AccountLevel: account.AccountLevel,
			}

			// Find account type category
			for _, at := range accountTypes {
				if at.ID == account.AccountTypeID {
					switch at.TypeCategory {
					case "asset":
						assets.Accounts = append(assets.Accounts, line)
					case "liability":
						liabilities.Accounts = append(liabilities.Accounts, line)
					case "equity":
						equity.Accounts = append(equity.Accounts, line)
					}
					break
				}
			}

			// Use debit/credit to avoid compiler warnings
			_ = debit
			_ = credit
		}

		logger.Info("balance sheet generated",
			zap.String("as_of_date", asOfDate.Format("2006-01-02")),
		)

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"balance_sheet": map[string]interface{}{
				"assets":      assets,
				"liabilities": liabilities,
				"equity":      equity,
			},
			"as_of_date": asOfDate.Format("2006-01-02"),
		})
	}
}

// IncomeStatementHandler handles GET /api/v1/organizations/{org_id}/reports/income-statement
func IncomeStatementHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Safe context extraction
		orgID, err := getOrganizationID(r)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Parse date parameters
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		var startDate, endDate time.Time
		var parseErr error

		if startDateStr != "" {
			startDate, parseErr = time.Parse("2006-01-02", startDateStr)
			if parseErr != nil {
				respondError(w, logger, apperrors.BadRequest("Invalid start_date format. Use YYYY-MM-DD"))
				return
			}
		} else {
			// Default to current year start
			year := time.Now().Year()
			startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		if endDateStr != "" {
			endDate, parseErr = time.Parse("2006-01-02", endDateStr)
			if parseErr != nil {
				respondError(w, logger, apperrors.BadRequest("Invalid end_date format. Use YYYY-MM-DD"))
				return
			}
		} else {
			endDate = time.Now()
		}

		// Create repository and service
		repo := postgres.NewAccountingRepository(db)
		service := accounting.NewService(repo)

		// Get chart of accounts
		accounts, err := service.ListChartOfAccounts(ctx, orgID, map[string]interface{}{
			"is_active": true,
		})
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Get account types
		accountTypes, err := repo.ListAccountTypes(ctx)
		if err != nil {
			respondError(w, logger, err)
			return
		}

		// Build income statement structure
		type AccountLine struct {
			AccountID     uuid.UUID `json:"account_id"`
			AccountCode   string    `json:"account_code"`
			AccountName   string    `json:"account_name"`
			Amount        string    `json:"amount"`
			AccountLevel  int       `json:"account_level"`
		}

		type IncomeStatementSection struct {
			Title    string         `json:"title"`
			Accounts []AccountLine  `json:"accounts"`
			Total    string         `json:"total"`
		}

		revenue := IncomeStatementSection{Title: "Revenue", Accounts: []AccountLine{}, Total: "0"}
		expenses := IncomeStatementSection{Title: "Expenses", Accounts: []AccountLine{}, Total: "0"}

		// Group accounts by type
		for _, account := range accounts {
			if account.IsHeaderAccount {
				continue
			}

			// Get account balance as of end date
			_, _, balance, err := service.GetAccountBalance(ctx, orgID, account.ID, endDate)
			if err != nil {
				logger.Warn("failed to get account balance", zap.Error(err), zap.String("account_id", account.ID.String()))
				continue
			}

			line := AccountLine{
				AccountID:    account.ID,
				AccountCode:  account.AccountCode,
				AccountName:  account.AccountName,
				Amount:       balance,
				AccountLevel: account.AccountLevel,
			}

			// Find account type category
			for _, at := range accountTypes {
				if at.ID == account.AccountTypeID {
					switch at.TypeCategory {
					case "revenue", "income":
						revenue.Accounts = append(revenue.Accounts, line)
					case "expense", "cost_of_goods_sold", "operating_expense":
						expenses.Accounts = append(expenses.Accounts, line)
					}
					break
				}
			}
		}

		// Calculate net income (simplified - needs proper decimal arithmetic)
		netIncome := "0"

		logger.Info("income statement generated",
			zap.String("start_date", startDate.Format("2006-01-02")),
			zap.String("end_date", endDate.Format("2006-01-02")),
		)

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"income_statement": map[string]interface{}{
				"revenue":    revenue,
				"expenses":   expenses,
				"net_income": netIncome,
			},
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		})
	}
}
