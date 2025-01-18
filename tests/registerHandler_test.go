package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestRegisterHandler(t *testing.T) {
	const host = "http://localhost:8000"
	const urlPath = "/api/v1/register"
	
	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewAssertReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("Invalid email", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"username": "name",
			"password": "123456",
			"email":    "namemail.com",
		}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

			testCase.ContainsKey("error").Value("error").String().NotEmpty()
			testCase.ContainsKey("message").ValueEqual("message", "invalid data")
	})

	t.Run("Invalid password", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"username": "name",
			"password": "1234", // min 6 symbols
			"email":    "name@mail.com",
		}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

			testCase.ContainsKey("error").Value("error").String().NotEmpty()
			testCase.ContainsKey("message").ValueEqual("message", "invalid data")
	})

	t.Run("invalid JSON body request (no name)", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"password": "123456", 
			"email":    "name@mail.com",
		}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

			testCase.ContainsKey("error").Value("error").String().NotEmpty()
			testCase.ContainsKey("message").ValueEqual("message", "invalid data")
	})
}
