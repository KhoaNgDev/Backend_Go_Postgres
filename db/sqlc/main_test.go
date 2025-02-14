package db

import (
	"database/sql"  
	"log"           
	"os"            
	"testing"      
	_ "github.com/lib/pq" // Import driver PostgreSQL nhưng không dùng trực tiếp
)

//! Cấu hình thông tin kết nối database
const (
	dbDriver = "postgres"  
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"  
)

//! Biến toàn cục để sử dụng chung trong các test case
var testQueries *Queries
var testDB *sql.DB

//! Hàm khởi tạo trước khi chạy các test case
//? Hàm này sẽ được chạy **trước** khi tất cả các test trong package `db` được thực thi.
func TestMain(m *testing.M) {
	var err error

	//? 1. Kết nối đến database
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("❌ Không thể kết nối đến database: ", err)
	}

	//? 2. Kiểm tra xem database có thể nhận truy vấn hay không
	err = testDB.Ping()
	if err != nil {
		log.Fatal("❌ Database không thể truy cập: ", err)
	}

	//? 3. Gán kết nối DB vào `testQueries` để sử dụng trong các test case
	testQueries = New(testDB)

	//? 4. Chạy tất cả test case và trả về mã lỗi phù hợp
	os.Exit(m.Run())
}
