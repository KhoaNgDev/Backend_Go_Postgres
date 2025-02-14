package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Kane",
		Balance:  10000,
		Currency: "VNĐ",
	}
	accounts, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	require.Equal(t, arg.Owner, accounts.Owner)
	require.Equal(t, arg.Balance, accounts.Balance)
	require.Equal(t, arg.Currency, accounts.Currency)

	require.NotZero(t, accounts.ID)
	require.NotZero(t, accounts.CreatedAt)
}
