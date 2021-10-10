package internal

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mdb struct {
	coll *mongo.Collection
}

func (db *mdb) Create(data Entity) error {
	_, err := db.coll.InsertOne(context.Background(), data)
	return err
}

func (db *mdb) Get() ([]Entity, error) {
	docs, err := db.coll.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	res := make([]Entity, 0)
	for docs.Next(context.Background()) {
		e := Entity{}
		if err := docs.Decode(e); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (db *mdb) GetOne(id string) (*Entity, error) {
	doc := db.coll.FindOne(context.Background(), bson.D{{"_id", id}})
	e := Entity{}
	if err := doc.Decode(e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (db *mdb) Update(id string, data Entity) error {
	upd := bson.D{
		{"$set", bson.D{{"name", data.Name}}},
		{"$set", bson.D{{"description", data.Description}}},
		{"$set", bson.D{{"enabled", data.Enabled}}},
	}
	_, err := db.coll.UpdateOne(context.Background(), bson.D{{"_id", id}}, upd)
	return err
}

func (db *mdb) Delete(id string) error {
	_, err := db.coll.DeleteOne(context.Background(), bson.D{{"_id", id}})
	return err
}

func NewMongodb(uri string) (DBAdapter, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	if err := client.Connect(context.Background()); err != nil {
		return nil, err
	}
	db := client.Database("entities")
	coll := db.Collection("entities")
	return &mdb{coll}, nil
}
