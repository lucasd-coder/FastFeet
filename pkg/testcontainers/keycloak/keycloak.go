package keycloak

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	keycloak "github.com/stillya/testcontainers-keycloak"
	"github.com/testcontainers/testcontainers-go"
)

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	testDataPath, err := FindTestDataDir()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(testDataPath, "realm-export.json")
	return keycloak.RunContainer(ctx,
		WithCustomOption(),
		keycloak.WithContextPath("/auth"),
		keycloak.WithRealmImportFile(fullPath),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
}

func WithCustomOption() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Image = "quay.io/keycloak/keycloak:22.0.1-2"
	}
}

func FindTestDataDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return findTestDataInInfraRecursive(currentDir)
}

func findTestDataInInfraRecursive(dir string) (string, error) {
	if dir == "" || dir == "/" {
		return "", errors.New("testdata folder not found")
	}

	infraTestDataPath, err := findTestDataInInfra(dir)
	if err != nil {
		return findTestDataInInfraRecursive(filepath.Dir(dir))
	}

	return infraTestDataPath, nil
}

func findTestDataInInfra(dir string) (string, error) {
	infraTestDataPath := filepath.Join(dir, "infra", "testdata")
	if _, err := os.Stat(infraTestDataPath); err == nil {
		return infraTestDataPath, nil
	}

	return "", errors.New("testdata folder not found")
}
