package util

import (
	"math/rand"
	"strings"
	"time"
)

//! Khởi tạo hạt giống cho bộ tạo số ngẫu nhiên
var randomGen = rand.New(rand.NewSource(time.Now().UnixNano()))

//! Danh sách các ký tự dùng để tạo chuỗi ngẫu nhiên
const alphabet = "abcdefghijklmnopqrstuvwxyz"

//! Danh sách tên và họ ngẫu nhiên
var firstNames = []string{
	"James", "John", "Robert", "Michael", "William",
	"David", "Richard", "Joseph", "Thomas", "Charles",
	"Daniel", "Matthew", "Anthony", "Mark", "Elizabeth",
	"Jennifer", "Emily", "Sarah", "Jessica", "Ashley",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones",
	"Miller", "Davis", "Garcia", "Rodriguez", "Wilson",
	"Martinez", "Anderson", "Taylor", "Thomas", "Hernandez",
	"Moore", "Martin", "Jackson", "Thompson", "White",
}

//! RandomInt - Tạo số nguyên ngẫu nhiên trong khoảng [min, max]
func RandomInt(min, max int64) int64 {
	return min + randomGen.Int63n(max-min+1)
}

//! RandomString - Tạo chuỗi ngẫu nhiên với độ dài `n`
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	// Lặp `n` lần để chọn ký tự ngẫu nhiên từ `alphabet`
	for i := 0; i < n; i++ {
		c := alphabet[randomGen.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//! RandomOwner - Tạo tên ngẫu nhiên từ danh sách firstNames và lastNames
func RandomOwner() string {
	first := firstNames[randomGen.Intn(len(firstNames))]
	last := lastNames[randomGen.Intn(len(lastNames))]
	return first + " " + last
}

//! RandomCurrency - Chọn đơn vị tiền tệ ngẫu nhiên
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "VNĐ"}
	return currencies[randomGen.Intn(len(currencies))]
}

//! RandomMoney - Tạo số tiền ngẫu nhiên trong khoảng 0 - 1000
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
