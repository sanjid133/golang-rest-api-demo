package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type Client struct {
	client *mongo.Client
	db string
}

func NewClient(uri, db string)(*Client, error)  {
	opts := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)
	c, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	return &Client{client:c, db: db}, nil
}

type Row struct {
	s   *mongo.SingleResult
	mu  sync.RWMutex
	raw bson.Raw
	err error
}



