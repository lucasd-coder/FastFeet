package kecloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/lucasd-coder/fast-feet/auth-service/config"
	"github.com/lucasd-coder/fast-feet/pkg/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewClient(ctx context.Context, cfg *config.Config) *gocloak.GoCloak {
	log := logger.FromContext(ctx)
	client := gocloak.NewClient(cfg.KeyCloakBaseURL,
		gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	client.RestyClient().
		SetDebug(*cfg.KeyCloakDebug).
		SetTimeout(cfg.KeyCloakRequestTimeout).
		SetLogger(log).
		EnableTrace().
		SetTransport(otelhttp.NewTransport(
			client.RestyClient().GetClient().Transport))

	return client
}
