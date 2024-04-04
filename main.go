package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	UserPb "rahulchhabra.io/proto/users"
)

type User struct {
	gorm.Model
	Firstname string
	Lastname  string
	Age       int64
	Username  string
	Password  string
}
type Config struct {
	UserPb.UnimplementedUserServiceServer
}

var dbConnector *gorm.DB

func (*Config) CreateUser(ctx context.Context, request *UserPb.CreateUserRequest) (response *UserPb.CreateUserResponse, err error) {
	username := request.GetUsername()
	/**
		Check if user already exist...
	**/
	var existingUser User
	userNotFoundError := dbConnector.Where("username = ?", username).First(&existingUser).Error
	if userNotFoundError != nil {
		password := request.GetPassword()
		firstname := request.GetFirstname()
		lastname := request.GetLastname()
		age := request.GetAge()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Couldn't hash password and the error is %s", err)
		}
		/** Creating a new User.
		**/
		newUser := &User{Username: username, Password: string(hashedPassword), Firstname: firstname, Lastname: lastname, Age: age}

		primaryKey := dbConnector.Create(newUser)
		if primaryKey.Error != nil {
			return nil, primaryKey.Error
		}
		return &UserPb.CreateUserResponse{
			Username: username,
		}, nil
	}
	fmt.Println("Username is used, please try another username ")
	return nil, userNotFoundError
}
func (*Config) LoginUser(ctx context.Context, request *UserPb.LoginUserRequest) (response *UserPb.LoginUserResponse, err error) {
	username := request.GetUsername()
	password := request.GetPassword()
	var existingUser User
	if err := dbConnector.Where("username = ?", username).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Record not found")
			return nil, err
		} else {
			fmt.Println("Error:", err)
			return nil, err
		}
	} else {
		// ** Compare Passwords....
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				fmt.Printf("Password does not match %s", err)
				return nil, err
			} else {
				log.Fatalf("Error : %s", err)
				return nil, err
			}
		}
		return &UserPb.LoginUserResponse{
			Username:  existingUser.Username,
			IsSuccess: true,
			Message:   "Logged In Successfully",
		}, nil
	}
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
