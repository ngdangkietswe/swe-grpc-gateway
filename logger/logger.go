package logger

import "github.com/ngdangkietswe/swe-go-common-shared/logger"

func NewZapLogger() (*logger.Logger, error) {
	instance, err := logger.NewLogger(
		"swe-gateway-service",
		"local",
		"debug",
		"logs/swe-gateway-service.log",
	)

	if err != nil {
		return nil, err
	}

	return instance, nil
}
