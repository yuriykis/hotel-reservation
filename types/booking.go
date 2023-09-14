package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID     primitive.ObjectID `bson:"roomID" json:"roomID"`
	UserID     primitive.ObjectID `bson:"userID" json:"userID"`
	NumPersons int                `bson:"numPersons" json:"numPersons"`
	FromDate   time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate   time.Time          `bson:"tillDate" json:"tillDate"`
	Canceled   bool               `bson:"canceled" json:"canceled"`
}

func NewBooking(roomID, userID primitive.ObjectID, numPersons int, fromDate, tillDate time.Time) *Booking {
	return &Booking{
		RoomID:     roomID,
		UserID:     userID,
		NumPersons: numPersons,
		FromDate:   fromDate,
		TillDate:   tillDate,
		Canceled:   false,
	}
}
