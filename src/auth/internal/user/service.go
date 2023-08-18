package user

import (
	"auth.com/internal/utils/authorization"
	"context"
	"github.com/victorvcruz/password_warehouse/protobuf/user_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"strconv"
)

const UserService = "USER"

type UserServiceClient interface {
	UserById(id int64) (*UserDTO, error)
	UserByEmail(email string) (*UserDTO, error)
}

type userService struct {
	client        user_pb.UserClient
	ctx           context.Context
	authorization authorization.AuthorizationClient
}

func NewUserService(_authorization authorization.AuthorizationClient) UserServiceClient {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(os.Getenv("USER_HOST"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	_client := user_pb.NewUserClient(conn)
	return &userService{
		client:        _client,
		authorization: _authorization,
		ctx:           context.Background(),
	}
}

func (u *userService) UserByEmail(email string) (*UserDTO, error) {
	apiToken, err := u.authorization.Login(UserService)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"userEmail": email, "api-token": apiToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	user, err := u.client.FindUserByData(ctx, &user_pb.Empty{})
	if err != nil {
		return nil, err
	}

	return u.userResponseToDto(user), nil
}

func (u *userService) UserById(id int64) (*UserDTO, error) {
	apiToken, err := u.authorization.Login(UserService)
	if err != nil {
		return nil, err
	}

	md := metadata.New(map[string]string{"id": strconv.FormatInt(id, 10), "api-token": apiToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	user, err := u.client.FindUserByData(ctx, &user_pb.Empty{})
	if err != nil {
		return nil, err
	}

	return u.userResponseToDto(user), nil
}

func (u *userService) userResponseToDto(response *user_pb.DetailedUserResponse) *UserDTO {
	return &UserDTO{
		Id:             response.Id,
		Name:           response.Name,
		Email:          response.Email,
		MasterPassword: response.MasterPassword,
		CreatedAt:      response.CreatedAt.AsTime(),
		UpdatedAt:      response.UpdatedAt.AsTime(),
		DeletedAt:      response.DeletedAt.AsTime(),
	}
}