package wallet_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/EvansTrein/RESTful_exchangerServer/models"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

const host = "http://localhost:8000"

var (
	testDataName     = "walletTest"
	testDataPassword = "123456"
	testDataEmail    = "walletTest@mail.com"
	token            = ""
)

func TestRegisterAndLogin(t *testing.T) {
	urlPathReg := "/api/v1/register"
	urlPathLog := "/api/v1/login"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("successful registration for tests", func(t *testing.T) {
		testCase := testHTTP.POST(urlPathReg).WithJSON(map[string]string{
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
		testCase := testHTTP.POST(urlPathLog).WithJSON(map[string]string{
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
	urlPathBalance := "/api/v1/balance"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("successful balance zero", func(t *testing.T) {
		testCase := testHTTP.GET(urlPathBalance).WithHeader("Authorization", "Bearer "+token).
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
	urlPathDeposit := "/api/v1/wallet/deposit"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("deposit fail currency not found", func(t *testing.T) {
		amount := 10000
		currency := "XXXX"

		testCase := testHTTP.POST(urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
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

		testCase := testHTTP.POST(urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
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
		amount := 1000
		currency := "USD"

		testCase := testHTTP.POST(urlPathDeposit).WithHeader("Authorization", "Bearer "+token).
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

		assert.EqualValues(t, 1000, v, "The value of the field is not equal, which is not expected")
	})
}

func TestBalanceDepositWithdraw(t *testing.T) {
	t.Skip()
	// urlPathBalance := "/api/v1/balance"
	// urlPathDeposit := "/api/v1/wallet/deposit"
	// urlPathWithdraw := "/api/v1/wallet/withdraw"

	// testHTTP := httpexpect.WithConfig(httpexpect.Config{
	// 	BaseURL:  host,
	// 	Reporter: httpexpect.NewRequireReporter(t),
	// 	Client:   http.DefaultClient,
	// })

	// t.Run("", func(t *testing.T) {
	// 	testCase := testHTTP.GET(urlPathBalance).WithHeader("Authorization", "Bearer "+token).
	// 		Expect().
	// 		Status(http.StatusOK).
	// 		JSON().Object().NotEmpty()

		
	// })

}

func TestDeleteUserAfter(t *testing.T) {
	urlPathDel := "/api/v1/delete"
	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("successful delete for tests", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPathDel).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("message").ValueEqual("message", "user successfully deleted")
	})
}
