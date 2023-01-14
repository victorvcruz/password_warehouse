package handlers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user.com/internal/auth"
	"user.com/internal/user"
	"user.com/internal/utils"
	"user.com/internal/utils/errors"
	"user.com/pkg/pb"
)

type UserHandler struct {
	pb.UnimplementedUserServer
	userService user.UserServiceClient
	authService auth.AuthServiceClient
	validate    *validator.Validate
}

func NewUserHandler(
	_userService user.UserServiceClient,
	_authService auth.AuthServiceClient,
	_validate *validator.Validate,
) *UserHandler {
	return &UserHandler{
		userService: _userService,
		authService: _authService,
		validate:    _validate,
	}
}

func (u *UserHandler) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	user := user.UserRequest{Name: req.Name, Email: req.Email, MasterPassword: req.MasterPassword}

	err := u.validate.Struct(user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, utils.RequestUserValidate(err))
	}

	err = u.userService.CreateUser(user)
	if err != nil {
		switch err.(type) {
		case *errors.ConflictEmailError:
			return nil, status.Error(codes.AlreadyExists, "Email already exist")
		default:
			return nil, status.Error(codes.Internal, "Internal server error")
		}
	}

	return &pb.UserResponse{Name: req.Name, Email: req.Email}, nil
}

func (u *UserHandler) FindUser(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) {

	token, err := utils.BearerToken(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Add token")
	}

	id, err := u.authService.AuthToken(token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := u.userService.FindUser(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.UserResponse{Name: user.Name, Email: user.Email}, nil
}

func (u *UserHandler) UpdateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	token, err := utils.BearerToken(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Add token")
	}

	id, err := u.authService.AuthToken(token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := user.UserRequest{Name: req.Name, Email: req.Email, MasterPassword: req.MasterPassword}
	err = u.validate.Struct(user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, utils.RequestUserValidate(err))
	}

	err = u.userService.UpdateUser(id, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.UserResponse{Name: req.Name, Email: req.Email}, nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, req *pb.Empty) (*pb.MessageResponse, error) {

	token, err := utils.BearerToken(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Add token")
	}

	id, err := u.authService.AuthToken(token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = u.userService.DeleteUser(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.MessageResponse{Message: "User deleted"}, nil
}
