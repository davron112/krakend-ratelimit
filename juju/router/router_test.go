package router

import (
	"encoding/json"
	"testing"

	"github.com/davron112/lura/v2/config"
)

func TestConfigGetter(t *testing.T) {
	serializedCfg := []byte(`{
		"github.com/devopsfaith/krakend-ratelimit/juju/router": {
			"max_rate":10,
			"capacity":10
		}
	}`)
	var dat config.ExtraConfig
	if err := json.Unmarshal(serializedCfg, &dat); err != nil {
		t.Error(err.Error())
	}
	cfg, err := ConfigGetter(dat)
	if cfg.MaxRate != 10 {
		t.Errorf("wrong value for MaxRate. Want: 10, have: %f", cfg.MaxRate)
	}
	if cfg.ClientMaxRate != 0 {
		t.Errorf("wrong value for ClientMaxRate. Want: 0, have: %f", cfg.ClientMaxRate)
	}
	if cfg.Capacity != 10 {
		t.Errorf("wrong value for Capacity. Want: 10, have: %d", cfg.Capacity)
	}
	if cfg.ClientCapacity != 0 {
		t.Errorf("wrong value for ClientCapacity. Want: 0, have: %d", cfg.ClientCapacity)
	}
	if cfg.Strategy != "" {
		t.Errorf("wrong value for Strategy. Want: '', have: %s", cfg.Strategy)
	}
	if cfg.Key != "" {
		t.Errorf("wrong value for Key. Want: '', have: %s", cfg.Key)
	}
	if err != nil {
		t.Error(err)
	}
}
