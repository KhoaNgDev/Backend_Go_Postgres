package db

import (
	"context"
	"database/sql"
	"fmt"
)

//! Store cung cấp tất cả các function để thực thi truy vấn DB và giao dịch
type Store struct {
	db *sql.DB
	*Queries
}

//! NewStore - Hàm khởi tạo Store
//? Nhận vào `db *sql.DB`, tạo một Store mới có thể thực thi truy vấn và giao dịch
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//! execTx - Thực thi một function trong phạm vi giao dịch (transaction)
//? Nếu xảy ra lỗi, giao dịch sẽ bị rollback
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Bắt đầu transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Tạo instance Queries từ transaction
	q := New(tx)

	// Thực thi function `fn`
	err = fn(q)
	if err != nil {
		// Rollback nếu có lỗi
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction nếu thành công
	return tx.Commit()
}

//! TransferTxParams - Struct chứa dữ liệu đầu vào của giao dịch chuyển tiền
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"` // ID tài khoản nguồn
	ToAccountID   int64 `json:"to_account_id"`   // ID tài khoản đích
	Amount        int64 `json:"amount"`          // Số tiền cần chuyển
}

//! TransferTxResult - Struct chứa kết quả của giao dịch chuyển tiền
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`    // Bản ghi về giao dịch chuyển tiền
	FromAccount Account  `json:"from_account"` // Tài khoản nguồn sau khi cập nhật số dư
	ToAccount   Account  `json:"to_account"`   // Tài khoản đích sau khi cập nhật số dư
	FromEntry   Entry    `json:"from_entry"`   // Bản ghi trừ tiền tài khoản nguồn
	ToEntry     Entry    `json:"to_entry"`     // Bản ghi cộng tiền tài khoản đích
}

//! TransferTx - Thực hiện giao dịch chuyển tiền giữa 2 tài khoản
//? Gồm các bước: tạo bản ghi giao dịch, cập nhật tài khoản, và cập nhật số dư trong transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Thực thi transaction
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//? 1. Tạo bản ghi giao dịch chuyển tiền
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//? 2. Tạo bản ghi trừ tiền từ tài khoản nguồn
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Trừ tiền
		})
		if err != nil {
			return err
		}

		//? 3. Tạo bản ghi cộng tiền vào tài khoản đích
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // Cộng tiền
		})
		if err != nil {
			return err
		}

		//? 4. Cập nhật số dư tài khoản nguồn và tài khoản đích theo thứ tự ID
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

//! addMoney - Cập nhật số dư cho 2 tài khoản trong transaction
//? accountID1: tài khoản đầu tiên, amount1: số tiền cập nhật cho tài khoản đầu tiên
//? accountID2: tài khoản thứ hai, amount2: số tiền cập nhật cho tài khoản thứ hai
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	// Cập nhật số dư tài khoản đầu tiên
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	// Cập nhật số dư tài khoản thứ hai
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
