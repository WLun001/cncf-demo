package main

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Item struct {
	ObjectID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Quantity int                `json:"quantity" bson:"quantity"`
}

type Handler struct {
	mongoClient *mongo.Client
	db          *mongo.Database
}

func (h *Handler) HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) GetItems(c echo.Context) error {
	ctx := c.Request().Context()
	cur, err := h.db.Collection(itemCol).Find(ctx, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var items []*Item

	for cur.Next(ctx) {
		var item Item
		if itemErr := cur.Decode(&item); itemErr != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, itemErr)
		}
		items = append(items, &item)
	}
	if cErr := cur.Close(ctx); cErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, cErr)
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) AddItem(c echo.Context) error {
	i := new(Item)
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	_, err := h.db.Collection(itemCol).InsertOne(ctx, i)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
