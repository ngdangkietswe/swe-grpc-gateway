package logger

import (
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"go.uber.org/fx"
)

func NewZapLogger() (*logger.Logger, error) {
	instance, err := logger.NewLogger(
		"swe-gateway-service",
		"local",
		"debug",
		"",
	)

	if err != nil {
		return nil, err
	}

	return instance, nil
}

var Module = fx.Provide(NewZapLogger)
