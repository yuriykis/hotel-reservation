package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/api"
	"github.com/yuriykis/hotel-reservation/api/middleware"
	"github.com/yuriykis/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		userHandler  = api.NewUserHandler(store)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(store.User)
		roomHandler  = api.NewRoomHandler(store)
		bookinghand  = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", middleware.JWTAutentication(store.User))
	)
	// Auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// vaerioned api routes
	// users handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// hotels handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// rooms handlers
	apiv1.Get("/room", roomHandler.HanldeGetRooms)
	apiv1.Post("room/:id/book", roomHandler.HandleBookRoom)

	// TODO: cancel booking

	// bookings handlers
	apiv1.Get("/booking", bookinghand.HandleGetBookings)
	apiv1.Get("/booking/:id", bookinghand.HandleGetBooking)

	app.Listen(*listenAddr)
}
