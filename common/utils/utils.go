package utils

import (
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

func NewUUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

// return ms
func GetNowTimeStamp() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}

/*
* GetLocalMACs
* get the local host's macaddress.
 */
func GetLocalMACs() []string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil
	}

	macAddrs := make([]string, 0)
	for _, ni := range netInterfaces {
		if !strings.HasPrefix(ni.Name, "e") &&
			!strings.HasPrefix(ni.Name, "w") &&
			!strings.HasPrefix(ni.Name, "p") {
			continue
		}
		macAddr := ni.HardwareAddr.String()
		if len(macAddr) > 6 {
			macAddrs = append(macAddrs, strings.ToUpper(macAddr))
		}
	}

	return macAddrs
}
