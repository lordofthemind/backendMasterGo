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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/lordofthemind/backendMasterGo/db/mock"
	db "github.com/lordofthemind/backendMasterGo/db/sqlc"
	"github.com/lordofthemind/backendMasterGo/token"
	"github.com/lordofthemind/backendMasterGo/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	// Generate a random account for testing
	account := randomAccount(user.Username)

	// Define test cases
	testCases := []struct {
		name          string // Name of the test case
		body          gin.H
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker) // ID of the account to fetch
		buildStubs    func(store *mockdb.MockStore)                                 // Function to build stubs for the mock store
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)       // Function to check the HTTP response
	}{
		{
			name: "Ok",
			body: gin.H{
				"owner":    account.Owner,
				"balance":  0,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the CreateAccount method to return the generated account
				store.EXPECT().CreateAccount(gomock.Any(), db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
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
			server := newTestServer(t, store)

			// Create a new HTTP recorder for recording the response
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			// Create a new HTTP request to create the account
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// If authentication setup is provided, call it to add authentication to the request
			if tc.setupAuth != nil {
				tc.setupAuth(t, request, server.tokenMaker)
			}

			// Serve the HTTP request and record the response
			server.router.ServeHTTP(recorder, request)

			// Check the HTTP response against the expected response
			tc.checkResponse(t, recorder)
		})
	}

}

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	// Generate a random account for testing
	account := randomAccount(user.Username)

	// Define test cases
	testCases := []struct {
		name          string // Name of the test case
		accountID     int64
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker) // ID of the account to fetch
		buildStubs    func(store *mockdb.MockStore)                                 // Function to build stubs for the mock store
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)       // Function to check the HTTP response
	}{
		{
			name:      "Ok",       // Test case name
			accountID: account.ID, // ID of the account to fetch
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			name:      "UnauthorisedUser", // Test case name
			accountID: account.ID,         // ID of the account to fetch
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, "unauthorised_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the GetAccount method to return the generated account
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization", // Test case name
			accountID: account.ID,        // ID of the account to fetch
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the GetAccount method to return the generated account
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NotFound", // Test case name
			accountID: account.ID, // ID of the account to fetch
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
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
			server := newTestServer(t, store)

			// Create a new HTTP recorder for recording the response
			recorder := httptest.NewRecorder()

			// Create a new HTTP request to fetch the account
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Set up the authentication for the request
			tc.setupAuth(t, request, server.tokenMaker)

			// Serve the HTTP request and record the response
			server.router.ServeHTTP(recorder, request)

			// Check the HTTP response against the expected response
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	rg := utils.NewRandomGenerator()
	return db.Account{
		ID:       rg.RandomInt(1, 1000),
		Owner:    owner,
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

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
