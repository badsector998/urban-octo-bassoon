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

func createAccoutEntry() Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	if err != nil {
		fmt.Errorf("Create Account error : %v", err)
	}

	return account
}

func createNewEntry(t *testing.T, accID Account) Entry {

	//create Entry
	args := NewEntryParams{
		AccountID: accID.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.NewEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, entry.Amount, args.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestNewEntry(t *testing.T) {
	accEntry := createAccoutEntry()
	createNewEntry(t, accEntry)
}

func TestGetEntry(t *testing.T) {
	entry := createNewEntry(t, createAccoutEntry())

	entryExcute, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryExcute)
	require.Equal(t, entryExcute.ID, entry.ID)
	require.WithinDuration(t, entry.CreatedAt, entryExcute.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {

	accEntry := createAccoutEntry()
	for i := 0; i < 10; i++ {
		createNewEntry(t, accEntry)
	}

	argsListEntry := ListEntriesParams{
		AccountID: accEntry.ID,
		Limit:     5,
		Offset:    5,
	}

	entryExcute, err := testQueries.ListEntries(context.Background(), argsListEntry)
	require.NoError(t, err)
	require.Len(t, entryExcute, 5)

	for _, each_entry := range entryExcute {
		require.NotEmpty(t, each_entry)
	}
}

func TestDeleteEntry(t *testing.T) {
	newEntry := createNewEntry(t, createAccoutEntry())

	err := testQueries.DeleteEntry(context.Background(), newEntry.ID)
	require.NoError(t, err)

	delEntry, err := testQueries.GetEntry(context.Background(), newEntry.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, delEntry)
}
