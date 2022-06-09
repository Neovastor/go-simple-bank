package db

import (
	"context"
	"testing"
	"time"

	"github.com/Neovastor/go-simple-bank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) (Entry, error) {
	arg := CreateEntryParams{
		AccountID: util.RandomInteger(7, 10),
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	return entry, err
}

func TestCreateEntry(t *testing.T) {

	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	createdEntry, _ := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), createdEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, createdEntry.ID, entry2.ID)
	require.Equal(t, createdEntry.AccountID, entry2.AccountID)
	require.Equal(t, createdEntry.Amount, entry2.Amount)
	require.WithinDuration(t, createdEntry.CreatedAt, entry2.CreatedAt, time.Second)

}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 3,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)		
	}

}
