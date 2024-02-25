package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func RunContainer(ctx context.Context, queueName, exchangeName string) (*rabbitmq.RabbitMQContainer, error) {
	return rabbitmq.RunContainer(ctx,
		testcontainers.WithImage("rabbitmq:3.13.0-management-alpine"),
		rabbitmq.WithAdminUsername("admin"),
		rabbitmq.WithAdminPassword("password"),
		testcontainers.WithAfterReadyCommand(VirtualHost{
			Name:    "fastfeet",
			Tracing: true,
		}),
		testcontainers.WithAfterReadyCommand(Queue{
			Name:  queueName,
			VHost: "fastfeet",
		}),
		testcontainers.WithAfterReadyCommand(Exchange{
			Name:  exchangeName,
			Type:  "direct",
			VHost: "fastfeet",
		}),
		testcontainers.WithAfterReadyCommand(NewBindingWithVHost("fastfeet", exchangeName, queueName)),
	)
}

// --------- Exchange ---------

type Exchange struct {
	testcontainers.ExecOptions
	Name       string
	VHost      string
	Type       string
	AutoDelete bool
	Internal   bool
	Durable    bool
	Args       map[string]interface{}
}

func (e Exchange) AsCommand() []string {
	cmd := []string{"rabbitmqadmin"}

	if e.VHost != "" {
		cmd = append(cmd, "--vhost="+e.VHost)
	}

	cmd = append(cmd, "declare", "exchange", fmt.Sprintf("name=%s", e.Name), fmt.Sprintf("type=%s", e.Type))

	if e.AutoDelete {
		cmd = append(cmd, "auto_delete=true")
	}
	if e.Internal {
		cmd = append(cmd, "internal=true")
	}
	if e.Durable {
		cmd = append(cmd, fmt.Sprintf("durable=%t", e.Durable))
	}

	if len(e.Args) > 0 {
		bytes, err := json.Marshal(e.Args)
		if err != nil {
			return cmd
		}

		cmd = append(cmd, "arguments="+string(bytes))
	}

	return cmd
}

// --------- Queue ---------

type Queue struct {
	testcontainers.ExecOptions
	Name       string
	VHost      string
	AutoDelete bool
	Durable    bool
	Args       map[string]interface{}
}

func (q Queue) AsCommand() []string {
	cmd := []string{"rabbitmqadmin"}

	if q.VHost != "" {
		cmd = append(cmd, "--vhost="+q.VHost)
	}

	cmd = append(cmd, "declare", "queue", fmt.Sprintf("name=%s", q.Name))

	if q.AutoDelete {
		cmd = append(cmd, "auto_delete=true")
	}
	if q.Durable {
		cmd = append(cmd, fmt.Sprintf("durable=%t", q.Durable))
	}

	if len(q.Args) > 0 {
		bytes, err := json.Marshal(q.Args)
		if err != nil {
			return cmd
		}

		cmd = append(cmd, "arguments="+string(bytes))
	}

	return cmd
}

// --------- Bindings ---------

type Binding struct {
	testcontainers.ExecOptions
	VHost           string
	Source          string
	Destination     string
	DestinationType string
	RoutingKey      string
	// additional arguments, that will be serialized to JSON when passed to the container
	Args map[string]interface{}
}

func NewBinding(source string, destination string) Binding {
	return Binding{
		Source:      source,
		Destination: destination,
	}
}

func NewBindingWithVHost(vhost string, source string, destination string) Binding {
	return Binding{
		VHost:       vhost,
		Source:      source,
		Destination: destination,
	}
}

func (b Binding) AsCommand() []string {
	cmd := []string{"rabbitmqadmin"}

	if b.VHost != "" {
		cmd = append(cmd, fmt.Sprintf("--vhost=%s", b.VHost))
	}

	cmd = append(cmd, "declare", "binding", fmt.Sprintf("source=%s", b.Source), fmt.Sprintf("destination=%s", b.Destination))

	if b.DestinationType != "" {
		cmd = append(cmd, fmt.Sprintf("destination_type=%s", b.DestinationType))
	}
	if b.RoutingKey != "" {
		cmd = append(cmd, fmt.Sprintf("routing_key=%s", b.RoutingKey))
	}

	if len(b.Args) > 0 {
		bytes, err := json.Marshal(b.Args)
		if err != nil {
			return cmd
		}

		cmd = append(cmd, "arguments="+string(bytes))
	}

	return cmd
}

// --------- Virtual Hosts --------

type VirtualHost struct {
	testcontainers.ExecOptions
	Name    string
	Tracing bool
}

func (v VirtualHost) AsCommand() []string {
	cmd := []string{"rabbitmqadmin", "declare", "vhost", fmt.Sprintf("name=%s", v.Name)}

	if v.Tracing {
		cmd = append(cmd, "tracing=true")
	}

	return cmd
}
