package keycloak_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	keycloak "github.com/stillya/testcontainers-keycloak"
)

var keycloakContainer *keycloak.KeycloakContainer

func Test_Example(t *testing.T) {
	ctx := context.Background()

	authServerURL, err := keycloakContainer.GetAuthServerURL(ctx)
	if err != nil {
		t.Errorf("GetAuthServerURL() error = %v", err)
		return
	}

	fmt.Println(authServerURL)
	// Output:
	// http://localhost:32768/auth
}

func TestMain(m *testing.M) {
	defer func() {
		if r := recover(); r != nil {
			shutDown()
			fmt.Println("Panic")
		}
	}()
	setup()
	code := m.Run()
	shutDown()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()
	keycloakContainer, err = keycloak.RunContainer(ctx)
	if err != nil {
		panic(err)
	}
}

func shutDown() {
	ctx := context.Background()
	err := keycloakContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}
