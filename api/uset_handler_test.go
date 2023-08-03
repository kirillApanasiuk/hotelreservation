package api

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hotelreservation/types"
	"net/http/httptest"
	"testing"
)

func TestPostUers(t *testing.T) {
	testDb := setup(t)
	defer testDb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testDb.Store.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Kirill",
		LastName:  "Apanasiuk",
		Email:     "some@foo.com",
		Password:  "admin",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be includded in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected username %s but got %s", params.FirstName, user.FirstName)

	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)

	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)

	}
}
