package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/yuriykis/hotel-reservation/api"
	"github.com/yuriykis/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config
// 1. Mongo db
// 2. Listen address for HTTP server
// JWT secret
// MongoDBName

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongodburi := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodburi))
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
		apiv1 = app.Group("/api/v1", api.JWTAutentication(store.User))
		admin = apiv1.Group("/admin", api.AdminAuth)
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

	// bookings handlers
	apiv1.Get("/booking/:id", bookinghand.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookinghand.HandleCancelBooking)

	// admin routes
	admin.Get("/booking", bookinghand.HandleGetBookings)

	httpListenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(httpListenAddr)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
