/*
Package router provides several rate-limit routers using the golang.org/x/time/rate lib.

Sample endpoint extra config

	...
	"extra_config": {
		...
		"github.com/devopsfaith/krakend-ratelimit/rate/router": {
			"maxRate": 2000,
			"strategy": "header",
			"clientMaxRate": 100,
			"key": "X-Private-Token",
		},
		...
	},
	...

The ratelimit package provides an efficient token bucket implementation. See https://golang.org/x/time/rate
and http://en.wikipedia.org/wiki/Token_bucket for more details.
*/
package router

import (
	"fmt"

	"github.com/davron112/lura/v2/config"
)

// Namespace is the key to use to store and access the custom config data for the router
const Namespace = "github.com/devopsfaith/krakend-ratelimit/rate/router"

// Config is the custom config struct containing the params for the router middlewares
type Config struct {
	MaxRate       int
	Strategy      string
	ClientMaxRate int
	Key           string
}

// ZeroCfg is the zero value for the Config struct
var ZeroCfg = Config{}

// ConfigGetter implements the config.ConfigGetter interface. It parses the extra config for the
// rate adapter and returns a ZeroCfg if something goes wrong.
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return ZeroCfg
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return ZeroCfg
	}
	cfg := Config{}
	if v, ok := tmp["maxRate"]; ok {
		switch val := v.(type) {
		case int:
			cfg.MaxRate = val
		case float64:
			cfg.MaxRate = int(val)
		}
	}
	if v, ok := tmp["strategy"]; ok {
		cfg.Strategy = fmt.Sprintf("%v", v)
	}
	if v, ok := tmp["clientMaxRate"]; ok {
		switch val := v.(type) {
		case int:
			cfg.ClientMaxRate = val
		case float64:
			cfg.ClientMaxRate = int(val)
		}
	}
	if v, ok := tmp["key"]; ok {
		cfg.Key = fmt.Sprintf("%v", v)
	}
	return cfg
}
