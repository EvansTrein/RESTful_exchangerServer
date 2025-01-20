package auth_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

const host = "http://localhost:8000"

var (
	testDataName     = "testName"
	testDataPassword = "123456"
	testDataEmail    = "test@mail.com"
	token            = ""
)

func TestRegisterHandler(t *testing.T) {
	urlPath := "/api/v1/register"

	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
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

	t.Run("Invalid JSON body request (no name)", func(t *testing.T) {
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

	t.Run("successful registration", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
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

	t.Run("Invalid email already exists", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"username": testDataName,
			"password": testDataPassword,
			"email":    testDataEmail,
		}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").ValueEqual("error", "email already exists")
		testCase.ContainsKey("message").ValueEqual("message", "failed to save a new user")
	})
}

func TestLoginHandler(t *testing.T) {
	urlPath := "/api/v1/login"
	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("Invalid email not found", func(t *testing.T) {

		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"email":    "failEmail@mail.com",
			"password": testDataPassword,
		}).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "user not found")
	})

	t.Run("Invalid password", func(t *testing.T) {
		urlPath := "/api/v1/login"

		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"email":    testDataEmail,
			"password": "invalidPass",
		}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "invalid email or password")
	})

	t.Run("Invalid JSON body request (no email)", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
			"password": testDataPassword,
		}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "invalid data")
	})

	t.Run("successful login", func(t *testing.T) {
		testCase := testHTTP.POST(urlPath).WithJSON(map[string]string{
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

func TestDeleteHandlerAndLoggingMiddleware(t *testing.T) {
	urlPath := "/api/v1/delete"
	testHTTP := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  host,
		Reporter: httpexpect.NewRequireReporter(t),
		Client:   http.DefaultClient,
	})

	t.Run("Invalid not header Authorization", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("Invalid header Authorization", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).WithHeader("Authori", "Bearer "+token).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("Invalid not token", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).WithHeader("Authorization", "Bearer "+"").
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("Invalid token format", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).WithHeader("Authorization", "Bearer "+"jher34234koq").
			Expect().
			Status(http.StatusInternalServerError).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "failed to check the token")
	})

	t.Run("Invalid Authorization header format", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).WithHeader("Authorization", "Lalala"+token).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "unauthorized user")
	})

	t.Run("successful delete", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("message").ValueEqual("message", "user successfully deleted")
	})

	t.Run("Invalid repeated removal", func(t *testing.T) {
		testCase := testHTTP.DELETE(urlPath).WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound).
			JSON().Object().NotEmpty()

		testCase.ContainsKey("error").Value("error").String().NotEmpty()
		testCase.ContainsKey("message").ValueEqual("message", "user does not exist")
	})
}
