package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qweartyuiopasdfghjklzxcvbnm"
var randomGen = rand.New(rand.NewSource(time.Now().UnixNano()))
var firstNames = []string{"James", "John", "Robert", "Michael", "William", "David", "Richard", "Joseph", "Thomas", "Charles", "Daniel", "Matthew", "Anthony", "Mark", "Elizabeth", "Jennifer", "Emily", "Sarah", "Jessica", "Ashley"}
var lastNames = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Miller", "Davis", "Garcia", "Rodriguez", "Wilson", "Martinez", "Anderson", "Taylor", "Thomas", "Hernandez", "Moore", "Martin", "Jackson", "Thompson", "White"}

func RandomInt(min, max int64) int64 {
	return min + randomGen.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder 
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[randomGen.Intn(k)] 
		sb.WriteByte(c)                  
	}
	return sb.String()
}

func RandomOwner() string {
	first := firstNames[randomGen.Intn(len(firstNames))]
	last := lastNames[randomGen.Intn(len(lastNames))]
	return first + " " + last 
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "VNĐ"} 
	n := len(currencies)
	return currencies[randomGen.Intn(n)] 
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
