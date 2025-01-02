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

func UpdateUser(ctx context.Context, db *mongo.Database, userID primitive.ObjectID, updateData bson.M) error {
    collection := db.Collection("users")

    // Update the user in the database
    result, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": updateData})
    if err != nil {
        return fmt.Errorf("failed to update user: %v", err)
    }

    if result.MatchedCount == 0 {
        return fmt.Errorf("user with ID %s not found", userID.Hex())
    }

    return nil
}

// DeleteUser deletes a user from the database
func DeleteUser(ctx context.Context, db *mongo.Database, userID primitive.ObjectID) error {
	collection := db.Collection("users")

	// Delete the user from the database
	_, err := collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
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

func LogOut(ctx context.Context, db *mongo.Database, token string) error {
    // Validasi apakah token sudah diblacklist sebelumnya
    count, err := db.Collection("blacklist_tokens").CountDocuments(ctx, bson.M{"token": token})
    if err != nil {
        return fmt.Errorf("failed to check blacklist: %v", err)
    }
    if count > 0 {
        return fmt.Errorf("token already blacklisted")
    }

    // Tambahkan token ke blacklist dengan waktu kadaluwarsa
    _, err = db.Collection("blacklist_tokens").InsertOne(ctx, bson.M{
        "token":     token,
        "expiresAt": primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour)),
    })
    if err != nil {
        return fmt.Errorf("failed to insert token into blacklist: %v", err)
    }

    return nil
}