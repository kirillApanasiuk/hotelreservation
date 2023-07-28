package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hotelreservation/db"
	"hotelreservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "Kirill",
		LastName:  "Apanasiuk",
		Email:     "milidanex@gmail.com",
		Password:  "admin",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	testDb := setup(t)

	defer testDb.teardown(t)
	insertTestUser(t, testDb.UserStore)
	app := fiber.New()
	authhandler := NewAuthHandler(testDb.UserStore)
	app.Post("/auth", authhandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "milidanex@gmail.com",
		Password: "admin",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but god %d", resp.StatusCode)
	}

	var genREsp genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genREsp); err != nil {
		t.Fatal(err)
	}
	if genREsp.Type != "error" {
		t.Fatalf("expected gen response type to be error %s", genREsp.Type)
	}
	if genREsp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response to be invalid credential %s", genREsp.Msg)
	}
}
