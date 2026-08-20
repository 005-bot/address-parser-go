package address

import (
	"context"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module returns an fx option that provides *Parser and stops it on
// application shutdown. A Config provider must be supplied by the app.
func Module() fx.Option {
	return fx.Module(
		"address-parser",
		logger.WithNamedLogger("address-parser"),
		fx.Provide(NewParser),
		fx.Invoke(func(lc fx.Lifecycle, p *Parser, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return nil
				},
				OnStop: func(_ context.Context) error {
					p.Stop()
					logger.Info("address parser stopped")
					return nil
				},
			})
		}),
	)
}
