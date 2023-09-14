package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/db/fixtures"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "james", "foo", false)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "james@foo.com",
		Password: "james_foo",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		var body genericResp
		err := json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("expected status code 200, got %d, msg: %s", resp.StatusCode, body.Msg)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response")
	}

	// Set the encrypted password to empty string to be able to compare the user
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResp.User)
		t.Fatal("expected the user to be present in the auth response")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(tdb.Store, "James1", "Foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "eqwfwe@sverw.com",
		Password: "password1",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status code 401, got %d", resp.StatusCode)
	}
	var authResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Type != "error" {
		t.Fatalf("expected the response type to be error, but got %s", authResp.Type)
	}
	if authResp.Msg != "invalid credentials" {
		t.Fatalf("expected the response msg to be 'invalid credentials', but got %s", authResp.Msg)
	}
}
