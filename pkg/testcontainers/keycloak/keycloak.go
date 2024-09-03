package keycloak

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	// keycloak "github.com/stillya/testcontainers-keycloak"

	"github.com/lucasd-coder/fast-feet/pkg/testcontainers/keycloak/container"
	"github.com/testcontainers/testcontainers-go"
)

func RunContainer(ctx context.Context) (*container.KeycloakContainer, error) {
	testDataPath, err := FindTestDataDir()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(testDataPath, "realm-export.json")
	return container.RunContainer(ctx,
		WithCustomOption(),
		container.WithContextPath("/auth"),
		container.WithRealmImportFile(fullPath),
		container.WithAdminUsername("admin"),
		container.WithAdminPassword("admin"),
		testcontainers.WithEnv(map[string]string{
			"KEYCLOAK_LOGLEVEL": "DEBUG",
			"KEYCLOAK_USER":     "admin",
			"KEYCLOAK_PASSWORD": "admin",
		}),
	)
}

func WithCustomOption() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Image = "quay.io/keycloak/keycloak:22.0.1-2"
		return nil
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
