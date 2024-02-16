/*
Package proxy provides a rate-limit proxy middleware using the github.com/juju/ratelimit lib.

Sample backend extra config

	...
	"extra_config": {
		...
		"github.com/devopsfaith/krakend-ratelimit/rate/proxy": {
			"max_rate": 100,
			"capacity": 100
		},
		...
	},
	...

Adding the middleware to your proxy stack

	import juju "github.com/davron112/krakend-ratelimit/juju/proxy"

	...

	var p proxy.Proxy
	var backend *config.Backend

	...

	p = juju.NewMiddleware(backend)(p)

	...

The ratelimit package provides an efficient token bucket implementation. See https://github.com/juju/ratelimit
and http://en.wikipedia.org/wiki/Token_bucket for more details.
*/
package proxy

import (
	"context"
	"errors"

	"github.com/davron112/lura/v2/config"
	"github.com/davron112/lura/v2/logging"
	"github.com/davron112/lura/v2/proxy"

	krakendrate "github.com/davron112/krakend-ratelimit/v2"
	"github.com/davron112/krakend-ratelimit/v2/juju"
)

// Namespace is the key to use to store and access the custom config data for the proxy
const Namespace = "github.com/devopsfaith/krakend-ratelimit/juju/proxy"

// Config is the custom config struct containing the params for the limiter
type Config struct {
	MaxRate  float64
	Capacity int64
}

// BackendFactory adds a ratelimiting middleware wrapping the internal factory
func BackendFactory(logger logging.Logger, next proxy.BackendFactory) proxy.BackendFactory {
	return func(cfg *config.Backend) proxy.Proxy {
		return NewMiddleware(logger, cfg)(next(cfg))
	}
}

// NewMiddleware builds a middleware based on the extra config params or fallbacks to the next proxy
func NewMiddleware(logger logging.Logger, remote *config.Backend) proxy.Middleware {
	logPrefix := "[BACKEND: " + remote.URLPattern + "][Ratelimit]"
	cfg, err := ConfigGetter(remote.ExtraConfig)
	if err != nil {
		if err != ErrNoExtraCfg {
			logger.Error(logPrefix, err)
		}
		return proxy.EmptyMiddleware
	}
	if cfg.MaxRate <= 0 {
		return proxy.EmptyMiddleware
	}
	tb := juju.NewLimiter(cfg.MaxRate, cfg.Capacity)
	logger.Debug(logPrefix, "Enabling the rate limiter")
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}
		return func(ctx context.Context, request *proxy.Request) (*proxy.Response, error) {
			if !tb.Allow() {
				return nil, krakendrate.ErrLimited
			}
			return next[0](ctx, request)
		}
	}
}

// ZeroCfg is the zero value for the Config struct
var ZeroCfg = Config{}

var (
	ErrNoExtraCfg    = errors.New("no extra config")
	ErrWrongExtraCfg = errors.New("wrong extra config")
)

// ConfigGetter parses the extra config for the rate adapter and returns
// a ZeroCfg and an error if something goes wrong.
func ConfigGetter(e config.ExtraConfig) (Config, error) {
	v, ok := e[Namespace]
	if !ok {
		return ZeroCfg, ErrNoExtraCfg
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return ZeroCfg, ErrWrongExtraCfg
	}
	cfg := Config{}
	if v, ok := tmp["max_rate"]; ok {
		switch val := v.(type) {
		case float64:
			cfg.MaxRate = val
		case int:
			cfg.MaxRate = float64(val)
		case int64:
			cfg.MaxRate = float64(val)
		}
	}
	if v, ok := tmp["capacity"]; ok {
		switch val := v.(type) {
		case int64:
			cfg.Capacity = val
		case int:
			cfg.Capacity = int64(val)
		case float64:
			cfg.Capacity = int64(val)
		}
	}
	return cfg, nil
}
