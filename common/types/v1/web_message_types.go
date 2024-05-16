package v1

import "github.com/edgehook/ithings/common/dbm/model"

type DeviceTwinMessage struct {
	DeviceID string `json:"DeviceID"`
	// DeviceTwin provides a logical representation of A series of control properties on a device
	// Optional:
	//It's the twin's history
	DeviceTwinsHistory []*DeviceTwin `json:"device_twins,omitempty"`
	//desired twins we want to set. and it's just used for updateDesiredTwins API.
	//Optional:
	DesiredTwins []*TwinProperty `json:"desired_twins,omitempty"`
}

func NewDeviceTwinMessage(deviceID string) *DeviceTwinMessage {
	return &DeviceTwinMessage{
		DeviceID:           deviceID,
		DeviceTwinsHistory: make([]*DeviceTwin, 0),
	}
}

func NewDeviceTwinMessageV2(deviceID string) *DeviceTwinMessage {
	return &DeviceTwinMessage{
		DeviceID:     deviceID,
		DesiredTwins: make([]*TwinProperty, 0),
	}
}

func (dtm *DeviceTwinMessage) AddTwins(desired, reported []*TwinProperty) {
	twin := NewDeviceTwin(desired, reported)
	dtm.DeviceTwinsHistory = append(dtm.DeviceTwinsHistory, twin)
}

type DesiredPropertyMsg struct {
	DeviceId     string      `form:"deviceId" json:"deviceId"`
	ServiceName  string      `form:"serviceName" json:"serviceName"`
	PropertyName string      `form:"propertyName" json:"propertyName"`
	Value        interface{} `form:"value" json:"value,omitempty"`
}

// Device instance access config web api
type InstanceConfig struct {
	// Required: The device property name.
	//device Instance, property instance, event instance, command instance
	Name         string `form:"name" json:"name"`
	DeviceId     string `form:"deviceId" json:"deviceId"`
	ServiceName  string `form:"serviceName" json:"serviceName"`
	AccessConfig string `form:"accessConfig" json:"accessConfig"`
	InstanceType string `form:"type" json:"type"`
}

// Alert web api
type Alert struct {
	Name         string `form:"name" json:"name" binding:"required"`
	Description  string `form:"description" json:"description"`
	Level        *int64 `form:"level" json:"level" binding:"required"`
	Notification string `form:"notification" json:"notification"`
}

// Alert log api
type AlertLog struct {
	Record string `form:"record" json:"record"`
	Status *int32 `form:"status" json:"status"`
	Level  *int64 `form:"level" json:"level"`
}

// RuleLinkage web api  status: enable/disable
type RuleLinkage struct {
	Name            string                `form:"name" json:"name"`
	Description     *string               `form:"description" json:"description"`
	DeviceModelName string                `form:"deviceModelName" json:"deviceModelName,omitempty"`
	Trigger         []*RuleLinkageTrigger `form:"trigger" json:"trigger"`
	Filter          string                `form:"filter" json:"filter"`
	Action          string                `form:"action" json:"action"`
	Status          string                `form:"status" json:"status"`
}

// DataForward web data
type DataForward struct {
	Name        string               `form:"name" json:"name"`
	Description *string              `form:"description" json:"description"`
	Source      []*DataForwardSource `form:"source" json:"source,omitempty"`
	Destination string               `form:"destination" json:"destination"`
	Status      string               `form:"status" json:"status"`
}

type DeviceInstanceData struct {
	*model.DeviceInstance `json:",inline,omitempty"`
	EdgeName              string `json:"edgeName"`
}

// fluxdb config
type InfluxDbConfig struct {
	Persistence *bool  `form:"persistence" json:"persistence" binding:"required"`
	Address     string `form:"address" json:"address" binding:"required"`
	Username    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	DbName      string `form:"dbName" json:"dbName" binding:"required"`
	Duration    string `form:"duration" json:"duration" binding:"required"`
}
