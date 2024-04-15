package cep

import (
	"context"

	"github.com/google/wire"
	"github.com/lucasd-coder/fast-feet/business-service/config"
	"github.com/lucasd-coder/fast-feet/business-service/internal/shared"
)

const (
	spanErrRequest         = "Request Error"
	spanErrResponseStatus  = "Response Status Error"
	spanErrExtractResponse = "Error Extract Response"
)

type Repository interface {
	GetAddress(ctx context.Context, cep string) (*shared.AddressResponse, error)
}

func NewRepository(cfg *config.Config) Repository {
	if *cfg.BrasilAbertoEnabled {
		wire.Build(wire.InterfaceValue(new(Repository), NewBrasilAbertoRepository))
	}
	if *cfg.ViaCepEnabled {
		wire.Build(wire.InterfaceValue(new(Repository), NewViaCepRepository))
	}
	return nil
}
