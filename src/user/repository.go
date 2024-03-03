package user

import (
	"context"
	"errors"
	"fmt"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, collectionName string) *UserRepository {
	collection := client.Database("FASTER").Collection("users")
	 
	return &UserRepository{
		Collection: collection,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *User) error {

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	_, err := ur.Collection.InsertOne(ctx, user)
	return err
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := ur.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, userID string, updatedUser *User) error {
	updatedUser.Updated_at = time.Now()
	_, err := ur.Collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": updatedUser})
	return err
}

func (ur *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	_, err := ur.Collection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := ur.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		// Check if the error is due to no documents found
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user with email '%s' not found", email)
		}
		return nil, err
	}
	fmt.Println("CHECK@@@@@@@",user)

	return &user, nil
}

func (ur *UserRepository) ListAllUsers(ctx context.Context) ([]*User, error) {
	cursor, err := ur.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (ur *UserRepository) ListOne(ctx context.Context, userID string) (*User, error) {
	var user User
	err := ur.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
