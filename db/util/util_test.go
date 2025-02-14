package util

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min, max := int64(1), int64(100)
	result := RandomInt(min, max)
	require.GreaterOrEqual(t, result, min)
	require.LessOrEqual(t, result, max)
}

func TestRandomString(t *testing.T) {
	length := 10
	result := RandomString(length)
	require.Len(t, result, length)
}

func TestRandomOwner(t *testing.T) {
	result := RandomOwner()
	require.NotEmpty(t, result)
}

func TestRandomCurrency(t *testing.T) {
	result := RandomCurrency()
	expectedCurrencies := []string{"EUR", "USD", "VNĐ"}
	require.Contains(t, expectedCurrencies, result)
}

func TestRandomMoney(t *testing.T) {
	result := RandomMoney()
	require.GreaterOrEqual(t, result, int64(0))
	require.LessOrEqual(t, result, int64(1000))
}
