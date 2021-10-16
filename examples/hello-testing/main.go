package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	e           *echo.Echo
	dbName      = "store"
	itemCol     = "items"
)

func connectToDB() (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		return nil, nil, err
	}

	mongoClient = client

	log.Println("db client connected")
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, err
	}
	log.Println("db client ping")

	db := client.Database(dbName)
	return client, db, nil
}

func run() (*Handler, *echo.Echo, error) {
	c, d, err := connectToDB()
	if err != nil {
		return nil, nil, err
	}

	handler := Handler{
		mongoClient: c,
		db:          d,
	}

	e = echo.New()
	e.GET("/", handler.HelloWorld)
	e.GET("/items", handler.GetItems)
	e.POST("/items", handler.AddItem)

	return &handler, e, nil
}

func stop() {
	ctx := context.TODO()
	_ = e.Shutdown(ctx)
	_ = mongoClient.Disconnect(ctx)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	_, ee, err := run()
	if err != nil {
		panic(err)
	}
	ee.Logger.Fatal(ee.Start(":1323"))
}
