package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId     primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	RoomId     primitive.ObjectID `bson:"roomId,omitempty" json:"roomId,omitempty"`
	FromDate   time.Time          `bson:"fromDate" json:"fromDate"`
	TillDate   time.Time          `bson:"tillDate" json:"tillDate"`
	NumPersons int                `bson:"numPersons" json:"numPersons"`
	Canceled   bool               `bson:"canceled,omitempty" json:"canceled,omitempty"`
}
