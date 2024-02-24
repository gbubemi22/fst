package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	First_name *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password   *string            `json:"password" validate:"required,min=6"`
	Email      *string            `json:"email" validate:"email,required"`
	Image      *string            `json:"image"`
	User_type  *string            `json:"user_type" validate:"omitempty,eq=ADMIN|eq=USER"`
	Phone      *string            `json:"phone" validate:"required"`

	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id    string    `json:"user_id"`
}
