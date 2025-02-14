package db

import (
	"context"          
	"testing"        
	"github.com/stretchr/testify/require" 
)


func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Kane",  
		Balance:  10000,   // Số dư
		Currency: "VNĐ",   // Loại tiền tệ 
	}

	// Gọi function CreateAccount để thêm tài khoản mới vào database.
	accounts, err := testQueries.CreateAccount(context.Background(), arg)

	// Kiểm tra lỗi:
	require.NoError(t, err)   // Đảm bảo rằng không có lỗi khi tạo tài khoản.
	require.NotEmpty(t, accounts) // Đảm bảo rằng tài khoản được tạo không phải là rỗng.

	// Kiểm tra dữ liệu trả về có chính xác không:
	require.Equal(t, arg.Owner, accounts.Owner)       // So sánh chủ tài khoản.
	require.Equal(t, arg.Balance, accounts.Balance)   // So sánh số dư.
	require.Equal(t, arg.Currency, accounts.Currency) // So sánh loại tiền tệ.

	// Kiểm tra các trường quan trọng:
	require.NotZero(t, accounts.ID)        // Đảm bảo rằng ID được gán (có nghĩa là tài khoản đã được lưu vào DB).
	require.NotZero(t, accounts.CreatedAt) // Đảm bảo rằng thời gian tạo tài khoản được ghi nhận.

	// 🛠 **Cách debug khi gặp lỗi:**
	// - Nếu gặp lỗi "relation 'accounts' does not exist" → Kiểm tra lại migration DB.
	// - Nếu ID hoặc CreatedAt bằng 0 → Có thể lỗi do không insert đúng vào DB.
	// - Nếu Balance, Owner hoặc Currency sai → Kiểm tra logic trong CreateAccount.
}
