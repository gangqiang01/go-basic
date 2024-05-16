package config

import (
	"k8s.io/klog/v2"
)

type ISyncorConfig struct {
	ServerAddr string
	CertFile   string
}

func GetISyncorConfig() *ISyncorConfig {
	cfg := &ISyncorConfig{}

	cfg.ServerAddr = ITHINGS_CONFIG.GetString("transport.isyncor.server_address")
	if cfg.ServerAddr == "" {
		klog.Warningf("isyncor.server_address is empty, we use the default 127.0.0.1:8082")
		cfg.ServerAddr = "127.0.0.1:8082"
	}

	cfg.CertFile = ITHINGS_CONFIG.GetString("transport.isyncor.cert_file")

	return cfg
}
