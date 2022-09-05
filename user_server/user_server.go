package main

import (
	"context"
	"errors"
	pb "example.com/grpc_jwt_crud_api/user_protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"math"
	"net"
	"time"
)

type server struct {
	pb.UserServicesServer
}

type UserDetails struct {
	Uid         int32
	Name        string
	Nationality string
	Zip         int32
}

var users = []UserDetails{
	{
		Uid:         1,
		Name:        "Josh Winters",
		Nationality: "American",
		Zip:         10111,
	},
	{
		Uid:         2,
		Name:        "Brian Stone",
		Nationality: "British",
		Zip:         20212,
	},
}

func ToUser(data UserDetails) *pb.User {
	return &pb.User{
		Uid:         data.Uid,
		Name:        data.Name,
		Nationality: data.Nationality,
		Zip:         data.Zip,
	}
}

func FromUser(user *pb.User) UserDetails {
	return UserDetails{
		Uid:         user.GetUid(),
		Name:        user.GetName(),
		Nationality: user.GetNationality(),
		Zip:         user.GetZip(),
	}
}

func (server *server) FetchUser(ctx context.Context, req *pb.FetchUserRequest) (*pb.FetchUserResponse, error) {
	log.Println("Fetching User")
	uid := req.GetUid()

	for _, user := range users {
		if user.Uid == uid {
			return &pb.FetchUserResponse{
				User: ToUser(user),
			}, nil
		}
	}

	return nil, errors.New("User not found ")
}

func (server *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("Creating User")
	user := req.GetUser()

	data := FromUser(user)

	users = append(users, data)
	return &pb.CreateUserResponse{
		User: ToUser(data),
	}, nil
}

func (server *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Println("Updating User")
	user := req.GetUser()

	data := FromUser(user)
	for i, user := range users {
		if user.Uid == data.Uid {
			users[i] = data
			return &pb.UpdateUserResponse{
				User: ToUser(data),
			}, nil
		}
	}
	return nil, errors.New("Couldn't update user ")
}

func (server *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Println("Deleting User")
	uid := req.GetUid()
	var tempUsers []UserDetails

	for i, user := range users {
		if user.Uid == uid {
			tempUsers = append(users[:i], users[i+1:]...)
			users = tempUsers
			log.Printf("User with %d had been deleted.\n", uid)
			return &pb.DeleteUserResponse{
				Uid: uid,
			}, nil
		}
	}
	return nil, errors.New("User Doesn't exist ")
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:5005")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt64),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Timeout: 5 * time.Second,
			}),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterUserServicesServer(s, &server{})

	log.Println("Starting server ... ")
	log.Printf("Hosting server on : %s \n", lis.Addr().String())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}
}
