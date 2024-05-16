package core

import (
	"encoding/json"
	"strings"
	"time"

	db "github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types/v1"
)

func EdgeIsOnline(edgeID string) bool {
	if edgeID == "" {
		return false
	}

	val, ok := v1.AgentStatusCache[edgeID]
	if !ok {
		return false
	}

	if strings.Contains(val.Status, global.DeviceStatusOnline) {
		return true
	}

	return false
}

func ResetAllMapperFetchStateInThisEdge(edgeID string) {
	for defaultICore == nil {
		time.Sleep(time.Millisecond * 100)
	}

	dsm := defaultICore.iCore.devStatusMgr

	protoMap := db.GetAllProtocolTypeInThisEdge(edgeID)

	for proto, _ := range protoMap {
		dsm.UpdateFetchHistory(edgeID, proto, false)
	}
}

func MakeAllDeviceOfflineInThisEdge(edgeID string) {
	for defaultICore == nil {
		time.Sleep(time.Millisecond * 100)
	}

	dsm := defaultICore.iCore.devStatusMgr

	//make all device offline in this edge.
	dsm.UpdateAllDeviceStatusInThisEdge(edgeID, global.DeviceStatusOffline)
}

func DeviceIsInactive(deviceID string) bool {
	for defaultICore == nil {
		time.Sleep(time.Millisecond * 100)
	}

	dsm := defaultICore.iCore.devStatusMgr

	devInfo := dsm.GetDeviceStatus(deviceID)
	if devInfo == nil {
		return false
	}

	return devInfo.DeviceStatus == global.DeviceStatusInactive
}

func DeviceIsStarted(deviceID string) bool {
	for defaultICore == nil {
		time.Sleep(time.Millisecond * 100)
	}

	dsm := defaultICore.iCore.devStatusMgr

	devInfo := dsm.GetDeviceStatus(deviceID)
	if devInfo == nil {
		return false
	}

	return devInfo.State == global.DeviceStateStarted
}

func decodeReportDevicesMessage(content string) (*v1.ReportDevicesMessage, error) {
	msg := &v1.ReportDevicesMessage{}

	//decode data.
	err := json.Unmarshal([]byte(content), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func decodeReportDevicesStatusMessage(content string) (*v1.DevicesStatusMessage, error) {
	msg := &v1.DevicesStatusMessage{}

	//decode data.
	err := json.Unmarshal([]byte(content), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
