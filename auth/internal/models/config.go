package models

type Config struct {
    Database struct {
        Host     string
        Port     int
        User     string
        Password string
        DBName   string
    }
    JWT struct {
        Secret string
    }
}