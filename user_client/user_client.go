package main

import (
	"context"
	"google.golang.org/grpc"
	"log"

	pb "example.com/grpc_jwt_crud_api/user_protobuf"
)

func main() {
	log.Println("Starting Client ")

	cc, err := grpc.Dial("localhost:5005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error Connecting: %v", err)
	}
	defer cc.Close()

	c := pb.NewUserServicesClient(cc)

	// User Create
	log.Println("Calling CreateUser() ")
	user := &pb.User{
		Uid:         3,
		Name:        "Mya Mya",
		Nationality: "Rakhine",
		Zip:         202,
	}
	createUserReponse, err := c.CreateUser(context.Background(), &pb.CreateUserRequest{User: user})
	if err != nil {
		log.Fatalf("Error occur in creating user : %v", err)
	}
	log.Printf("User has been created : %v \n", createUserReponse)

	// User Update
	log.Println("Calling UpdateUser() ")
	updateUser := &pb.User{
		Uid:         2,
		Name:        "Shwe Shwe",
		Nationality: "Burma",
		Zip:         90,
	}
	updatedUserResponse, err := c.UpdateUser(context.Background(), &pb.UpdateUserRequest{User: updateUser})
	if err != nil {
		log.Fatalf("Error occur in updating user : %v", err)
	}
	log.Printf("User has been updated: %v \n", updatedUserResponse)

	// Delete User
	log.Println("Calling DeleteUser() ")
	var uid int32 = 3
	deleteUserResponse, err := c.DeleteUser(context.Background(), &pb.DeleteUserRequest{Uid: uid})
	if err != nil {
		log.Fatalf("Error occur in deleting user: %v", err)
	}
	log.Printf("User had been deleted: %v\n", deleteUserResponse)

	// Fetch user
	log.Println("Calling FetchUser() ")
	uid = 1
	fetchUserResponse, err := c.FetchUser(context.Background(), &pb.FetchUserRequest{Uid: uid})
	if err != nil {
		log.Fatalf("Error occur in fetching user: %v", err)
	}
	log.Printf("Fetched user: %v\n", fetchUserResponse)

}
