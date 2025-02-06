package repository

import "go.mongodb.org/mongo-driver/mongo"

type Generate interface {
    GenerateUrls(count int)
    GetUrls() ([]string, error)
}

type Repository struct {
    Generator Generate
}

func NewRepository(db *mongo.Database) *Repository {
    return &Repository{
        Generator: NewMongoGenerator(db),
    }
}