package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	transaction := NewTransaction(testDB)
	senderAccount := createRandomAccount(t)
	recieverAccount := createRandomAccount(t)
	fmt.Println(">> before balance: ", senderAccount.Balance, recieverAccount.Balance)

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTXResults)
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := transaction.TransferTX(ctx, TransferTXParams{
				FromAccountID: senderAccount.ID,
				ToAccountID:   recieverAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	exist := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, senderAccount.ID, transfer.FromAccountID)
		require.Equal(t, recieverAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = transaction.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, senderAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = transaction.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, recieverAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = transaction.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		from_account := result.FromAccount
		require.NotEmpty(t, from_account)
		require.Equal(t, from_account.ID, senderAccount.ID)

		to_account := result.ToAccount
		require.NotEmpty(t, to_account)
		require.Equal(t, to_account.ID, recieverAccount.ID)

		fmt.Println(">> tx: ", from_account.Balance, to_account.Balance)

		senderDiff := senderAccount.Balance - from_account.Balance
		reciverDiff := to_account.Balance - recieverAccount.Balance
		require.Equal(t, senderDiff, reciverDiff)
		require.True(t, senderDiff > 0)
		require.True(t, senderDiff%amount == 0)

		k := int(senderDiff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exist, k)
		exist[k] = true
	}

	senderUpdatedAccount, err := transaction.GetAccount(context.Background(), senderAccount.ID)
	require.NoError(t, err)
	require.Equal(t, senderAccount.ID, senderUpdatedAccount.ID)

	recieverUpdatedAccount, err := transaction.GetAccount(context.Background(), recieverAccount.ID)
	require.NoError(t, err)
	require.Equal(t, recieverUpdatedAccount.ID, recieverAccount.ID)

	fmt.Println(">> after balance: ", senderUpdatedAccount.Balance, recieverUpdatedAccount.Balance)

	require.Equal(t, senderUpdatedAccount.Balance, senderAccount.Balance-int64(n)*amount)
	require.Equal(t, recieverUpdatedAccount.Balance, recieverAccount.Balance+int64(n)*amount)
}

func TestTransferTXDeadLock(t *testing.T) {
	transaction := NewTransaction(testDB)
	senderAccount := createRandomAccount(t)
	recieverAccount := createRandomAccount(t)
	fmt.Println(">> before balance: ", senderAccount.Balance, recieverAccount.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		senderAccountID := senderAccount.ID
		recieverAccountID := recieverAccount.ID

		if i%2 == 1 {
			senderAccountID = recieverAccount.ID
			recieverAccountID = senderAccount.ID
		}

		go func() {
			ctx := context.Background()
			_, err := transaction.TransferTX(ctx, TransferTXParams{
				FromAccountID: senderAccountID,
				ToAccountID:   recieverAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	senderUpdatedAccount, err := transaction.GetAccount(context.Background(), senderAccount.ID)
	require.NoError(t, err)
	require.Equal(t, senderAccount.ID, senderUpdatedAccount.ID)

	recieverUpdatedAccount, err := transaction.GetAccount(context.Background(), recieverAccount.ID)
	require.NoError(t, err)
	require.Equal(t, recieverUpdatedAccount.ID, recieverAccount.ID)

	fmt.Println(">> after balance: ", senderUpdatedAccount.Balance, recieverUpdatedAccount.Balance)

	require.Equal(t, senderUpdatedAccount.Balance, senderAccount.Balance)
	require.Equal(t, recieverUpdatedAccount.Balance, recieverAccount.Balance)
}
