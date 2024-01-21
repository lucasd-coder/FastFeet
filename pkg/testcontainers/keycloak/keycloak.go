package keycloak

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	keycloak "github.com/stillya/testcontainers-keycloak"
)

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	testDataPath, err := FindTestDataDir()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(testDataPath, "realm-export.json")
	return keycloak.RunContainer(ctx,
		keycloak.WithContextPath("/auth"),
		keycloak.WithRealmImportFile(fullPath),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
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
