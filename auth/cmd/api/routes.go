package api

import (
	"auth.com/cmd/api/handlers"
	"auth.com/pkg/pb"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
)

func New(user *handlers.AuthHandler) error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("API_PORT")))
	if err != nil {
		return err
	}

	app := grpc.NewServer()
	pb.RegisterAuthServer(app, user)

	err = app.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}
