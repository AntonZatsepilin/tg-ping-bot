package repository

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoGenerator struct {
    db *mongo.Database
}

func NewMongoGenerator(db *mongo.Database) *MongoGenerator {
    return &MongoGenerator{db: db}
}


func (r *MongoGenerator) GenerateUrls(count int) {
    collection := r.db.Collection("urls")
    for i := 0; i < count; i++ {
        url := gofakeit.URL()
        _, err := collection.InsertOne(context.TODO(), bson.M{"url": url})
        if err != nil {
            log.Printf("failed to insert author: %v", err)
        }
    }
    logrus.Infof("Generated %d URLs", count)	
}

func (r *MongoGenerator) GetUrls() ([]string, error) {
    collection := r.db.Collection("urls")
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    var urls []string
    for cursor.Next(context.TODO()) {
        var doc struct {
            URL string `bson:"url"`
        }
        if err := cursor.Decode(&doc); err != nil {
            log.Printf("Failed to decode document: %v", err)
            continue
        }
        urls = append(urls, doc.URL)
    }
    return urls, nil
}