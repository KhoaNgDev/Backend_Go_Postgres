package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Log
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	nTransfers := 2
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < nTransfers; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < nTransfers; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check Entries

		// FROM
		fromEntry := result.FromEntry

		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// TO
		toEntry := result.ToEntry

		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Kiểm tra tài khoản nguồn (FromAccount)
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)              // Đảm bảo tài khoản nguồn không rỗng
		require.Equal(t, account1.ID, fromAccount.ID) // Xác nhận ID tài khoản nguồn khớp với tài khoản ban đầu

		// Kiểm tra tài khoản đích (ToAccount)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)              // Đảm bảo tài khoản đích không rỗng
		require.Equal(t, account2.ID, toAccount.ID) // Xác nhận ID tài khoản đích khớp với tài khoản ban đầu

		// Kiểm tra số dư tài khoản sau giao dịch
		// Log
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		
		diff1 := account1.Balance - fromAccount.Balance // Sự thay đổi số dư của tài khoản nguồn
		diff2 := toAccount.Balance - account2.Balance   // Sự thay đổi số dư của tài khoản đích

		require.Equal(t, diff1, diff2)     // Đảm bảo số tiền rút từ tài khoản nguồn bằng số tiền cộng vào tài khoản đích
		require.True(t, diff1 > 0)         // Đảm bảo số tiền bị trừ là số dương (tức là tài khoản nguồn đã bị giảm tiền)
		require.True(t, diff1%amount == 0) // Đảm bảo sự thay đổi số dư là bội số của `amount` (1 * amount, 2 * amount, ..., n * amount)

		k := int(diff1 / amount)                   // Xác định số lần giao dịch được thực hiện
		require.True(t, k >= 1 && k <= nTransfers) // Đảm bảo số lần giao dịch nằm trong phạm vi hợp lệ (từ 1 đến n)
		require.NotContains(t, existed, k)
		existed[k] = true

		// Check the final updated balances
		updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
		require.NoError(t, err)

		updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
		require.NoError(t, err)

		// Log
		fmt.Println(">> after:", updateAccount1.Balance, updateAccount2.Balance)

		require.Equal(t, account1.Balance-int64(nTransfers)*amount, updateAccount1.Balance)
		require.Equal(t, account2.Balance+int64(nTransfers)*amount, updateAccount2.Balance)

	}
}
