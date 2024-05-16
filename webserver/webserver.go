package webserver

import (
	"crypto/tls"
	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/webserver/router"
	"github.com/jwzl/beehive/pkg/core"
	"k8s.io/klog"
	"net/http"
	"time"
)

const (
	WebServerName = "webserver"
)

type WebServer struct {
}

// Register this module.
func Register() {
	ws := &WebServer{}
	core.Register(ws)
}

// Name
func (ws *WebServer) Name() string {
	return WebServerName
}

// Group
func (ws *WebServer) Group() string {
	return WebServerName
}

// Enable indicates whether this module is enabled
func (ws *WebServer) Enable() bool {
	//The module is always enabled!
	return true
}

// Start this module.
func (ws *WebServer) Start() {
	var err error

	initRouter := router.InitRouter()
	cfg := config.GetWebServerConfig()
	klog.Infof("Start web server on %s ", cfg.BindAddress)
	s := &http.Server{
		Addr:           cfg.BindAddress,
		Handler:        initRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if cfg.SSL {
		s.TLSConfig = createServerTLSConfiguration()
		err = s.ListenAndServeTLS(cfg.SSLCert, cfg.SSLKey)
	} else {
		err = s.ListenAndServe()
	}

	if err != nil {
		klog.Errorf("Start web server with error: %v", err)
		return
	}
}

// createServerTLSConfiguration creates a basic tls.Config to be used by servers with recommended TLS settings
func createServerTLSConfiguration() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
}
