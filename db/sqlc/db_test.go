package db

import (
    "testing"
)

func TestDBConnection(t *testing.T) {
    if testQueries == nil {
        t.Fatal("testQueries is nil, database connection failed")
    }
}
