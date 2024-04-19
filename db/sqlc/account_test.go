package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lordofthemind/backendMasterGo/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {

	rg := utils.NewRandomGenerator()
	arg := CreateAccountParams{
		Owner:    rg.RandomOwner(),
		Balance:  rg.RandomMoney(),
		Currency: rg.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {

	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccounts(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.NotZero(t, account2.ID)
	require.NotZero(t, account2.CreatedAt)
}

func TestUpdateAccount(t *testing.T) {
	rg := utils.NewRandomGenerator()
	account1 := CreateRandomAccount(t)
	arg := UpdateAccountsParams{
		ID:      account1.ID,
		Balance: rg.RandomMoney(),
	}
	err := testQueries.UpdateAccounts(context.Background(), arg)
	require.NoError(t, err)

	// Retrieve the updated account to verify the changes
	account2, err := testQueries.GetAccounts(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// Add assertions here to validate the behavior of UpdateAccount method
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccounts(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccounts(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
