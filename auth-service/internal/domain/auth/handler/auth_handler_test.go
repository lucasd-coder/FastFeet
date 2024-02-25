package handler_test

import (
	"context"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth/handler"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/kecloak"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/provider/validator"
	"github.com/lucasd-coder/fast-feet/auth-service/pkg/pb"
	testcontainers "github.com/lucasd-coder/fast-feet/pkg/testcontainers/keycloak"
	keycloak "github.com/stillya/testcontainers-keycloak"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type AuthHandlerSuite struct {
	suite.Suite
	srv               *grpc.Server
	lis               *bufconn.Listener
	cfg               config.Config
	authHandler       *handler.AuthHandler
	ctx               context.Context
	keycloakContainer *keycloak.KeycloakContainer
	conn              *grpc.ClientConn
}

func (suite *AuthHandlerSuite) SetupSuite() {
	suite.ctx = context.Background()
	var err error
	suite.keycloakContainer, err = testcontainers.RunContainer(suite.ctx)
	if err != nil {
		suite.T().Fatal(err)
	}

	authServerURL, err := suite.keycloakContainer.GetAuthServerURL(suite.ctx)
	if err != nil {
		suite.T().Errorf("GetAuthServerURL() error = %v", err)
		return
	}

	baseDir, err := os.Getwd()
	if err != nil {
		suite.T().Errorf("os.Getwd() error = %v", err)
		return
	}

	staticDir := filepath.Join(baseDir, "..", "..", "..", "..", "/config/config-test.yml")

	slog.Info("config lod", "dir", staticDir)
	err = cleanenv.ReadConfig(staticDir, &suite.cfg)
	if err != nil {
		suite.T().Errorf("cleanenv.ReadConfig() error = %v", err)
		return
	}
	suite.cfg.KeyCloakBaseURL = authServerURL
	config.ExportConfig(&suite.cfg)

	client := kecloak.NewClient(suite.ctx, &suite.cfg)

	token, err := client.LoginAdmin(suite.ctx, "admin", "admin", "master")
	if err != nil {
		suite.T().Fatalf("get token %v", err)
	}

	createRoles := []gocloak.Role{
		{
			Name: gocloak.StringP("user"),
		},
	}

	for _, role := range createRoles {
		roleName, err := client.CreateRealmRole(suite.ctx, token.AccessToken, suite.cfg.KeyCloakRealm, role)
		if err != nil {
			suite.T().Fatalf("CreateRealmRole %v", err)
		}

		slog.Info("CreateRealmRole", "roleName", roleName)
	}
}

func (suite *AuthHandlerSuite) TearDownSuite() {
	if err := suite.keycloakContainer.Terminate(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *AuthHandlerSuite) SetupTest() {
	lis := bufconn.Listen(1024 * 1024)
	suite.T().Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	suite.T().Cleanup(func() {
		srv.Stop()
	})

	val := validator.NewValidation()
	repo := kecloak.NewRepository(&suite.cfg)
	svc := auth.NewService(val, repo)
	hdler := handler.NewHandler(svc, &suite.cfg)
	suite.authHandler = handler.NewAuthHandler(*hdler)

	pb.RegisterRegisterHandlerServer(srv, suite.authHandler)
	pb.RegisterAuthHandlerServer(srv, suite.authHandler)

	go func() {
		if err := srv.Serve(lis); err != nil {
			suite.T().Error(err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(suite.ctx, 3*time.Second)
	suite.T().Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.T().Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		suite.T().Fatalf("grpc.DialContext %v", err)
	}
	suite.conn = conn
}

func (suite *AuthHandlerSuite) TestCreateUser() {
	client := pb.NewRegisterHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "maria@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")
}

func (suite *AuthHandlerSuite) TestCreateUser_ValidateFailure() {
	client := pb.NewRegisterHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{}

	_, err := client.CreateUser(ctx, in)
	suite.Error(err)
	suite.Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *AuthHandlerSuite) TestFindUserByEmail() {
	client := pb.NewRegisterHandlerClient(suite.conn)
	client2 := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "maria12@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")

	md := metadata.New(map[string]string{"email": in.GetUsername()})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	findUser, err := client2.FindUserByEmail(ctx, &inEmpty)
	suite.NoError(err)

	suite.Equal(findUser.GetEmail(), in.GetUsername())
	suite.True(findUser.GetEnabled())
}

func (suite *AuthHandlerSuite) TestFindUserByEmail_ValidateFailure() {
	client := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	md := metadata.New(map[string]string{"email": "invalid email"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	_, err := client.FindUserByEmail(ctx, &inEmpty)

	suite.Error(err)
	suite.Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *AuthHandlerSuite) TestFindUserByEmail_UserNotFound() {
	client := pb.NewRegisterHandlerClient(suite.conn)
	client2 := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "maria123@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")

	md := metadata.New(map[string]string{"email": "joao123@gmail.com"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	_, err = client2.FindUserByEmail(ctx, &inEmpty)
	suite.Error(err)
	suite.Equal(status.Code(err), codes.NotFound)
}

func (suite *AuthHandlerSuite) TestGetRoles() {
	client := pb.NewRegisterHandlerClient(suite.conn)
	client2 := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "maria1234@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")

	md := metadata.New(map[string]string{"id": resp.GetId()})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	getRoles, err := client2.GetRoles(ctx, &inEmpty)
	suite.NoError(err)

	suite.NotEmpty(getRoles.GetRoles())
	suite.Contains(getRoles.GetRoles(), strings.ToLower(in.GetAuthority().String()))
}

func (suite *AuthHandlerSuite) TestGetRoles_ValidateFailure() {
	client := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	md := metadata.New(map[string]string{"id": "id invalid"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	_, err := client.GetRoles(ctx, &inEmpty)

	suite.Error(err)
	suite.Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *AuthHandlerSuite) TestGetRoles_UserNotFound() {
	client := pb.NewRegisterHandlerClient(suite.conn)
	client2 := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "maria12345@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")

	md := metadata.New(map[string]string{"id": "fdac84ee-7a63-48dd-906f-2bb13f5b6b8e"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	_, err = client2.GetRoles(ctx, &inEmpty)

	suite.Error(err)
	suite.Equal(status.Code(err), codes.NotFound)
}

func (suite *AuthHandlerSuite) TestIsActiveUser() {
	client := pb.NewRegisterHandlerClient(suite.conn)
	client2 := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "joao@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")

	md := metadata.New(map[string]string{"id": resp.GetId()})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	isActiveUser, err := client2.IsActiveUser(ctx, &inEmpty)
	suite.NoError(err)

	suite.True(isActiveUser.GetActive())
}

func (suite *AuthHandlerSuite) TestIsActiveUser_ValidateFailure() {
	client := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	md := metadata.New(map[string]string{"id": "id invalid"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	_, err := client.IsActiveUser(ctx, &inEmpty)
	suite.Error(err)
	suite.Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *AuthHandlerSuite) TestIsActiveUser_UserNotFound() {
	client := pb.NewRegisterHandlerClient(suite.conn)
	client2 := pb.NewAuthHandlerClient(suite.conn)

	ctx := suite.ctx

	in := &pb.RegisterRequest{
		FirstName: "Maria",
		Username:  "maria12345678@gmail.com",
		Password:  "12345@#&",
		Authority: pb.Roles_USER,
	}

	resp, err := client.CreateUser(ctx, in)
	suite.NoError(err)
	suite.NotEqual(resp.GetId(), "")

	md := metadata.New(map[string]string{"id": "fdac84ee-7a63-48dd-906f-2bb13f5b6b8e"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	inEmpty := pb.EmptyRequest{}

	_, err = client2.IsActiveUser(ctx, &inEmpty)
	suite.Error(err)
	suite.Equal(status.Code(err), codes.NotFound)
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}
