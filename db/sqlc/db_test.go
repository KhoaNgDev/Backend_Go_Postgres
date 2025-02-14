package db

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestDBConnection(t *testing.T) {
    t.Run("✅ Kiểm tra testDB", func(t *testing.T) {
        require.NotNil(t, testDB, "❌ Kết nối tới database thất bại")
        
        // Thử ping để đảm bảo DB có phản hồi
        err := testDB.Ping()
        require.NoError(t, err, "❌ Không thể ping database, có thể kết nối bị lỗi")
    })

    t.Run("✅ Kiểm tra testQueries", func(t *testing.T) {
        require.NotNil(t, testQueries, "❌Không thể thực hiện do đối tượng truy vấn chưa được khởi tạo")

        // Thử truy vấn giả để kiểm tra DB có hoạt động không
        _, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
            Limit:  1,
            Offset: 0,
        })
        require.NoError(t, err, "❌ Không thể thực hiện truy vấn mẫu, có thể migration chưa chạy")
    })
}
