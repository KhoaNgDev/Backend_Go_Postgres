package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all funcs to exec db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// (Hàm này thực thi một function trong một transaction của database)
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// 1️⃣ Bắt đầu một transaction mới
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err // Nếu không thể mở transaction, trả về lỗi ngay
	}

	// 2️⃣ Tạo một instance của Queries, sử dụng transaction thay vì db trực tiếp
	q := New(tx)

	// 3️⃣ Gọi function `fn`, thực thi các truy vấn database trong transaction
	err = fn(q)
	if err != nil {
		// 4️⃣ Nếu có lỗi, rollback transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			// Nếu rollback cũng bị lỗi, báo lỗi kép
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// Trả về lỗi từ function `fn`
		return err
	}

	// 5️⃣ Nếu không có lỗi, commit transaction để lưu thay đổi vào database
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single db transaction

// creates a transfer record
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update account's balance

		return nil
	})
	return result, err
}
