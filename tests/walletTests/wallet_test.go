package wallet_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

const (
	host       = "http://localhost:8000"
	apiVersion = "/api/v1"
)

var (
	testDataName     = "walletTest"
	testDataPassword = "123456"
	testDataEmail    = "walletTest@mail.com"
	token            = ""
)

func TestRegisterAndLogin(t *testing.T) {
	urlPathReg := "/register"
	urlPathLog := "/login"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("successful registration for tests", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion + urlPathReg).WithJSON(map[string]string{
			"username": testDataName,
			"password": testDataPassword,
			"email":    testDataEmail,
		}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("id").Value("id").Number()
		testCase.ContainsKey("message").ValueEqual("message", "user successfully created")
	})

	t.Run("successful login for tests", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion + urlPathLog).WithJSON(map[string]string{
			"email":    testDataEmail,
			"password": testDataPassword,
		}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("token").Value("token").String().NotEmpty()
		token = testCase.Value("token").String().Raw()
	})
}

func TestBalance(t *testing.T) {
	urlPathBalance := "/balance"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("balance fail not header Authorization", func(t *testing.T) {
		testCase := testHTTP.GET(apiVersion + urlPathBalance).WithJSON(map[string]interface{}{
			"amount":   1200,
			"currency": "EUR",
		}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("successful balance zero", func(t *testing.T) {
		testCase := testHTTP.GET(apiVersion+urlPathBalance).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("balance")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var balanceResponse models.BalanceResponse
		err = json.Unmarshal(jsonData, &balanceResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		for _, v := range balanceResponse.Balance {
			assert.EqualValues(t, 0.0, v, "Field value is not zero, which is not expected")
		}
	})
}

func TestDeposit(t *testing.T) {
	urlPathDeposit := "/wallet/deposit"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("deposit fail not header Authorization", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion + urlPathDeposit).WithJSON(map[string]interface{}{
			"amount":   1200,
			"currency": "EUR",
		}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("deposit fail currency not found", func(t *testing.T) {
		amount := 10000
		currency := "XXXX"

		testCase := testHTTP.POST(apiVersion+urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "currency not found")
	})

	t.Run("deposit fail user account not found", func(t *testing.T) {
		// when creating a user, accounts in currencies are created automatically, the test is irrelevant
		t.Skip("Skipping this test as it is currently inactive")
		amount := 1000
		currency := "GBP"

		testCase := testHTTP.POST(apiVersion+urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "account not found")
	})

	t.Run("deposit successful", func(t *testing.T) {
		amount := 2000
		currency := "USD"

		testCase := testHTTP.POST(apiVersion+urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("new_balance")
		testCase.ContainsKey("message").ValueEqual("message", "successfully deposit")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var depositResponse models.AccountOperationResponse
		err = json.Unmarshal(jsonData, &depositResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		v, ok := depositResponse.NewBalance[currency]
		if !ok {
			t.Errorf("recharged account - %s is not in the response", currency)
		}

		assert.EqualValues(t, amount, v, "The value of the field is not equal, which is not expected")
	})
}

func TestWithdraw(t *testing.T) {
	urlPathWithdraw := "/wallet/withdraw"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("withdraw fail not header Authorization", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion + urlPathWithdraw).WithJSON(map[string]interface{}{
			"amount":   1200,
			"currency": "EUR",
		}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("withdraw fail insufficient funds", func(t *testing.T) {
		amount := 2500
		currency := "USD"

		testCase := testHTTP.POST(apiVersion+urlPathWithdraw).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusPaymentRequired).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "insufficient funds")
	})

	t.Run("withdraw fail currency not found", func(t *testing.T) {
		amount := 5000
		currency := "XXXX"

		testCase := testHTTP.POST(apiVersion+urlPathWithdraw).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "currency is not supported")
	})

	t.Run("withdraw successful", func(t *testing.T) {
		amount := 1000
		currency := "USD"

		testCase := testHTTP.POST(apiVersion+urlPathWithdraw).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("new_balance")
		testCase.ContainsKey("message").ValueEqual("message", "successfully withdrawn")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var depositResponse models.AccountOperationResponse
		err = json.Unmarshal(jsonData, &depositResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		v, ok := depositResponse.NewBalance[currency]
		if !ok {
			t.Errorf("debit currency account - %s is missing in the response", currency)
		}

		// 2000 (from the test above) - 1000 = 1000
		assert.EqualValues(t, amount, v, "The value of the field is not equal, which is not expected")
	})
}

func TestBalanceDepositWithdraw(t *testing.T) {
	urlPathBalance := "/balance"
	urlPathDeposit := "/wallet/deposit"
	urlPathWithdraw := "/wallet/withdraw"

	amount := 3000
	currency := "EUR"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("deposit EUR", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion+urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("new_balance")
		testCase.ContainsKey("message").ValueEqual("message", "successfully deposit")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var depositResponse models.AccountOperationResponse
		err = json.Unmarshal(jsonData, &depositResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		v, ok := depositResponse.NewBalance[currency]
		if !ok {
			t.Errorf("recharged account - %s is not in the response", currency)
		}

		assert.EqualValues(t, amount, v, "The value of the field is not equal, which is not expected")
	})

	t.Run("check balance and account EUR", func(t *testing.T) {
		testCase := testHTTP.GET(apiVersion+urlPathBalance).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("balance")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var balanceResponse models.BalanceResponse
		err = json.Unmarshal(jsonData, &balanceResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		v, ok := balanceResponse.Balance[currency]
		if !ok {
			t.Errorf("debit currency account - %s is missing in the response", currency)
		}

		if v != float32(amount) {
			t.Errorf("incorrect account %s balance after top-up %v", currency, amount)
		}
	})

	t.Run("withdraw EUR", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion+urlPathWithdraw).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"amount":   amount,
				"currency": currency,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("new_balance")
		testCase.ContainsKey("message").ValueEqual("message", "successfully withdrawn")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var depositResponse models.AccountOperationResponse
		err = json.Unmarshal(jsonData, &depositResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		v, ok := depositResponse.NewBalance[currency]
		if !ok {
			t.Errorf("debit currency account - %s is missing in the response", currency)
		}

		// 3000 - 3000 = 0
		assert.EqualValues(t, 0, v, "The value of the field is not equal, which is not expected")
	})

	t.Run("final check balance", func(t *testing.T) {
		testCase := testHTTP.GET(apiVersion+urlPathBalance).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("balance")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var balanceResponse models.BalanceResponse
		err = json.Unmarshal(jsonData, &balanceResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		v, ok := balanceResponse.Balance[currency]
		if !ok {
			t.Errorf("debit currency account - %s is missing in the response", currency)
		}

		if v != 0 {
			t.Errorf("incorrect account %s balance after withdraw %v", currency, amount)
		}

		for k, v := range balanceResponse.Balance {
			if v < 0 {
				t.Errorf("negative balance %s - %v", k, v)
			}
		}
	})
}

func TestAllRates(t *testing.T) {
	urlPathAllRates := "/exchange/rates"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("exchange all rates fail not header Authorization", func(t *testing.T) {
		testCase := testHTTP.GET(apiVersion + urlPathAllRates).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("successful exchange all rates", func(t *testing.T) {
		testCase := testHTTP.GET(apiVersion+urlPathAllRates).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("rates")
		testCase.ContainsKey("message").ValueEqual("message", "data successfully received")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var balanceResponse models.ExchangeRatesResponse
		err = json.Unmarshal(jsonData, &balanceResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		for k, v := range balanceResponse.Rates {
			if v < 0 {
				t.Errorf("negative rate %s - %v", k, v)
			}
		}
	})
}

func TestExchange(t *testing.T) {
	urlPathExchange := "/exchange"
	fromCurrency := "USD"
	toCurrency := "CNY"
	amount := 500

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("exchange fail not header Authorization", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion + urlPathExchange).WithJSON(map[string]interface{}{
			"from_currency": "RUB",
			"to_currency":   "EUR",
			"amount":        5000,
		}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("Invalid JSON body request (no to_currency)", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion+urlPathExchange).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"from_currency": "RUB",
				"amount":        5000,
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "invalid data")
	})

	t.Run("successful exchange", func(t *testing.T) {
		testCase := testHTTP.POST(apiVersion+urlPathExchange).WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"from_currency": fromCurrency,
				"to_currency":   toCurrency,
				"amount":        amount,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("exchange_rate")
		testCase.ContainsKey("spent_accoutn")
		testCase.ContainsKey("received_account")
		testCase.ContainsKey("new_balance")
		testCase.ContainsKey("message").ValueEqual("message", "currency exchange successfully")

		jsonData, err := json.Marshal(testCase.Raw())
		if err != nil {
			t.Errorf("Failed to marshal raw data to JSON: %v", err)
		}

		var exchangeResponse models.ExchangeResponse
		err = json.Unmarshal(jsonData, &exchangeResponse)
		if err != nil {
			t.Errorf("Failed to decode JSON response: %v", err)
		}

		assert.Greater(t, exchangeResponse.ExchangeRate, float32(0), "ExchangeRate must be greater than zero")
		assert.Equal(t, fromCurrency, exchangeResponse.SpentAccoutn.Currency, "SpentAccoutn.Currency must be equal to USD")
		// 1000 USD (from the test above) - 500 = 500
		assert.Equal(t, float32(500), exchangeResponse.SpentAccoutn.Amount, "SpentAccoutn.Amount must equal 500")
		assert.Equal(t, toCurrency, exchangeResponse.ReceivedAccount.Currency, "ReceivedAccount.Currency must be equal to CNY")
		assert.Greater(t, exchangeResponse.ReceivedAccount.Amount, float32(0), "ReceivedAccount.Amount must be greater than zero")
	})
}

func TestDeleteUserAfter(t *testing.T) {
	urlPathDel := "/delete"
	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("successful delete for tests", func(t *testing.T) {
		testCase := testHTTP.DELETE(apiVersion+urlPathDel).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("message").ValueEqual("message", "user successfully deleted")
	})
}
