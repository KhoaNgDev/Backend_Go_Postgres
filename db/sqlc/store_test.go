package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// ? TestTransferTx kiểm tra tính đúng đắn của giao dịch chuyển tiền giữa hai tài khoản
func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	// * Tạo hai tài khoản ngẫu nhiên để kiểm tra
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5          // ! Số lượng giao dịch song song
	amount := int64(10) // ! Số tiền chuyển mỗi giao dịch

	errChan := make(chan error) // * Kênh để nhận lỗi từ các goroutine
	resultChan := make(chan TransferTxResult) // * Kênh để nhận kết quả từ các giao dịch

	// TODO: Thực hiện n giao dịch chuyển khoản đồng thời
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errChan <- err
			resultChan <- result
		}()
	}

	existed := make(map[int]bool) // * Lưu trữ số lần chuyển tiền duy nhất để tránh trùng lặp

	// TODO: Xử lý kết quả của các giao dịch
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err) // ! Đảm bảo không có lỗi xảy ra

		result := <-resultChan
		require.NotEmpty(t, result) // ! Kiểm tra kết quả không rỗng

		// ? Kiểm tra thông tin giao dịch
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// * Kiểm tra dữ liệu trong bảng Transfer
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// ? Kiểm tra Entries của tài khoản gửi
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// ? Kiểm tra Entries của tài khoản nhận
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: Kiểm tra số dư tài khoản sau giao dịch
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		// * Xác thực tính toán số dư
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// TODO: Kiểm tra số dư cuối cùng của tài khoản
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

// ? TestTransferTxDeadlock kiểm tra deadlock khi thực hiện nhiều giao dịch giữa 2 tài khoản theo chiều ngược nhau
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10         // ! Số lượng giao dịch đồng thời
	amount := int64(10)
	errChan := make(chan error)

	// TODO: Thực hiện các giao dịch đồng thời với hướng ngược nhau để kiểm tra deadlock
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errChan <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err) // ! Kiểm tra không có lỗi xảy ra
	}

	// TODO: Kiểm tra số dư cuối cùng
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
