package repository

import (
	"context"

	"learnonbe/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create Role
func CreateRole(ctx context.Context, db *mongo.Database, role *model.Roles) error {
	collection := db.Collection("roles")
	_, err := collection.InsertOne(ctx, role)
	if err != nil {
		return err
	}

	return nil
}

// GetRoleByID
func GetRoleByID(ctx context.Context, db *mongo.Database, roleID primitive.ObjectID) (*model.Roles, error) {
	var role model.Roles
	err := db.Collection("roles").FindOne(ctx, bson.M{"_id": roleID}).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}