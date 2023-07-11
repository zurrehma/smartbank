package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zurrehma/bank-accounts/util"
)

func createRandomAccount(t *testing.T) Account {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	match_account, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, match_account)

	require.Equal(t, match_account.ID, account.ID)
	require.Equal(t, match_account.Owner, account.Owner)
	require.Equal(t, match_account.Currency, account.Currency)
	require.Equal(t, match_account.Balance, account.Balance)
	require.WithinDuration(t, match_account.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}

	updated_account, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updated_account)

	require.Equal(t, updated_account.ID, account.ID)
	require.Equal(t, updated_account.Owner, account.Owner)
	require.Equal(t, updated_account.Currency, account.Currency)
	require.Equal(t, args.Balance, updated_account.Balance)
	require.WithinDuration(t, updated_account.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	// fmt.Printf("created account: %v", account)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	match_account, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, match_account)
}

func TestListAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, accounts := range accounts {
		require.NotEmpty(t, accounts)
	}

}
