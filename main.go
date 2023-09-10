package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/api"
	"github.com/yuriykis/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userCollection = "users"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{
			"error": err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5001", "server listen address")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	apiv1.Post("/user", userHandler.HandlePostUser)

	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	app.Listen(*listenAddr)
}
