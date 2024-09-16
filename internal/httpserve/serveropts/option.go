package serveropts

import (
	"github.com/rs/zerolog/log"
	echoprometheus "github.com/theopenlane/echo-prometheus"
	echo "github.com/theopenlane/echox"
	"github.com/theopenlane/echox/middleware"

	"github.com/theopenlane/openlane-cloud/internal/httpserve/config"

	"github.com/theopenlane/core/pkg/middleware/cachecontrol"
	"github.com/theopenlane/core/pkg/middleware/cors"
	"github.com/theopenlane/core/pkg/middleware/mime"
	"github.com/theopenlane/core/pkg/middleware/ratelimit"
	"github.com/theopenlane/core/pkg/middleware/redirect"
	"github.com/theopenlane/core/pkg/openlaneclient"
	"github.com/theopenlane/echox/middleware/echocontext"
)

type ServerOption interface {
	apply(*ServerOptions)
}

type applyFunc struct {
	applyInternal func(*ServerOptions)
}

func (fso *applyFunc) apply(s *ServerOptions) {
	fso.applyInternal(s)
}

func newApplyFunc(apply func(option *ServerOptions)) *applyFunc {
	return &applyFunc{
		applyInternal: apply,
	}
}

// WithConfigProvider supplies the config for the server
func WithConfigProvider(cfgProvider config.ConfigProvider) ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		s.ConfigProvider = cfgProvider
	})
}

// WithOpenlaneClient supplies the openlane client for the server
func WithOpenlaneClient() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		var err error

		creds := openlaneclient.Authorization{
			BearerToken: s.Config.Settings.Server.Openlane.Token,
		}

		s.Config.Handler.OpenlaneClient, err = openlaneclient.NewWithDefaults(
			openlaneclient.WithCredentials(creds))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create openlane client")
		}
	})
}

// WithHTTPS sets up TLS config settings for the server
func WithHTTPS() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		if !s.Config.Settings.Server.TLS.Enabled {
			// this is set to enabled by WithServer
			// if TLS is not enabled, move on
			return
		}

		s.Config.WithTLSDefaults()

		if !s.Config.Settings.Server.TLS.AutoCert {
			s.Config.WithTLSCerts(s.Config.Settings.Server.TLS.CertFile, s.Config.Settings.Server.TLS.CertKey)
		}
	})
}

// WithMiddleware adds the middleware to the server
func WithMiddleware() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		// Initialize middleware if null
		if s.Config.DefaultMiddleware == nil {
			s.Config.DefaultMiddleware = []echo.MiddlewareFunc{}
		}

		// default middleware
		s.Config.DefaultMiddleware = append(s.Config.DefaultMiddleware,
			middleware.RequestID(), // add request id
			middleware.Recover(),   // recover server from any panic/fatal error gracefully
			middleware.LoggerWithConfig(middleware.LoggerConfig{
				Output: log.Logger.Hook(LevelNameHook{}),
				Format: "remote_ip=${remote_ip}, method=${method}, uri=${uri}, status=${status}, session=${header:Set-Cookie}, host=${host}, referer=${referer}, user_agent=${user_agent}, route=${route}, path=${path}",
			}),
			echoprometheus.MetricsMiddleware(),                                                       // add prometheus metrics
			echocontext.EchoContextToContextMiddleware(),                                             // adds echo context to parent
			cors.New(s.Config.Settings.Server.CORS.AllowOrigins),                                     // add cors middleware
			mime.NewWithConfig(mime.Config{DefaultContentType: echo.MIMEApplicationJSONCharsetUTF8}), // add mime middleware
			cachecontrol.New(),                        // add cache control middleware
			middleware.Secure(),                       // add XSS middleware
			redirect.NewWithConfig(redirect.Config{}), // add redirect middleware
		)
	})
}

// WithRateLimiter sets up the rate limiter for the server
func WithRateLimiter() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		if s.Config.Settings.Ratelimit.Enabled {
			s.Config.DefaultMiddleware = append(s.Config.DefaultMiddleware, ratelimit.RateLimiterWithConfig(&s.Config.Settings.Ratelimit))
		}
	})
}
