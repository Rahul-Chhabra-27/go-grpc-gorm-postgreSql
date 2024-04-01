package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	UserPb "rahulchhabra.io/proto/users"
)

type User struct {
	gorm.Model
	Username string
	Password string
}
type Config struct {
	UserPb.UnimplementedUserServiceServer
}

var dbConnector *gorm.DB

func (*Config) CreateUser(ctx context.Context, request *UserPb.CreateUserRequest) (response *UserPb.CreateUserResponse, err error) {
	username := request.GetUsername()
	password := request.GetPassword()

	newUser := &User{Username: username, Password: password}
	primaryKey := dbConnector.Create(newUser)
	if primaryKey.Error != nil {
		return nil, primaryKey.Error
	}
	return &UserPb.CreateUserResponse{
		Username: username,
	}, nil
}

func main() {
	dbConnector = Connect()
	grpcServer := grpc.NewServer()
	listner, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Couldn't start grpc server")
	}
	UserPb.RegisterUserServiceServer(grpcServer, &Config{})

	reflection.Register(grpcServer)
	if err = grpcServer.Serve(listner); err != nil {
		log.Fatalf("Couldn't connect with grpc server")
	}
}
