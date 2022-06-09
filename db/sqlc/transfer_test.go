package db

import (
	"context"
	"testing"
	"time"

	"github.com/Neovastor/go-simple-bank/db/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) (Transfer, error) {
	arg := CreateTransferParams{
		ToAccountID: util.RandomInteger(1,10),
		FromAccountID: util.RandomInteger(1,10),
		Amount: util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotEmpty(t, transfer.ID)
	require.NotEmpty(t, transfer.CreatedAt)

	return transfer, err
}

func TestCreateTransfer(t *testing.T) {

	createRandomEntry(t)
}

func TestGetTransfer(t *testing.T) {
	createdTransfer, _ := createRandomTransfer(t)

	transfer, err := testQueries.GetTransfer(context.Background(), createdTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, createdTransfer.ID, transfer.ID)
	require.Equal(t, createdTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, createdTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, createdTransfer.Amount, transfer.Amount)

	require.WithinDuration(t, createdTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 3,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)		
	}

}
