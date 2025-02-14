package util

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// 🧪 TestRandomInt kiểm tra việc tạo số ngẫu nhiên trong khoảng min-max
func TestRandomInt(t *testing.T) {
	min, max := int64(1), int64(100)

	t.Run("✅ Giá trị trong khoảng", func(t *testing.T) {
		for i := 0; i < 10; i++ { // Lặp 10 lần để kiểm tra tính ngẫu nhiên
			result := RandomInt(min, max)
			require.GreaterOrEqual(t, result, min) // Phải >= min
			require.LessOrEqual(t, result, max)    // Phải <= max
		}
	})
}

// 🧪 TestRandomString kiểm tra việc tạo chuỗi ngẫu nhiên
func TestRandomString(t *testing.T) {
	length := 10

	t.Run("✅ Độ dài chính xác", func(t *testing.T) {
		result := RandomString(length)
		require.Len(t, result, length) // Chuỗi phải đúng độ dài yêu cầu
	})

	t.Run("✅ Chỉ chứa ký tự hợp lệ", func(t *testing.T) {
		result := RandomString(length)
		for _, c := range result {
			require.Contains(t, alphabet, string(c)) // Kiểm tra từng ký tự có trong alphabet không
		}
	})
}

// 🧪 TestRandomOwner kiểm tra việc tạo tên người ngẫu nhiên
func TestRandomOwner(t *testing.T) {
	t.Run("✅ Không được rỗng", func(t *testing.T) {
		result := RandomOwner()
		require.NotEmpty(t, result) // Tên không được rỗng
	})

	t.Run("✅ Định dạng đúng 'First Last'", func(t *testing.T) {
		result := RandomOwner()
		parts := strings.Split(result, " ")
		require.Len(t, parts, 2) // Kiểm tra có đúng 2 phần (First Last) không
	})
}

// 🧪 TestRandomCurrency kiểm tra việc chọn tiền tệ ngẫu nhiên
func TestRandomCurrency(t *testing.T) {
	expectedCurrencies := []string{"EUR", "USD", "VNĐ"}

	t.Run("✅ Phải thuộc danh sách tiền tệ hợp lệ", func(t *testing.T) {
		result := RandomCurrency()
		require.Contains(t, expectedCurrencies, result) // Kết quả phải là 1 trong các loại tiền tệ hợp lệ
	})
}

// 🧪 TestRandomMoney kiểm tra số tiền ngẫu nhiên
func TestRandomMoney(t *testing.T) {
	t.Run("✅ Giá trị trong khoảng 0 - 1000", func(t *testing.T) {
		result := RandomMoney()
		require.GreaterOrEqual(t, result, int64(0)) // Phải >= 0
		require.LessOrEqual(t, result, int64(1000)) // Phải <= 1000
	})
}
