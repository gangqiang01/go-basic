package webserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edgehook/ithings/webserver/router"
	"github.com/jwzl/beehive/pkg/core"
	"k8s.io/klog"
)

const (
	WebServerName = "webserver"
)

var BindAddress = ":9001"

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
	if !strings.Contains(BindAddress, ":") {
		BindAddress = fmt.Sprintf(":%s", BindAddress)
	}
	klog.Infof("Start web server on %s ", BindAddress)
	s := &http.Server{
		Addr:           BindAddress,
		Handler:        initRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err = s.ListenAndServe(); err != nil {
		klog.Errorf("Start web server with error: %v", err)
		return
	}
}
