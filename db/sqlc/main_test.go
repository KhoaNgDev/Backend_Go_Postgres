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


