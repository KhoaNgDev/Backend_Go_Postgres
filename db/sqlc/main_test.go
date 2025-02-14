package db

import (
	"database/sql"  
	"log"           
	"os"            
	"testing"      
	_ "github.com/lib/pq" 
)


const (
	dbDriver = "postgres"  
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"  
)

// Biến toàn cục để dùng chung trong các test case.
var testQueries *Queries

func TestMain(m *testing.M) {
	// Kết nối đến database.
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	// Kiểm tra xem database có thể thực hiện truy vấn không.
	err = conn.Ping()
	if err != nil {
		log.Fatal("Database is not reachable: ", err)
	}

	// Gán kết nối DB vào testQueries để dùng cho các test case.
	testQueries = New(conn)

	// Chạy toàn bộ test case và thoát chương trình với mã lỗi thích hợp.
	os.Exit(m.Run())
}


/*
 **Cách debug khi gặp lỗi**:
1. **Lỗi "cannot connect to db"**:
   - Kiểm tra xem PostgreSQL có đang chạy không: `docker ps`
   - Kiểm tra kết nối DB với lệnh: `psql -U root -d simple_bank`
   - Kiểm tra `dbSource` có đúng với config của database không.
   - Nếu dùng Docker, đảm bảo container đang chạy: `docker start postgres12`.

2. **Lỗi "relation does not exist" khi chạy test**:
   - Kiểm tra lại xem đã chạy migration chưa (`make migrateup`).
   - Kiểm tra thư mục `db/migration/` có chứa các file migration hợp lệ không.

3. **Test bị panic hoặc bị lỗi kết nối giữa chừng**:
   - Xem có phải do DB bị reset hoặc container bị stop giữa test không.
   - Kiểm tra log PostgreSQL bằng: `docker logs postgres12`.

4. **Lỗi testQueries bị nil hoặc không khởi tạo**:
   - Đảm bảo `sql.Open()` không bị lỗi trước khi gọi `New(conn)`.
   - Kiểm tra lại `dbDriver` và `dbSource`.

*/
