package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/giftcards"
)

type GiftCardRepository struct {
	db *DB
}

func NewGiftCardRepository(db *DB) *GiftCardRepository {
	return &GiftCardRepository{db: db}
}

// ============= GIFT CARDS =============

func (r *GiftCardRepository) CreateGiftCard(ctx context.Context, card *giftcards.GiftCard) error {
	if err := r.db.SetOrganizationContext(ctx, card.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO gift_cards (
			id, organization_id, card_number, pin_code, customer_id,
			original_value, current_balance, issued_date, expiry_date,
			status, issued_by_user_id, issued_location_id, notes,
			created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		card.ID, card.OrganizationID, card.CardNumber, card.PINCode, card.CustomerID,
		card.OriginalValue, card.CurrentBalance, card.IssuedDate, card.ExpiryDate,
		card.Status, card.IssuedByUserID, card.IssuedLocationID, card.Notes,
		card.CreatedBy, card.CreatedAt, card.UpdatedAt,
	)
	return err
}

func (r *GiftCardRepository) GetGiftCard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*giftcards.GiftCard, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, card_number, pin_code, customer_id,
		       original_value, current_balance, issued_date, expiry_date,
		       status, issued_by_user_id, issued_location_id, notes,
		       created_by, created_at, updated_at, deleted_at
		FROM gift_cards
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var card giftcards.GiftCard
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&card.ID, &card.OrganizationID, &card.CardNumber, &card.PINCode, &card.CustomerID,
		&card.OriginalValue, &card.CurrentBalance, &card.IssuedDate, &card.ExpiryDate,
		&card.Status, &card.IssuedByUserID, &card.IssuedLocationID, &card.Notes,
		&card.CreatedBy, &card.CreatedAt, &card.UpdatedAt, &card.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *GiftCardRepository) GetGiftCardByNumber(ctx context.Context, orgID uuid.UUID, cardNumber string) (*giftcards.GiftCard, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, card_number, pin_code, customer_id,
		       original_value, current_balance, issued_date, expiry_date,
		       status, issued_by_user_id, issued_location_id, notes,
		       created_by, created_at, updated_at, deleted_at
		FROM gift_cards
		WHERE organization_id = $1 AND card_number = $2 AND deleted_at IS NULL
	`

	var card giftcards.GiftCard
	err := r.db.Pool.QueryRow(ctx, query, orgID, cardNumber).Scan(
		&card.ID, &card.OrganizationID, &card.CardNumber, &card.PINCode, &card.CustomerID,
		&card.OriginalValue, &card.CurrentBalance, &card.IssuedDate, &card.ExpiryDate,
		&card.Status, &card.IssuedByUserID, &card.IssuedLocationID, &card.Notes,
		&card.CreatedBy, &card.CreatedAt, &card.UpdatedAt, &card.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *GiftCardRepository) ListGiftCards(ctx context.Context, orgID uuid.UUID, filters giftcards.GiftCardFilters) ([]giftcards.GiftCard, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, card_number, pin_code, customer_id,
		       original_value, current_balance, issued_date, expiry_date,
		       status, issued_by_user_id, issued_location_id, notes,
		       created_by, created_at, updated_at, deleted_at
		FROM gift_cards
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.Customer != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.Customer)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND card_number ILIKE $%d", argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []giftcards.GiftCard
	for rows.Next() {
		var card giftcards.GiftCard
		err := rows.Scan(
			&card.ID, &card.OrganizationID, &card.CardNumber, &card.PINCode, &card.CustomerID,
			&card.OriginalValue, &card.CurrentBalance, &card.IssuedDate, &card.ExpiryDate,
			&card.Status, &card.IssuedByUserID, &card.IssuedLocationID, &card.Notes,
			&card.CreatedBy, &card.CreatedAt, &card.UpdatedAt, &card.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, rows.Err()
}

func (r *GiftCardRepository) UpdateGiftCard(ctx context.Context, card *giftcards.GiftCard) error {
	if err := r.db.SetOrganizationContext(ctx, card.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE gift_cards
		SET current_balance = $1, status = $2, customer_id = $3, expiry_date = $4, updated_at = $5
		WHERE organization_id = $6 AND id = $7 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		card.CurrentBalance, card.Status, card.CustomerID, card.ExpiryDate, card.UpdatedAt,
		card.OrganizationID, card.ID,
	)
	return err
}

func (r *GiftCardRepository) DeleteGiftCard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE gift_cards
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// ============= GIFT CARD TRANSACTIONS =============

func (r *GiftCardRepository) CreateGiftCardTransaction(ctx context.Context, txn *giftcards.GiftCardTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, txn.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO gift_card_transactions (
			id, organization_id, gift_card_id, transaction_type, amount,
			balance_after, sale_id, payment_id, user_id, location_id,
			notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		txn.ID, txn.OrganizationID, txn.GiftCardID, txn.TransactionType, txn.Amount,
		txn.BalanceAfter, txn.SaleID, txn.PaymentID, txn.UserID, txn.LocationID,
		txn.Notes, txn.CreatedAt,
	)
	return err
}

func (r *GiftCardRepository) GetGiftCardTransaction(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*giftcards.GiftCardTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, gift_card_id, transaction_type, amount,
		       balance_after, sale_id, payment_id, user_id, location_id,
		       notes, created_at, deleted_at
		FROM gift_card_transactions
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var txn giftcards.GiftCardTransaction
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&txn.ID, &txn.OrganizationID, &txn.GiftCardID, &txn.TransactionType, &txn.Amount,
		&txn.BalanceAfter, &txn.SaleID, &txn.PaymentID, &txn.UserID, &txn.LocationID,
		&txn.Notes, &txn.CreatedAt, &txn.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *GiftCardRepository) ListGiftCardTransactions(ctx context.Context, orgID uuid.UUID, cardID uuid.UUID) ([]giftcards.GiftCardTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, gift_card_id, transaction_type, amount,
		       balance_after, sale_id, payment_id, user_id, location_id,
		       notes, created_at, deleted_at
		FROM gift_card_transactions
		WHERE organization_id = $1 AND gift_card_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []giftcards.GiftCardTransaction
	for rows.Next() {
		var txn giftcards.GiftCardTransaction
		err := rows.Scan(
			&txn.ID, &txn.OrganizationID, &txn.GiftCardID, &txn.TransactionType, &txn.Amount,
			&txn.BalanceAfter, &txn.SaleID, &txn.PaymentID, &txn.UserID, &txn.LocationID,
			&txn.Notes, &txn.CreatedAt, &txn.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}

	return txns, rows.Err()
}

// ============= STORE CREDIT =============

func (r *GiftCardRepository) CreateStoreCreditAccount(ctx context.Context, account *giftcards.StoreCreditAccount) error {
	if err := r.db.SetOrganizationContext(ctx, account.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO customer_store_credit_accounts (
			id, organization_id, customer_id, current_balance, credit_limit,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		account.ID, account.OrganizationID, account.CustomerID, account.CurrentBalance,
		account.CreditLimit, account.IsActive, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

func (r *GiftCardRepository) GetStoreCreditAccount(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*giftcards.StoreCreditAccount, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, customer_id, current_balance, credit_limit,
		       is_active, created_at, updated_at, deleted_at
		FROM customer_store_credit_accounts
		WHERE organization_id = $1 AND customer_id = $2 AND deleted_at IS NULL
	`

	var account giftcards.StoreCreditAccount
	err := r.db.Pool.QueryRow(ctx, query, orgID, customerID).Scan(
		&account.ID, &account.OrganizationID, &account.CustomerID, &account.CurrentBalance,
		&account.CreditLimit, &account.IsActive, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *GiftCardRepository) ListStoreCreditAccounts(ctx context.Context, orgID uuid.UUID, filters giftcards.StoreCreditFilters) ([]giftcards.StoreCreditAccount, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, customer_id, current_balance, credit_limit,
		       is_active, created_at, updated_at, deleted_at
		FROM customer_store_credit_accounts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.Customer != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.Customer)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []giftcards.StoreCreditAccount
	for rows.Next() {
		var account giftcards.StoreCreditAccount
		err := rows.Scan(
			&account.ID, &account.OrganizationID, &account.CustomerID, &account.CurrentBalance,
			&account.CreditLimit, &account.IsActive, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *GiftCardRepository) UpdateStoreCreditAccount(ctx context.Context, account *giftcards.StoreCreditAccount) error {
	if err := r.db.SetOrganizationContext(ctx, account.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customer_store_credit_accounts
		SET current_balance = $1, is_active = $2, updated_at = $3
		WHERE organization_id = $4 AND id = $5 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		account.CurrentBalance, account.IsActive, account.UpdatedAt,
		account.OrganizationID, account.ID,
	)
	return err
}

// ============= STORE CREDIT TRANSACTIONS =============

func (r *GiftCardRepository) CreateStoreCreditTransaction(ctx context.Context, txn *giftcards.StoreCreditTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, txn.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO store_credit_transactions (
			id, organization_id, store_credit_account_id, transaction_type, amount,
			balance_after, sale_id, payment_id, user_id, location_id, notes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		txn.ID, txn.OrganizationID, txn.StoreCreditAccountID, txn.TransactionType, txn.Amount,
		txn.BalanceAfter, txn.SaleID, txn.PaymentID, txn.UserID, txn.LocationID, txn.Notes, txn.CreatedAt,
	)
	return err
}

func (r *GiftCardRepository) ListStoreCreditTransactions(ctx context.Context, orgID uuid.UUID, accountID uuid.UUID) ([]giftcards.StoreCreditTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, store_credit_account_id, transaction_type, amount,
		       balance_after, sale_id, payment_id, user_id, location_id, notes, created_at, deleted_at
		FROM store_credit_transactions
		WHERE organization_id = $1 AND store_credit_account_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []giftcards.StoreCreditTransaction
	for rows.Next() {
		var txn giftcards.StoreCreditTransaction
		err := rows.Scan(
			&txn.ID, &txn.OrganizationID, &txn.StoreCreditAccountID, &txn.TransactionType, &txn.Amount,
			&txn.BalanceAfter, &txn.SaleID, &txn.PaymentID, &txn.UserID, &txn.LocationID, &txn.Notes, &txn.CreatedAt, &txn.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}

	return txns, rows.Err()
}

// ============= RETURN REASONS =============

func (r *GiftCardRepository) ListReturnReasons(ctx context.Context, orgID *uuid.UUID) ([]giftcards.ReturnReason, error) {
	query := `
		SELECT id, organization_id, reason_code, reason_name, requires_approval,
		       affects_inventory, is_restockable, is_active, display_order, created_at, deleted_at
		FROM return_reasons
		WHERE (organization_id = $1 OR organization_id IS NULL) AND deleted_at IS NULL
		ORDER BY display_order ASC, reason_name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reasons []giftcards.ReturnReason
	for rows.Next() {
		var reason giftcards.ReturnReason
		err := rows.Scan(
			&reason.ID, &reason.OrganizationID, &reason.ReasonCode, &reason.ReasonName,
			&reason.RequiresApproval, &reason.AffectsInventory, &reason.IsRestockable,
			&reason.IsActive, &reason.DisplayOrder, &reason.CreatedAt, &reason.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		reasons = append(reasons, reason)
	}

	return reasons, rows.Err()
}

func (r *GiftCardRepository) GetReturnReason(ctx context.Context, id uuid.UUID) (*giftcards.ReturnReason, error) {
	query := `
		SELECT id, organization_id, reason_code, reason_name, requires_approval,
		       affects_inventory, is_restockable, is_active, display_order, created_at, deleted_at
		FROM return_reasons
		WHERE id = $1 AND deleted_at IS NULL
	`

	var reason giftcards.ReturnReason
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&reason.ID, &reason.OrganizationID, &reason.ReasonCode, &reason.ReasonName,
		&reason.RequiresApproval, &reason.AffectsInventory, &reason.IsRestockable,
		&reason.IsActive, &reason.DisplayOrder, &reason.CreatedAt, &reason.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reason, nil
}

// ============= SALE RETURNS =============

func (r *GiftCardRepository) CreateSaleReturn(ctx context.Context, return_ *giftcards.SaleReturn) error {
	if err := r.db.SetOrganizationContext(ctx, return_.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO sale_returns (
			id, organization_id, return_number, original_sale_id, customer_id,
			location_id, user_id, return_date, total_amount, refund_amount,
			restocking_fee, refund_method, status, notes, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		return_.ID, return_.OrganizationID, return_.ReturnNumber, return_.OriginalSaleID, return_.CustomerID,
		return_.LocationID, return_.UserID, return_.ReturnDate, return_.TotalAmount, return_.RefundAmount,
		return_.RestockingFee, return_.RefundMethod, return_.Status, return_.Notes, return_.CreatedBy,
		return_.CreatedAt, return_.UpdatedAt,
	)
	return err
}

func (r *GiftCardRepository) GetSaleReturn(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*giftcards.SaleReturn, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, return_number, original_sale_id, customer_id,
		       location_id, user_id, return_date, total_amount, refund_amount,
		       restocking_fee, refund_method, status, approved_by, approved_at,
		       notes, created_by, created_at, updated_at, deleted_at
		FROM sale_returns
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var return_ giftcards.SaleReturn
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&return_.ID, &return_.OrganizationID, &return_.ReturnNumber, &return_.OriginalSaleID, &return_.CustomerID,
		&return_.LocationID, &return_.UserID, &return_.ReturnDate, &return_.TotalAmount, &return_.RefundAmount,
		&return_.RestockingFee, &return_.RefundMethod, &return_.Status, &return_.ApprovedBy, &return_.ApprovedAt,
		&return_.Notes, &return_.CreatedBy, &return_.CreatedAt, &return_.UpdatedAt, &return_.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &return_, nil
}

func (r *GiftCardRepository) GetSaleReturnByNumber(ctx context.Context, orgID uuid.UUID, returnNumber string) (*giftcards.SaleReturn, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, return_number, original_sale_id, customer_id,
		       location_id, user_id, return_date, total_amount, refund_amount,
		       restocking_fee, refund_method, status, approved_by, approved_at,
		       notes, created_by, created_at, updated_at, deleted_at
		FROM sale_returns
		WHERE organization_id = $1 AND return_number = $2 AND deleted_at IS NULL
	`

	var return_ giftcards.SaleReturn
	err := r.db.Pool.QueryRow(ctx, query, orgID, returnNumber).Scan(
		&return_.ID, &return_.OrganizationID, &return_.ReturnNumber, &return_.OriginalSaleID, &return_.CustomerID,
		&return_.LocationID, &return_.UserID, &return_.ReturnDate, &return_.TotalAmount, &return_.RefundAmount,
		&return_.RestockingFee, &return_.RefundMethod, &return_.Status, &return_.ApprovedBy, &return_.ApprovedAt,
		&return_.Notes, &return_.CreatedBy, &return_.CreatedAt, &return_.UpdatedAt, &return_.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &return_, nil
}

func (r *GiftCardRepository) ListSaleReturns(ctx context.Context, orgID uuid.UUID, filters giftcards.ReturnFilters) ([]giftcards.SaleReturn, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, return_number, original_sale_id, customer_id,
		       location_id, user_id, return_date, total_amount, refund_amount,
		       restocking_fee, refund_method, status, approved_by, approved_at,
		       notes, created_by, created_at, updated_at, deleted_at
		FROM sale_returns
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.Customer != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.Customer)
	}

	if filters.Location != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.Location)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND return_date >= $%d", argCount)
		args = append(args, *filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND return_date <= $%d", argCount)
		args = append(args, *filters.DateTo)
	}

	query += " ORDER BY return_date DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var returns []giftcards.SaleReturn
	for rows.Next() {
		var return_ giftcards.SaleReturn
		err := rows.Scan(
			&return_.ID, &return_.OrganizationID, &return_.ReturnNumber, &return_.OriginalSaleID, &return_.CustomerID,
			&return_.LocationID, &return_.UserID, &return_.ReturnDate, &return_.TotalAmount, &return_.RefundAmount,
			&return_.RestockingFee, &return_.RefundMethod, &return_.Status, &return_.ApprovedBy, &return_.ApprovedAt,
			&return_.Notes, &return_.CreatedBy, &return_.CreatedAt, &return_.UpdatedAt, &return_.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		returns = append(returns, return_)
	}

	return returns, rows.Err()
}

func (r *GiftCardRepository) UpdateSaleReturn(ctx context.Context, return_ *giftcards.SaleReturn) error {
	if err := r.db.SetOrganizationContext(ctx, return_.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE sale_returns
		SET status = $1, approved_by = $2, approved_at = $3, total_amount = $4,
		    refund_amount = $5, refund_method = $6, updated_at = $7
		WHERE organization_id = $8 AND id = $9 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		return_.Status, return_.ApprovedBy, return_.ApprovedAt, return_.TotalAmount,
		return_.RefundAmount, return_.RefundMethod, return_.UpdatedAt,
		return_.OrganizationID, return_.ID,
	)
	return err
}

// ============= SALE RETURN ITEMS =============

func (r *GiftCardRepository) CreateSaleReturnItem(ctx context.Context, item *giftcards.SaleReturnItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO sale_return_items (
			id, organization_id, sale_return_id, original_sale_item_id, product_id,
			product_variant_id, quantity, unit_price, subtotal, tax_amount,
			discount_amount, total_amount, return_reason_id, return_reason_notes,
			item_condition, is_restockable, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.OrganizationID, item.SaleReturnID, item.OriginalSaleItemID, item.ProductID,
		item.ProductVariantID, item.Quantity, item.UnitPrice, item.Subtotal, item.TaxAmount,
		item.DiscountAmount, item.TotalAmount, item.ReturnReasonID, item.ReturnReasonNotes,
		item.ItemCondition, item.IsRestockable, item.CreatedAt,
	)
	return err
}

func (r *GiftCardRepository) ListSaleReturnItems(ctx context.Context, orgID uuid.UUID, returnID uuid.UUID) ([]giftcards.SaleReturnItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_return_id, original_sale_item_id, product_id,
		       product_variant_id, quantity, unit_price, subtotal, tax_amount,
		       discount_amount, total_amount, return_reason_id, return_reason_notes,
		       item_condition, is_restockable, created_at, deleted_at
		FROM sale_return_items
		WHERE organization_id = $1 AND sale_return_id = $2 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, returnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []giftcards.SaleReturnItem
	for rows.Next() {
		var item giftcards.SaleReturnItem
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.SaleReturnID, &item.OriginalSaleItemID, &item.ProductID,
			&item.ProductVariantID, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.TaxAmount,
			&item.DiscountAmount, &item.TotalAmount, &item.ReturnReasonID, &item.ReturnReasonNotes,
			&item.ItemCondition, &item.IsRestockable, &item.CreatedAt, &item.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
