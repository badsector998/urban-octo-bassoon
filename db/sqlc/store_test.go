package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)
	account1 := accountRandomTest(t)
	Account2 := accountRandomTest(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	result_chan := make(chan TransferTxResult)

	// run n cocurent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   Account2.ID,
				Amount:        amount,
			})

			errs <- err
			result_chan <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-result_chan
		require.NotEmpty(t, result)

		// result store multiple objects.
		// check one by one

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, Account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		// check entry account1

		// check entry account2
	}

}
