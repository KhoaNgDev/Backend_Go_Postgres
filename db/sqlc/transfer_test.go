package db

import (
	"context"
	"testing"
	"time"

	"github.com/KhoaNgDev/Backend_Go_Postgres/db/util"
	"github.com/stretchr/testify/require"
)

// ? createRandomTransfer tạo một giao dịch chuyển tiền ngẫu nhiên giữa hai tài khoản
func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(), // * Chọn số tiền ngẫu nhiên
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)       // ! Đảm bảo không có lỗi khi tạo giao dịch
	require.NotEmpty(t, transfer) // ! Kiểm tra giao dịch không rỗng

	// ? Kiểm tra thông tin giao dịch có khớp với tham số đầu vào không
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)        // ! Đảm bảo ID giao dịch đã được tạo
	require.NotZero(t, transfer.CreatedAt) // ! Đảm bảo thời gian tạo không rỗng

	return transfer
}

// ? TestCreateTransfer kiểm tra chức năng tạo giao dịch
func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t) // * Tạo tài khoản 1
	account2 := createRandomAccount(t) // * Tạo tài khoản 2
	createRandomTransfer(t, account1, account2)
}

// ? TestGetTransfer kiểm tra chức năng lấy thông tin một giao dịch
func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2) // * Tạo giao dịch để kiểm tra

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)       // ! Kiểm tra không có lỗi xảy ra
	require.NotEmpty(t, transfer2) // ! Đảm bảo giao dịch tồn tại

	// ? Kiểm tra thông tin giao dịch lấy ra có đúng với dữ liệu gốc không
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second) // ! Đảm bảo thời gian gần nhau
}

// ? TestListTransfer kiểm tra chức năng lấy danh sách giao dịch
func TestListTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2) // * Tạo giao dịch từ account1 -> account2
		createRandomTransfer(t, account2, account1) // * Tạo giao dịch từ account2 -> account1
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID, // * Lọc giao dịch từ account1
		ToAccountID:   account1.ID, // * Lọc giao dịch đến account1
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)       // ! Đảm bảo không có lỗi khi lấy danh sách
	require.Len(t, transfers, 5)  // ! Kiểm tra danh sách có đúng số lượng

	// ? Duyệt danh sách và kiểm tra từng giao dịch
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer) // ! Đảm bảo giao dịch không rỗng
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID) // ! Kiểm tra giao dịch liên quan đến account1
	}
}