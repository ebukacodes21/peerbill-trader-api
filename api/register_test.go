package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "peerbill-trader-server/db/mock"
	db "peerbill-trader-server/db/sqlc"
	"peerbill-trader-server/utils"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateTraderParamsMatcher struct {
	arg      db.CreateTraderParams
	password string
}

func (e eqCreateTraderParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateTraderParams)
	if !ok {
		return false
	}

	err := utils.VerifyPassword(arg.Password, e.password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateTraderParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateTraderParams(arg db.CreateTraderParams, password string) gomock.Matcher {
	return eqCreateTraderParamsMatcher{arg, password}
}

func TestCreateTrader(t *testing.T) {
	trader, password := randomTrader(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockDatabaseContract)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":   trader.Username,
				"first_name": trader.FirstName,
				"last_name":  trader.LastName,
				"email":      trader.Email,
				"password":   password,
				"country":    trader.Country,
				"phone":      trader.Phone,
			},
			buildStubs: func(store *mockdb.MockDatabaseContract) {
				arg := db.CreateTraderParams{
					Username:  trader.Username,
					FirstName: trader.FirstName,
					LastName:  trader.LastName,
					Email:     trader.Email,
					Country:   trader.Country,
					Phone:     trader.Phone,
				}
				store.EXPECT().
					CreateTrader(gomock.Any(), EqCreateTraderParams(arg, password)).
					Times(1).
					Return(trader, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTrader(t, recorder.Body, trader)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":   "invalid-trader#1",
				"first_name": trader.FirstName,
				"last_name":  trader.LastName,
				"email":      trader.Email,
				"password":   password,
				"country":    trader.Country,
				"phone":      trader.Phone,
			},
			buildStubs: func(store *mockdb.MockDatabaseContract) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":   trader.Username,
				"first_name": trader.FirstName,
				"last_name":  trader.LastName,
				"email":      "email",
				"password":   password,
				"country":    trader.Country,
				"phone":      trader.Phone,
			},
			buildStubs: func(store *mockdb.MockDatabaseContract) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":   trader.Username,
				"first_name": trader.FirstName,
				"last_name":  trader.LastName,
				"email":      trader.Email,
				"password":   password,
				"country":    trader.Country,
				"phone":      trader.Phone,
			},
			buildStubs: func(store *mockdb.MockDatabaseContract) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Trader{}, db.ErrUniqueViolation)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":   trader.Username,
				"first_name": trader.FirstName,
				"last_name":  trader.LastName,
				"email":      trader.Email,
				"password":   password,
				"country":    trader.Country,
				"phone":      trader.Phone,
			},
			buildStubs: func(store *mockdb.MockDatabaseContract) {
				store.EXPECT().
					CreateTrader(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Trader{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	// test proper
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repository := mockdb.NewMockDatabaseContract(ctrl)
			tc.buildStubs(repository)

			config := utils.Config{
				HTTPServerAddr: utils.RandomOwner(),
				DBDriver:       utils.RandomOwner(),
				DBSource:       utils.RandomOwner(),
			}

			server, _ := NewServer(config, repository)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/register-trader"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomTrader(t *testing.T) (trader db.Trader, password string) {
	password = utils.RandomString(6)
	Password, err := utils.HashPassword(password)
	require.NoError(t, err)

	trader = db.Trader{
		FirstName: utils.RandomOwner(),
		LastName:  utils.RandomOwner(),
		Username:  utils.RandomOwner(),
		Password:  Password,
		Email:     utils.RandomEmail(),
		Country:   utils.RandomOwner(),
		Phone:     utils.RandomPhone(),
	}
	return
}

func requireBodyMatchTrader(t *testing.T, body *bytes.Buffer, user db.Trader) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTrader db.Trader
	err = json.Unmarshal(data, &gotTrader)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotTrader.Username)
	require.Equal(t, user.FirstName, gotTrader.FirstName)
	require.Equal(t, user.Email, gotTrader.Email)
	require.Empty(t, gotTrader.Password)
}
