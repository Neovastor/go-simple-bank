package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store := NewStore(testDb)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	//run 5 transfers concurrently
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		var result TransferTxResult = <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		//check transfer data
		transfer := result.Transfer

		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		getTransfer, _ := testQueries.GetTransfer(context.Background(), transfer.ID)

		require.NotEmpty(t, getTransfer)

		//check ToEntry data
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		getToEntry, _ := testQueries.GetEntry(context.Background(), toEntry.ID)

		require.NotEmpty(t, getToEntry)

		//check FromEntry data
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account2.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, fromEntry.ID)


		getFromEntry, _ := testQueries.GetEntry(context.Background(), fromEntry.ID)

		require.NotEmpty(t, getFromEntry)

	}
}
