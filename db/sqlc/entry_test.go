package db

import (
	"context"
	"testing"
	"time"

	"github.com/KhoaNgDev/Backend_Go_Postgres/db/util"
	"github.com/stretchr/testify/require"
)

//! Tạo một entry (giao dịch) ngẫu nhiên trong database
//? Mỗi entry liên kết với một tài khoản cụ thể và có một số tiền ngẫu nhiên
func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,    // Liên kết entry với tài khoản đã cho
		Amount:    util.RandomMoney(), // Sinh số tiền ngẫu nhiên
	}

	//? Gọi hàm để tạo entry trong database
	entry, err := testQueries.CreateEntry(context.Background(), arg)

	//? Kiểm tra xem entry có được tạo thành công không
	require.NoError(t, err)       // Đảm bảo không có lỗi
	require.NotEmpty(t, entry)    // Entry không được rỗng

	//? Xác nhận dữ liệu entry có khớp với dữ liệu đầu vào không
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	//? Đảm bảo entry có ID hợp lệ và thời gian tạo không bị rỗng
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

//! Kiểm tra việc tạo entry có hoạt động đúng không
func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t) // Tạo tài khoản mới
	createRandomEntry(t, account)     // Tạo entry cho tài khoản đó
}

//! Kiểm tra lấy thông tin một entry từ database
func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t) // Tạo tài khoản mới
	entry1 := createRandomEntry(t, account) // Tạo entry mới

	//? Lấy entry từ database theo ID
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	//? Kiểm tra entry có tồn tại và không có lỗi khi truy vấn
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	//? Kiểm tra dữ liệu entry có khớp không
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)

	//? Kiểm tra thời gian tạo của entry có hợp lệ (chênh lệch tối đa 1s)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

//! Kiểm tra lấy danh sách các entry từ database
func TestListEntries(t *testing.T) {
	account := createRandomAccount(t) // Tạo tài khoản mới

	//? Tạo nhiều entry để kiểm tra việc lấy danh sách
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID, // Lấy entry của tài khoản này
		Limit:     5,          // Giới hạn số lượng kết quả trả về
		Offset:    0,          // Không bỏ qua entry nào
	}

	//? Lấy danh sách entry từ database
	entries, err := testQueries.ListEntries(context.Background(), arg)

	//? Kiểm tra danh sách entry có hợp lệ không
	require.NoError(t, err)
	require.Len(t, entries, 5) // Chắc chắn chỉ lấy đúng 5 entry

	//? Kiểm tra từng entry trong danh sách có hợp lệ không
	for _, entry := range entries {
		require.NotEmpty(t, entry) // Entry không được rỗng
		require.Equal(t, account.ID, entry.AccountID) // Entry phải thuộc về tài khoản đã tạo
	}
}
