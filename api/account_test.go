package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/lordofthemind/backendMasterGo/db/mock"
	db "github.com/lordofthemind/backendMasterGo/db/sqlc"
	"github.com/lordofthemind/backendMasterGo/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {
	// Generate a random account for testing
	account := randomAccount()

	// Define test cases
	testCases := []struct {
		name          string                                                  // Name of the test case
		body          gin.H                                                   // Request body
		buildStubs    func(store *mockdb.MockStore)                           // Function to build stubs for the mock store
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder) // Function to check the HTTP response
	}{
		{
			name: "Ok",
			body: gin.H{
				"owner":    account.Owner,
				"balance":  0,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Stub the CreateAccount method to return the generated account
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check if the HTTP status code is OK (200) and the response body matches the generated account
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"balance":  0,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, errors.New("internal server error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidData",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations, as no request should be made
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingOwnerField",
			body: gin.H{
				"balance":  0,
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations, as no request should be made
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		// {
		// 	name: "MissingBalanceField",
		// 	body: gin.H{
		// 		"owner":    account.Owner,
		// 		"currency": account.Currency,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		// No expectations, as no request should be made
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },

		{
			name: "MissingCurrencyField",
			body: gin.H{
				"owner":   account.Owner,
				"balance": 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// No expectations, as no request should be made
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			// Create a new HTTP request to create the account
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			// Serve the HTTP request and record the response
			server.router.ServeHTTP(recorder, request)

			// Check the HTTP response against the expected response
			tc.checkResponse(t, recorder)
		})
	}

}

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

// func TestListAccountsAPI(t *testing.T) {
// 	user, _ := randomUser(t)

// 	n := 5
// 	accounts := make([]db.Account, n)
// 	for i := 0; i < n; i++ {
// 		accounts[i] = randomAccount(user.Username)
// 	}

// 	type Query struct {
// 		pageID   int
// 		pageSize int
// 	}

// 	testCases := []struct {
// 		name          string
// 		query         Query
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.ListAccountsParams{
// 					Owner:  user.Username,
// 					Limit:  int32(n),
// 					Offset: 0,
// 				}

// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(accounts, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccounts(t, recorder.Body, accounts)
// 			},
// 		},
// 		{
// 			name: "NoAuthorization",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return([]db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageID",
// 			query: Query{
// 				pageID:   -1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageSize",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: 100000,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			url := "/accounts"
// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			// Add query parameters to request URL
// 			q := request.URL.Query()
// 			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
// 			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
// 			request.URL.RawQuery = q.Encode()

// 			tc.setupAuth(t, request, server.tokenMaker)
// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

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

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
