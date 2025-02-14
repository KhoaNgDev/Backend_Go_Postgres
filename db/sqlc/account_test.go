package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/KhoaNgDev/Backend_Go_Postgres/db/util"
	"github.com/stretchr/testify/require"
)

//? Tạo một tài khoản ngẫu nhiên để sử dụng trong các test case khác
func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),    // Tạo tên chủ tài khoản ngẫu nhiên
		Balance:  util.RandomMoney(),    // Số dư ngẫu nhiên
		Currency: util.RandomCurrency(), // Loại tiền tệ ngẫu nhiên
	}

	//! Gọi function tạo tài khoản trong database
	account, err := testQueries.CreateAccount(context.Background(), arg)

	//TODO: Đảm bảo không có lỗi khi tạo tài khoản
	require.NoError(t, err)
	require.NotEmpty(t, account)

	//TODO: Kiểm tra dữ liệu tài khoản có đúng không
	checkAccountEqual(t, arg, account)

	//TODO: Đảm bảo ID và thời gian tạo tài khoản hợp lệ
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

//? Hàm kiểm tra xem dữ liệu tài khoản có khớp với thông tin đầu vào hay không
func checkAccountEqual(t *testing.T, arg CreateAccountParams, account Account) {
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
}

//? Test kiểm tra việc tạo tài khoản
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t) // Nếu có lỗi, test sẽ fail
}

//? Test kiểm tra lấy thông tin tài khoản từ database
func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	//! Đảm bảo không có lỗi khi lấy tài khoản
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//TODO: Kiểm tra dữ liệu có khớp với tài khoản đã tạo không
	require.EqualValues(t, account1, account2)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

//? Test kiểm tra cập nhật số dư tài khoản
func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(), // Tạo số dư mới ngẫu nhiên
	}

	//! Gọi function cập nhật tài khoản
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//TODO: Kiểm tra dữ liệu có khớp không
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

//? Test kiểm tra xóa tài khoản
func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	//! Gọi function xóa tài khoản
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	//TODO: Thử lấy lại tài khoản đã xóa
	account2, err := testQueries.GetAccount(context.Background(), account.ID)

	//! Đảm bảo tài khoản không còn tồn tại
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

//? Test kiểm tra danh sách tài khoản
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5, // Bỏ qua 5 tài khoản đầu tiên, lấy 5 tài khoản tiếp theo
	}

	//! Lấy danh sách tài khoản từ database
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	//TODO: Kiểm tra từng tài khoản trong danh sách có hợp lệ không
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}