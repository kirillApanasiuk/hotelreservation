package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hotelreservation/db/fixtures"
	"hotelreservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "james", "test", true)
		hotel          = fixtures.AddHotel(db.Store, "Test hotel", "belarus", 1, nil)
		room           = fixtures.AddRoom(db.Store, "small", false, 12, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		bookingHandler = NewBookingHandler(db.Store)
		app            = fiber.New()
		admin          = app.Group("/", JWTAuthentication(db.Store.User), AdminAuth)
		bookings       = fixtures.AddBooking(db.Store, adminUser.ID, room.ID, from, till, 15)
	)

	_ = bookings

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 respose %v", resp.StatusCode)
	}
	var booking []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		t.Fatal(err)
	}

	if len(booking) != 1 {
		t.Fatalf("expected 1 booking got %v", len(booking))
	}
	fmt.Println(bookings)
}
