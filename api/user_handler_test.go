package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/db"
	"github.com/yuriykis/hotel-reservation/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	store := &db.Store{
		User: tdb.User,
	}
	userHandler := NewUserHandler(store)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "foo@bar.com",
		Password:  "password",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	user := types.User{}
	json.NewDecoder(res.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expected user id, got %s", user.ID)
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected not to be present to json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected %s, got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected %s, got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected %s, got %s", params.Email, user.Email)
	}
}
