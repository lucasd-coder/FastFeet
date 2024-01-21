package rabbitmq_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/lucasd-coder/fast-feet/pkg/testcontainers/rabbitmq"
	"github.com/testcontainers/testcontainers-go"
)

func TestRunContainer_withAllSettings(t *testing.T) {
	ctx := context.Background()

	container, err := rabbitmq.RunContainer(ctx, "queue1", "exchange1")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	if !assertEntity(t, container, "queues", "queue1") {
		t.Fatal("error verify queues")
	}

	if !assertEntity(t, container, "exchanges", "direct", "exchange1") {
		t.Fatal("error verify exchanges")
	}
}

func assertEntity(t *testing.T, container testcontainers.Container, listCommand string, entities ...string) bool {
	t.Helper()

	ctx := context.Background()

	cmd := []string{"rabbitmqadmin", "list", listCommand}

	_, out, err := container.Exec(ctx, cmd)
	if err != nil {
		t.Fatal(err)
	}

	check, err := io.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range entities {
		if !strings.Contains(string(check), e) {
			return false
		}
	}

	return true
}
