package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gofiber/fiber/v2"
)

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func TestRegisterAkun(t *testing.T) {
	app := fiber.New()
	mockCollection := new(MockCollection)

	t.Run("Success Registration", func(t *testing.T) {
		mockCollection.On("CountDocuments", mock.Anything, bson.M{"email": "newuser@example.com"}).Return(int64(0), nil)
		mockCollection.On("InsertOne", mock.Anything, mock.Anything).Return(&mongo.InsertOneResult{
			InsertedID: primitive.NewObjectID(),
		}, nil)

		body := map[string]string{
			"email":    "newuser@example.com",
			"password": "password123",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})
}
