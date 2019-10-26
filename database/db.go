package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Client) Insert(ctx context.Context, col string, data interface{}) (interface{}, error)  {
	r, err := c.client.Database(c.db).Collection(col).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return r.InsertedID, nil
}

func (c *Client) FindID(ctx context.Context, col, id string) (*Row, error)   {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filtr := bson.M{"_id": objID}
	res := c.client.Database(c.db).Collection(col).FindOne(ctx, filtr)
	return &Row{s: res }, nil
}