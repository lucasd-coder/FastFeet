package keycloak

import (
	"context"

	keycloak "github.com/stillya/testcontainers-keycloak"
	"github.com/testcontainers/testcontainers-go"
)

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	return keycloak.RunContainer(ctx,
		WithCustomOption(),
		keycloak.WithContextPath("/auth"),
		keycloak.WithRealmImportFile("../../../infra/keycloak/quarkus-realm.json"),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
}

func WithCustomOption() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Image = "quay.io/keycloak/keycloak:22.0.1-2"
	}
}
