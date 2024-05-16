package v1

import "github.com/edgehook/ithings/common/dbm/model"

// grpc action rulelinkage
type GrpcDeviceControl struct {
	EdgeId string `form:"edgeId" json:"edgeId"`
	Action string `form:"action" json:"action"`
	Value  string `form:"value" json:"value"`
}

// grpc report rulelinkage
type GrpcReport struct {
	EdgeId       string `form:"edgeId" json:"edgeId,omitempty"`
	DeviceName   string `form:"deviceName" json:"deviceName,omitempty"`
	Level        *int64 `form:"level" json:"level,omitempty"`
	Type         string `form:"type" json:"type,omitempty"`
	Details      string `form:"details" json:"details,omitempty"`
	Notification string `form:"notification" json:"notification,omitempty"`
}

type GrpcSyncStatus struct {
	Status    string `json:"status,omitempty"`
	AgentName string `json:"agentName,omitempty"`
}

// isync get agent status
var AgentStatusCache = make(map[string]*GrpcSyncStatus)

// Adge topo data
type EdgeTopoData struct {
	Name     string      `form:"name" json:"name,omitempty"`
	Children []*EdgeInfo `form:"children" json:"children,omitempty"`
}

type EdgeInfo struct {
	Name     string                  `form:"name" json:"name,omitempty"`
	EdgeId   string                  `form:"edgeId" json:"edgeId,omitempty"`
	Status   string                  `form:"status" json:"status,omitempty"`
	Children []*model.DeviceInstance `form:"children" json:"children,omitempty"`
}

type TopoDataParams struct {
	Keywords string `json:"keywords,omitempty"`
}

// receive msg from AppHub
const (
	AlertTag                 string = "alert"
	DeviceOnPeripheral       string = "deviceOnPeripheral"
	StopDeviceOnPeripheral   string = "stopDeviceOnPeripheral"
	CancelDeviceOnPeripheral string = "cancelDeviceOnPeripheral"
	DeleteEdge               string = "deleteEdge"
	EdgeOffline              string = "edgeOffline"
	// operation record event
	SubDeviceType      string = "subDevice"
	DeviceModelType    string = "deviceModel"
	RuleLinkageType    string = "ruleLinkage"
	DataForwardType    string = "dataForward"
	AlertConfigType    string = "alertConfig"
	AlertHistoryType   string = "alertHistory"
	AlertType          string = "alert"
	SystemConfigType   string = "systemConfig"
	Create             string = "0"
	Update             string = "1"
	Delete             string = "2"
	Clear              string = "6"
	UpdateLevel        string = "7"
	UpdateNotification string = "8"
	Enable             string = "9"
	UpdateRecords      string = "24"
	SetInfluxDb        string = "11"
	ImportSubDevice    string = "20"
	CopySubDevice      string = "21"
	ImportDeviceModel  string = "22"
)

type GrpcAlertInfo struct {
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	Level       int64  `form:"level" json:"level"`
	EdgeId      string `form:"edgeId" json:"edgeId"`
	EdgeName    string `form:"edgeName" json:"edgeName"`
}

type OperationRecord struct {
	Event    string `form:"event" json:"event"`
	UserName string `form:"username" json:"username"`
	Target   string `form:"target" json:"target"`
	CType    string `form:"type" json:"type"`
}
