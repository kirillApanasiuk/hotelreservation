package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/api"
	"hotelreservation/db"
	"log"
)

const mongodburi = "mongodb://localhost:27017"

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodburi))
	if err != nil {
		log.Fatal(err)
	}
	///handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	config := fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(map[string]string{"error": err.Error()})
		},
	}
	app := fiber.New(config)

	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	app.Listen(*listenAddr)
}
