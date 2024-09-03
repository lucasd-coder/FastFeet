package kecloak

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-resty/resty/v2"
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/domain/auth"
	"github.com/lucasd-coder/fast-feet/auth-service/internal/shared"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	spanErrRequest         = "Request Error"
	spanErrExtractResponse = "Error Extract Response"
)

type Repository struct {
	Config *config.Config
}

func NewRepository(cfg *config.Config) *Repository {
	return &Repository{
		Config: cfg,
	}
}

func (r *Repository) Register(ctx context.Context, pld *auth.Register) (string, error) {
	client := NewClient(ctx, r.Config)

	token, err := client.LoginAdmin(ctx, r.Config.KeyCloakUsername, r.Config.KeyCloakPassword, r.Config.KeyCloakRealm)
	if err != nil {
		return "", err
	}

	roles, err := r.rolesIdentity(ctx, auth.GetRolesFromString(pld.Roles))
	if err != nil {
		return "", err
	}
	user := gocloak.User{
		Username:  gocloak.StringP(pld.Username),
		FirstName: gocloak.StringP(pld.FirstName),
		LastName:  gocloak.StringP(pld.LastName),
		Email:     gocloak.StringP(pld.Username),
		Enabled:   gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(pld.Password),
			},
		},
		RealmRoles: roles,
	}

	id, err := client.CreateUser(ctx, token.AccessToken, r.Config.KeyCloakRealm, user)
	if err != nil {
		return "", err
	}

	if err := r.addRealmRoleToUser(ctx, id, *roles); err != nil {
		return "", err
	}

	return id, nil
}

func (repo *Repository) rolesIdentity(ctx context.Context, r auth.RolesEnum) (*[]string, error) {
	roles := []string{}
	switch r {
	case auth.Admin:
		roles = append(roles, auth.Admin.String(), auth.User.String())
	case auth.User:
		roles = append(roles, auth.User.String())
	default:
		roles = append(roles, auth.Unknown.String())
	}
	for _, role := range roles {
		if err := repo.ensureRoleExists(ctx, role); err != nil {
			return nil, err
		}
	}
	return &roles, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, pld *auth.FindUserByEmail) (*auth.UserRepresentation, error) {
	client := NewClient(ctx, r.Config)
	span := trace.SpanFromContext(ctx)
	token, err := client.LoginAdmin(ctx, r.Config.KeyCloakUsername, r.Config.KeyCloakPassword, r.Config.KeyCloakRealm)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetRequestWithBearerAuth(ctx, token.AccessToken).
		SetQueryParams(map[string]string{
			"email": pld.Email,
			"exact": "true",
		}).
		SetResult(&[]auth.UserRepresentation{}).
		Get(r.getURL("users"))
	if err := checkForError(resp, err); err != nil {
		return nil, err
	}

	result, ok := resp.Result().(*[]auth.UserRepresentation)
	if !ok {
		errMsg := fmt.Errorf("%w. Service: FindUserByEmail", shared.ErrExtractResponse)
		r.createSpanError(ctx, err, spanErrExtractResponse)
		return nil, errMsg
	}

	msg := fmt.Sprintf("keyCloak call successful. Service: FindUserByEmail, Response time: %s",
		resp.ReceivedAt().String())
	span.SetStatus(codes.Ok, msg)

	return extractUserRepresentation(result)
}

func (r *Repository) GetRoles(ctx context.Context, pld *auth.GetUserID) ([]string, error) {
	client := NewClient(ctx, r.Config)
	token, err := client.LoginAdmin(ctx, r.Config.KeyCloakUsername, r.Config.KeyCloakPassword, r.Config.KeyCloakRealm)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetRealmRolesByUserID(ctx, token.AccessToken, r.Config.KeyCloakRealm, pld.ID)
	if err != nil {
		r.createSpanError(ctx, err, spanErrRequest)
		return nil, err
	}

	return extractRoles(resp)
}

func (r *Repository) IsActiveUser(ctx context.Context, pld *auth.GetUserID) (bool, error) {
	user, err := r.findUserByID(ctx, pld.ID)
	if err != nil {
		return false, err
	}
	return *user.Enabled, nil
}

func (r *Repository) addRealmRoleToUser(ctx context.Context, userID string, roles []string) error {
	client := NewClient(ctx, r.Config)

	token, err := client.LoginAdmin(ctx, r.Config.KeyCloakUsername, r.Config.KeyCloakPassword, r.Config.KeyCloakRealm)
	if err != nil {
		return err
	}

	maxRoles := 100
	params := gocloak.GetRoleParams{
		Max: gocloak.IntP(maxRoles),
	}

	getRoles, err := client.GetRealmRoles(ctx, token.AccessToken, r.Config.KeyCloakRealm, params)
	if err != nil {
		return err
	}

	filterRoles := []gocloak.Role{}
	for _, role := range roles {
		for _, item := range getRoles {
			if strings.EqualFold(role, *item.Name) {
				filterRoles = append(filterRoles, *item)
			}
		}
	}
	return client.AddRealmRoleToUser(ctx, token.AccessToken,
		r.Config.KeyCloakRealm, userID, filterRoles)
}

func (r *Repository) findUserByID(ctx context.Context, userID string) (*gocloak.User, error) {
	client := NewClient(ctx, r.Config)
	token, err := client.LoginAdmin(ctx, r.Config.KeyCloakUsername, r.Config.KeyCloakPassword, r.Config.KeyCloakRealm)
	if err != nil {
		return nil, err
	}

	return client.GetUserByID(ctx, token.AccessToken,
		r.Config.KeyCloakRealm, userID)
}

func (r *Repository) getURL(path ...string) string {
	path = append([]string{r.Config.KeyCloakBaseURL, "admin/realms", r.Config.KeyCloakRealm}, path...)
	return makeURL(path...)
}

func makeURL(path ...string) string {
	return strings.Join(path, "/")
}

func checkForError(resp *resty.Response, err error) error {
	if err != nil {
		return &shared.HTTPError{
			StatusCode: 0,
			Message:    err.Error(),
		}
	}

	if resp == nil {
		return &shared.HTTPError{
			StatusCode: 500,
			Message:    "empty response",
		}
	}

	if resp.IsError() {
		var msg string

		if e, ok := resp.Error().(*gocloak.HTTPErrorResponse); ok && e.NotEmpty() {
			msg = fmt.Sprintf("%s: %s", resp.Status(), e)
		} else {
			msg = resp.Status()
		}

		return &shared.HTTPError{
			StatusCode: resp.StatusCode(),
			Message:    msg,
		}
	}

	return nil
}

func extractUserRepresentation(result *[]auth.UserRepresentation) (*auth.UserRepresentation, error) {
	if len(*result) == 0 {
		return nil, shared.ErrUserNotFound
	}

	for _, item := range *result {
		return &item, nil
	}
	return nil, nil
}

func (r *Repository) createSpanError(ctx context.Context, err error, msg string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, msg)
	span.RecordError(err)
}

func extractRoles(resp []*gocloak.Role) ([]string, error) {
	if len(resp) == 0 {
		return []string{}, nil
	}

	result := []string{}
	for _, item := range resp {
		result = append(result, *item.Name)
	}

	return result, nil
}

func (r *Repository) ensureRoleExists(ctx context.Context, roleName string) error {
	client := NewClient(ctx, r.Config)

	token, err := client.LoginAdmin(ctx, r.Config.KeyCloakUsername, r.Config.KeyCloakPassword, r.Config.KeyCloakRealm)
	if err != nil {
		return err
	}

	roles, err := client.GetRealmRoles(ctx, token.AccessToken, r.Config.KeyCloakRealm, gocloak.GetRoleParams{})
	if err != nil {
		return err
	}

	for _, role := range roles {
		if strings.EqualFold(roleName, *role.Name) {
			return nil
		}
	}

	_, err = client.CreateRealmRole(ctx, token.AccessToken, r.Config.KeyCloakRealm, gocloak.Role{
		Name: gocloak.StringP(roleName),
	})
	return err
}
