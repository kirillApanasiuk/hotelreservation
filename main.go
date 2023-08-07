package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/api"
	"hotelreservation/db"
	"log"
	"net/http"
	"os"
)

// Configuration
// 1. Mongodb endpoint
// 2. ListenAddress of our HTTP server
// 3. JWT  secret
// 4. MongoDBName

func main() {
	mongoDbUri := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		log.Fatal(err)
	}
	///handlers initialization
	hotelStore := db.NewMongoHotelStore(client)
	bookingStore := db.NewMongoBookingStore(client)
	userStore := db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(userStore)
	authHandler := api.NewAuthHandler(userStore)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}
	hotelHandler := api.NewHotelHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	roomHandler := api.NewRoomHandler(store)
	config := fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if apiError, ok := err.(api.Error); ok {
				return ctx.Status(apiError.Code).JSON(apiError)
			}
			return api.NewError(http.StatusInternalServerError, err.Error())
		},
	}
	app := fiber.New(config)

	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", api.AdminAuth)

	/// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	/// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// bookings handlers
	apiv1.Get("booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Put("booking/:id/cancel", bookingHandler.HandleCancelBooking)
	admin.Get("bookings", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	_ = app.Listen(listenAddr)

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
