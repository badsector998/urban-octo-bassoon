package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/badsector998/urban-octo-bassoon/util"
	"github.com/stretchr/testify/require"
)

func createAccounts() Account {
	//create two accounts as from and to.
	accTrfArg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), accTrfArg)
	if err != nil {
		fmt.Errorf("Errror Creating From Account error code : %v", err)
	}

	return account
}

func createNewTransfer(t *testing.T, fromAcc Account, toAcc Account) Transfer {

	//create transfer
	txArgs := NewTransferParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	}

	newTx, err := testQueries.NewTransfer(context.Background(), txArgs)
	require.NoError(t, err)
	require.NotEmpty(t, newTx)
	require.NotEqual(t, newTx.FromAccountID, newTx.ToAccountID)

	return newTx
}

func TestNewTransfer(t *testing.T) {
	createNewTransfer(t, createAccounts(), createAccounts())
}

func TestGetTransfer(t *testing.T) {
	tx1 := createNewTransfer(t, createAccounts(), createAccounts())

	tx2, err := testQueries.GetTransfer(context.Background(), tx1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tx2)
	require.Equal(t, tx2.ID, tx1.ID)
	require.WithinDuration(t, tx1.CreatedAt, tx2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	fromAcc := createAccounts()
	toAcc := createAccounts()
	for i := 0; i < 10; i++ {
		createNewTransfer(t, fromAcc, toAcc)
	}

	args := ListTransfersParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Limit:         5,
		Offset:        5,
	}

	txList, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, txList, 5)

	for _, listTx := range txList {
		require.NotEmpty(t, listTx)
	}
}

func TestDeleteTransfer(t *testing.T) {
	tx1 := createNewTransfer(t, createAccounts(), createAccounts())

	err := testQueries.DeleteTransfer(context.Background(), tx1.ID)
	require.NoError(t, err)

	tx2, err := testQueries.GetTransfer(context.Background(), tx1.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tx2)
}
