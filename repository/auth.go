package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"learnonbe/model" // Adjust the import path to your project structure

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user in the database
func CreateUser(ctx context.Context, db *mongo.Database, user *model.Users) error {
    // Check if the email already exists
    collection := db.Collection("users")
    count, err := collection.CountDocuments(ctx, bson.M{"email": user.Email})
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("email already exists")
    }

    // Hash the password before saving
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    // Set default role ID if not provided
    if user.RoleID.IsZero() {
        user.RoleID = primitive.NewObjectID()
    }

    // Set user ID and created at timestamp
    user.UserID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    // Insert the user into the database
    _, err = collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }

    return nil
}

func GetUserByID(ctx context.Context, db *mongo.Database, userID primitive.ObjectID) (*model.Users, error) {
    var user model.Users
    err := db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("user with ID %s not found", userID.Hex())
        }
        return nil, err
    }
    return &user, nil
}

// GetUserByEmail retrieves a user by email
const UserCollection = "users"
func GetUserByEmail(ctx context.Context, db *mongo.Database, email string) (*model.Users, error) {
    var user model.Users
    err := db.Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("user with email %s not found", email)
        }
        return nil, err
    }
    return &user, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(ctx context.Context, db *mongo.Database) ([]model.Users, error) {
	var users []model.Users
	cursor, err := db.Collection(UserCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.Users
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}