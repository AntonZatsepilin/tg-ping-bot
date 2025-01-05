package main

import (
	"goPingRobot/auth/internal/handler"
	"goPingRobot/auth/internal/repository"
	"goPingRobot/auth/internal/service"
	"goPingRobot/auth/proto"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetFormatter(new(logrus.TextFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error init configs: %s", err.Error())
	}

		if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

		db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBname:   viper.GetString("db.dbname"),
		SSLmode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

		if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
    authHandler := handler.NewAuthHandler(authService)

	grpcServer := grpc.NewServer()
    proto.RegisterAuthServiceServer(grpcServer, authHandler)

    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        logrus.Fatalf("Failed to start gRPC server: %v", err)
    }

    log.Println("Auth service is running on port 50051")
    if err := grpcServer.Serve(listener); err != nil {
        logrus.Fatalf("Failed to serve gRPC server: %v", err)
    }
}

func initConfig() error {
	viper.AddConfigPath("auth/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}