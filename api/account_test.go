package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/lordofthemind/backendMasterGo/db/mock"
	db "github.com/lordofthemind/backendMasterGo/db/sqlc"
	"github.com/lordofthemind/backendMasterGo/utils"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	// Generate a random account for testing
	account := randomAccount()

	// Define test cases
	testCases := []struct {
		name          string                                                  // Name of the test case
		accountID     int64                                                   // ID of the account to fetch
		buildStubs    func(store *mockdb.MockStore)                           // Function to build stubs for the mock store
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder) // Function to check the HTTP response
	}{
		{
			name:      "Ok",       // Test case name
			accountID: account.ID, // ID of the account to fetch
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the GetAccount method to return the generated account
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound", // Test case name
			accountID: account.ID, // ID of the account to fetch
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the GetAccount method to return the generated account
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError", // Test case name
			accountID: account.ID,      // ID of the account to fetch
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the GetAccount method to return the generated account
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidId", // Test case name
			accountID: 0,           // ID of the account to fetch
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the GetAccount method to return the generated account
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// Iterate through test cases
	for i := range testCases {
		// Get the current test case
		tc := testCases[i]

		// Run the test case
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock store
			store := mockdb.NewMockStore(ctrl)

			// Build stubs for the mock store
			tc.buildStubs(store)

			// Create a new server instance with the mock store
			server := NewServer(store)

			// Create a new HTTP recorder for recording the response
			recorder := httptest.NewRecorder()

			// Create a new HTTP request to fetch the account
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Serve the HTTP request and record the response
			server.router.ServeHTTP(recorder, request)

			// Check the HTTP response against the expected response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	rg := utils.NewRandomGenerator()
	return db.Account{
		ID:       rg.RandomInt(1, 1000),
		Owner:    rg.RandomOwner(),
		Balance:  rg.RandomMoney(),
		Currency: rg.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)

}
