package types

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"fistName"`
	LastName  string `bson:"lastName" json:"lastName"`
}