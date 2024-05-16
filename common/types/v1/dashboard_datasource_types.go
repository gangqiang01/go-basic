package v1

import (
	"encoding/json"
	"github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/transport/isync"
	"k8s.io/klog"
	"strconv"
)

const (
	DashboardQueryTypeTimeseries string = "timeserie"
	DashboardQueryTypeTable      string = "table"
)

type Dashboard_search_request struct {
	DataType string `form:"dataType" json:"dataType"`
	Type     string `form:"type" json:"type"`
	Device   string `form:"device" json:"device"`
	Ithing   string `form:"ithing" json:"ithing"`
	Module   string `form:"module" json:"module"`
	Property string `form:"property" json:"property"`
}
type Range struct {
	From string `from:"from" json:"from"`
	To   string `from:"to" json:"to"`
}

type Dashboard_query_request struct {
	Range         *Range    `from:"range" json:"range"`
	Targets       []*Target `form:"targets" json:"targets"`
	MaxDataPoints int64     `form:"maxDataPoints" json:"maxDataPoints"`
}

type Target struct {
	DataType    string `form:"dataType" json:"dataType"`
	Scene       string `form:"scene" json:"scene"`
	DeviceId    string `form:"deviceId" json:"deviceId"`
	Ithing      string `form:"ithing" json:"ithing"`
	Module      string `form:"module" json:"module"`
	Property    string `form:"property" json:"property"`
	Target      string `form:"target" json:"target"`
	DisplayName string `form:"displayName" json:"displayName"`
	RefId       string `form:"refId" json:"refId"`
	Type        string `form:"type" json:"type"`
}

type Dashboard_search_responce struct {
	Text  string `form:"text" json:"text,omitempty"`
	Value string `form:"value" json:"value,omitempty"`
}
type Dashboard_query_timeserie_responce struct {
	Target     string          `form:"target" json:"target,omitempty"`
	Datapoints [][]interface{} `form:"datapoints" json:"datapoints,omitempty"`
}

type Column struct {
	Text string `form:"text" json:"text"`
	Type string `form:"type" json:"type"`
}

type Dashboard_query_table_responce struct {
	Columns []*Column       `form:"columns" json:"columns"`
	Rows    [][]interface{} `form:"rows" json:"rows"`
	Type    string          `form:"type" json:"type"`
}

//scene

var (
	//overview series
	Dashboard_online_device         string = "Online Device"
	Dashboard_offline_device        string = "Offline Device"
	Dashboard_android_device        string = "Android Device"
	Dashboard_linux_device          string = "Linux Device"
	Dashboard_windows_device        string = "Windows Device"
	Dashboard_online_subDevice      string = "Online Sub Device"
	Dashboard_offline_subDevice     string = "Offline Sub Device"
	Dashboard_normal_device         string = "Normal Device"
	Dashboard_error_device          string = "Error Device"
	Dashboard_warning_device        string = "Warning Device"
	Dashboard_normal_subDevice      string = "Normal Sub Device"
	Dashboard_error_subDevice       string = "Error Sub Device"
	Dashboard_warning_subDevice     string = "Warning Sub Device"
	Dashboard_error_alert           string = "Error Alert"
	Dashboard_error_alert_handled   string = "Error Handled Alarm"
	Dashboard_warning_alert         string = "Warning Alert"
	Dashboard_warning_alert_handled string = "Warning Handled Alarm"
	//overview table
	Dashboard_error_alert_unhandled   string = "Unhandled Error Alarm"
	Dashboard_warning_alert_unhandled string = "Unhandled Warning Alarm"

	//information series
	Dashboard_total_memory_device      string = "Device Total Memory"
	Dashboard_free_memory_device       string = "Device Free Memory"
	Dashboard_usage_cpu_device         string = "Device CPU Usage"
	Dashboard_temp_cpu_device          string = "Device CPU temperature"
	Dashboard_total_storage_device     string = "Device Total Storage"
	Dashboard_free_storage_device      string = "Device Free Storage"
	Dashboard_battery_available_device string = "Device Battery Available"

	//information table
	Dashboard_monitor_app_device     string = "Device Monitor App"
	Dashboard_usb_device             string = "Device USB"
	Dashboard_monitor_process_device string = "Device Monitor Process"
	Dashboard_monitor_docker_device  string = "Device Monitor Docker"
)

var Dashboard_overview_scene_timeserie_resp = []*Dashboard_search_responce{
	//timeseries
	{Value: Dashboard_online_device, Text: Dashboard_online_device},
	{Value: Dashboard_offline_device, Text: Dashboard_offline_device},
	{Value: Dashboard_android_device, Text: Dashboard_android_device},
	{Value: Dashboard_linux_device, Text: Dashboard_linux_device},
	{Value: Dashboard_windows_device, Text: Dashboard_windows_device},

	{Value: Dashboard_online_subDevice, Text: Dashboard_online_subDevice},
	{Value: Dashboard_offline_subDevice, Text: Dashboard_offline_subDevice},
	{Value: Dashboard_normal_device, Text: Dashboard_normal_device},
	{Value: Dashboard_error_device, Text: Dashboard_error_device},
	{Value: Dashboard_warning_device, Text: Dashboard_warning_device},
	{Value: Dashboard_normal_subDevice, Text: Dashboard_normal_subDevice},
	{Value: Dashboard_error_subDevice, Text: Dashboard_error_subDevice},
	{Value: Dashboard_warning_subDevice, Text: Dashboard_warning_subDevice},
	{Value: Dashboard_error_alert, Text: Dashboard_error_alert},
	{Value: Dashboard_error_alert_handled, Text: Dashboard_error_alert_handled},
	{Value: Dashboard_warning_alert, Text: Dashboard_warning_alert},
	{Value: Dashboard_warning_alert_handled, Text: Dashboard_warning_alert_handled},
}
var Dashboard_overview_scene_table_resp = []*Dashboard_search_responce{
	//table
	{Value: Dashboard_error_alert_unhandled, Text: Dashboard_error_alert_unhandled},
	{Value: Dashboard_warning_alert_unhandled, Text: Dashboard_warning_alert_unhandled},
}

var Dashboard_information_scene_timeserie_resp = []*Dashboard_search_responce{
	//timeseries
	{Value: Dashboard_total_memory_device, Text: Dashboard_total_memory_device},
	{Value: Dashboard_free_memory_device, Text: Dashboard_free_memory_device},
	{Value: Dashboard_usage_cpu_device, Text: Dashboard_usage_cpu_device},
	{Value: Dashboard_temp_cpu_device, Text: Dashboard_temp_cpu_device},
	{Value: Dashboard_total_storage_device, Text: Dashboard_total_storage_device},
	{Value: Dashboard_free_storage_device, Text: Dashboard_free_storage_device},
	{Value: Dashboard_battery_available_device, Text: Dashboard_battery_available_device},
}

var Dashboard_information_scene_table_resp = []*Dashboard_search_responce{
	//table
	{Value: Dashboard_monitor_app_device, Text: Dashboard_monitor_app_device},
	{Value: Dashboard_usb_device, Text: Dashboard_usb_device},
	{Value: Dashboard_monitor_process_device, Text: Dashboard_monitor_process_device},
	{Value: Dashboard_monitor_docker_device, Text: Dashboard_monitor_docker_device},
}

func HandleSceneByTimeSeries(ctype, edgeId string) []interface{} {
	var value interface{}
	errorLevel := AlertErrorLevel
	warningLevel := AlertWarningLevel
	handledStatus := AlertLogResolved
	invalidStatus := AlertLogInvalid
	switch ctype {
	case Dashboard_online_device, Dashboard_offline_device,
		Dashboard_normal_device, Dashboard_error_device, Dashboard_warning_device,
		Dashboard_android_device, Dashboard_linux_device, Dashboard_windows_device:
		appHubResponse, err := isync.SendMsgToAppHub("overview", "")
		if err != nil {
			klog.Errorln("Get Agent info by grpc error: err", err)
			return nil
		}
		if appHubResponse.StatusCode == "200" {
			msg := appHubResponse.Msg
			//klog.Infof("agentInfo: %s", agentInfo)
			var overview *Dashboard_grpc_overview_responce
			err := json.Unmarshal([]byte(msg), &overview)
			if err != nil {
				klog.Errorln("Unmarshal error: %s", err.Error())
				return nil
			}
			switch ctype {
			case Dashboard_online_device:
				value = overview.Online
			case Dashboard_offline_device:
				value = overview.Offline
			case Dashboard_android_device:
				value = overview.Android
			case Dashboard_linux_device:
				value = overview.Linux
			case Dashboard_windows_device:
				value = overview.Windows
			case Dashboard_error_device:
				value = overview.Error
			case Dashboard_warning_device:
				value = overview.Warning
			case Dashboard_normal_device:
				value = overview.Normal
			}
		} else {
			klog.Errorln("Get Agent info by grpc error: err", appHubResponse.Msg)
			return nil
		}
	case Dashboard_online_subDevice:
		count, _ := model.GetDeviceInstanceCountByStatus("online")
		value = int(count)
	case Dashboard_offline_subDevice:
		count, _ := model.GetDeviceInstanceCountByStatus("offline")
		value = int(count)
	case Dashboard_normal_subDevice:
		count, _ := model.GetDeviceInstanceCountByStatusAndHealth("online", int64(0))
		value = int(count)
	case Dashboard_error_subDevice:
		count, _ := model.GetDeviceInstanceCountByStatusAndHealth("online", int64(2))
		value = int(count)
	case Dashboard_warning_subDevice:
		count, _ := model.GetDeviceInstanceCountByStatusAndHealth("online", int64(1))
		value = int(count)
	case Dashboard_error_alert:
		errorTotalCount, _ := model.GetAlertLogCountByCondition("", "", nil, &errorLevel, nil, nil, "")
		value = int(errorTotalCount)
	case Dashboard_error_alert_handled:
		errorHandledCount, _ := model.GetAlertLogCountByCondition("", "", &handledStatus, &errorLevel, nil, nil, "")
		errorInvalidCount, _ := model.GetAlertLogCountByCondition("", "", &invalidStatus, &errorLevel, nil, nil, "")
		value = int(errorHandledCount + errorInvalidCount)
	case Dashboard_warning_alert:
		warningTotalCount, _ := model.GetAlertLogCountByCondition("", "", nil, &warningLevel, nil, nil, "")
		value = int(warningTotalCount)
	case Dashboard_warning_alert_handled:
		warningHandledCount, _ := model.GetAlertLogCountByCondition("", "", &handledStatus, &warningLevel, nil, nil, "")
		warningInvalidCount, _ := model.GetAlertLogCountByCondition("", "", &invalidStatus, &warningLevel, nil, nil, "")
		value = int(warningHandledCount + warningInvalidCount)
	case Dashboard_total_memory_device, Dashboard_free_memory_device,
		Dashboard_usage_cpu_device, Dashboard_temp_cpu_device,
		Dashboard_total_storage_device, Dashboard_free_storage_device,
		Dashboard_battery_available_device:
		req := &Dashboard_grpc_information_request{
			EdgeId: edgeId,
			Type:   "hardware",
		}
		msg, _ := json.Marshal(req)
		appHubResponse, err := isync.SendMsgToAppHub(Dashboard_information, string(msg))
		if err != nil {
			klog.Errorln("Get Agent info by grpc error: err", err)
			return nil
		}
		if appHubResponse.StatusCode == "200" {
			info := appHubResponse.Msg
			var hardware *Dashboard_grpc_hardware
			err := json.Unmarshal([]byte(info), &hardware)
			if err != nil {
				klog.Errorln("Unmarshal error: %s", err.Error())
				return nil
			}
			switch ctype {
			case Dashboard_total_memory_device:
				val := hardware.Memory.Total
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}
			case Dashboard_free_memory_device:
				val := hardware.Memory.Free
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}

			case Dashboard_usage_cpu_device:
				val := hardware.Cpu.Usage
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}
			case Dashboard_temp_cpu_device:
				val := hardware.Cpu.Temp
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}
			case Dashboard_total_storage_device:
				val := hardware.storage.Total
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}
			case Dashboard_free_storage_device:
				val := hardware.storage.Free
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}

			case Dashboard_battery_available_device:
				val := hardware.Battery.Available
				value, err = strconv.ParseFloat(val, 32)
				if err != nil {
					klog.Errorln("Parse float error: %s", err.Error())
					return nil
				}
			}
		}

	default:
		return nil
	}
	dataPoint := []interface{}{
		value,
		utils.GetNowTimeStamp(),
	}
	return dataPoint
}

func HandleSceneByTable(ctype, edgeId string) *Dashboard_query_table_responce {
	rows := make([][]interface{}, 0)
	switch ctype {
	case Dashboard_error_alert_unhandled, Dashboard_warning_alert_unhandled:
		colomns := []*Column{
			{
				Text: "Time",
				Type: "time",
			},
			{
				Text: "Name",
				Type: "string",
			},
			{
				Text: "Device Name",
				Type: "string",
			},
			{
				Text: "Sub Device Name",
				Type: "string",
			},
			{
				Text: "Description",
				Type: "string",
			}}
		var (
			status = AlertLogUnsolved
			level  = AlertErrorLevel
		)
		if ctype == Dashboard_warning_alert_unhandled {
			level = AlertWarningLevel
		}
		alertLogs, err := model.GetAlertLogByCondition(
			"", "", &status, &level, nil, nil, "")
		if err != nil {
			return nil
		}
		for _, alertLog := range alertLogs {
			row := []interface{}{
				alertLog.CreateTimeStamp,
				alertLog.Name,
				alertLog.DeviceName,
				alertLog.EdgeName,
				alertLog.Description,
			}
			rows = append(rows, row)
		}

		table := &Dashboard_query_table_responce{
			Columns: colomns,
			Rows:    rows,
			Type:    DashboardQueryTypeTable,
		}
		return table
	case Dashboard_monitor_app_device:
		colomns := []*Column{
			{
				Text: "Status",
				Type: "string",
			},
			{
				Text: "Name",
				Type: "string",
			},
			{
				Text: "Package",
				Type: "string",
			},
			{
				Text: "Version",
				Type: "string",
			}}

		req := &Dashboard_grpc_information_request{
			EdgeId: edgeId,
			Type:   "appMonitor",
		}
		msg, _ := json.Marshal(req)
		appHubResponse, err := isync.SendMsgToAppHub(Dashboard_information, string(msg))
		if err != nil {
			klog.Errorln("Get monitor app info by grpc error: err", err)
			return nil
		}
		if appHubResponse.StatusCode == "200" {
			msg := appHubResponse.Msg
			//klog.Infof("agentInfo: %s", agentInfo)
			var apps []*Dashboard_grpc_monitor_app
			err := json.Unmarshal([]byte(msg), &apps)
			if err != nil {
				klog.Errorln("Unmarshal error: %s", err.Error())
				return nil
			}
			for _, app := range apps {
				row := []interface{}{
					app.Status,
					app.AppName,
					app.Package,
					app.Version,
				}
				rows = append(rows, row)
			}

			table := &Dashboard_query_table_responce{
				Columns: colomns,
				Rows:    rows,
				Type:    DashboardQueryTypeTable,
			}
			return table
		} else {
			klog.Errorln("Get monitor app info by grpc error: err", appHubResponse.Msg)
			return nil
		}
	case Dashboard_usb_device:
		colomns := []*Column{
			{
				Text: "Status",
				Type: "string",
			},
			{
				Text: "Usb",
				Type: "string",
			},
			{
				Text: "Manufacturer",
				Type: "string",
			},
			{
				Text: "Peripheral",
				Type: "string",
			}}
		req := &Dashboard_grpc_information_request{
			EdgeId: edgeId,
			Type:   "usb",
		}
		msg, _ := json.Marshal(req)
		appHubResponse, err := isync.SendMsgToAppHub(Dashboard_information, string(msg))
		if err != nil {
			klog.Errorln("Get msg by grpc error: err", err)
			return nil
		}
		if appHubResponse.StatusCode == "200" {
			msg := appHubResponse.Msg
			//klog.Infof("agentInfo: %s", agentInfo)
			var usbs []*Dashboard_grpc_monitor_usb
			err := json.Unmarshal([]byte(msg), &usbs)
			if err != nil {
				klog.Errorln("Unmarshal error: %s", err.Error())
				return nil
			}
			for _, usb := range usbs {
				row := []interface{}{
					usb.Status,
					usb.Usb,
					usb.Manufacturer,
					usb.Peripheral,
				}
				rows = append(rows, row)
			}

			table := &Dashboard_query_table_responce{
				Columns: colomns,
				Rows:    rows,
				Type:    DashboardQueryTypeTable,
			}
			return table
		} else {
			klog.Errorln("Get monitor usb by grpc error: err", appHubResponse.Msg)
			return nil
		}
	case Dashboard_monitor_process_device:
		colomns := []*Column{
			{
				Text: "Status",
				Type: "string",
			},
			{
				Text: "CMD",
				Type: "string",
			},
			{
				Text: "Process ID",
				Type: "string",
			},
			{
				Text: "CPU Loading",
				Type: "string",
			},
			{
				Text: "Memory Loading",
				Type: "string",
			}}
		req := &Dashboard_grpc_information_request{
			EdgeId: edgeId,
			Type:   "processMonitor",
		}
		msg, _ := json.Marshal(req)
		appHubResponse, err := isync.SendMsgToAppHub(Dashboard_information, string(msg))
		if err != nil {
			klog.Errorln("Get msg by grpc error: err", err)
			return nil
		}
		if appHubResponse.StatusCode == "200" {
			msg := appHubResponse.Msg
			var processes []*Dashboard_grpc_monitor_process
			err := json.Unmarshal([]byte(msg), &processes)
			if err != nil {
				klog.Errorln("Unmarshal error: %s", err.Error())
				return nil
			}
			for _, process := range processes {
				row := []interface{}{
					process.Status,
					process.Cmd,
					process.Id,
					process.CpuLoading,
					process.MemoryLoading,
				}
				rows = append(rows, row)
			}

			table := &Dashboard_query_table_responce{
				Columns: colomns,
				Rows:    rows,
				Type:    DashboardQueryTypeTable,
			}
			return table

		} else {
			klog.Errorln("Get monitor process by grpc error: err", appHubResponse.Msg)
			return nil
		}
	case Dashboard_monitor_docker_device:
		colomns := []*Column{
			{
				Text: "Status",
				Type: "string",
			},
			{
				Text: "Name",
				Type: "string",
			},
			{
				Text: "Image",
				Type: "string",
			},
			{
				Text: "Port",
				Type: "string",
			},
			{
				Text: "Created",
				Type: "string",
			}}
		req := &Dashboard_grpc_information_request{
			EdgeId: edgeId,
			Type:   "dockerMonitor",
		}
		msg, _ := json.Marshal(req)
		appHubResponse, err := isync.SendMsgToAppHub(Dashboard_information, string(msg))
		if err != nil {
			klog.Errorln("Get Agent info by grpc error: err", err)
			return nil
		}
		if appHubResponse.StatusCode == "200" {
			msg := appHubResponse.Msg
			var dockers []*Dashboard_grpc_monitor_docker
			err := json.Unmarshal([]byte(msg), &dockers)
			if err != nil {
				klog.Errorln("Unmarshal error: %s", err.Error())
				return nil
			}
			for _, docker := range dockers {
				row := []interface{}{
					docker.Status,
					docker.Name,
					docker.Image,
					docker.Ports,
					docker.Created,
				}
				rows = append(rows, row)
			}

			table := &Dashboard_query_table_responce{
				Columns: colomns,
				Rows:    rows,
				Type:    DashboardQueryTypeTable,
			}
			return table

		} else {
			klog.Errorln("Get Agent info by grpc error: err", appHubResponse.Msg)
			return nil
		}

	default:
		return nil
	}
	return nil
}

// data type
var (
	Dashboard_collection  string = "collection"
	Dashboard_overview    string = "overview"
	Dashboard_information string = "information"
)

// grpc overview
type Dashboard_grpc_overview_responce struct {
	//overview
	Android int `form:"android" json:"android"`
	Linux   int `form:"linux" json:"linux"`
	Windows int `form:"windows" json:"windows"`
	Error   int `form:"error" json:"error"`
	Warning int `form:"warning" json:"warning"`
	Normal  int `form:"normal" json:"normal"`
	Online  int `form:"online" json:"online"`
	Offline int `form:"offline" json:"offline"`
}

// grpc information
type Dashboard_grpc_information_request struct {
	EdgeId string `form:"edgeId" json:"edgeId"`
	Type   string `form:"type" json:"type"`
}
type Dashboard_grpc_information_responce struct {
	//overview
	Hardware       int `form:"hardware" json:"hardware"`
	AppMonitor     int `form:"linux" json:"linux"`
	Usb            int `form:"windows" json:"windows"`
	DockerMonitor  int `form:"error" json:"error"`
	ProcessMonitor int `form:"warning" json:"warning"`
	DeviceDetails  int `form:"normal" json:"normal"`
}

// {"battery":{"available":"N/A"},"cpu":{"temp":"27","usage":"44.6"},"memory":{"free":"56.3","total":"7.68","usage":"43.7"},"storage":{"free":59,"total":"74.93","usage":41}}
type Dashboard_grpc_battery struct {
	Available string `form:"available" json:"available"`
}

type Dashboard_grpc_cpu struct {
	Temp  string `form:"temp" json:"temp"`
	Usage string `form:"usage" json:"usage"`
}

type Dashboard_grpc_memory struct {
	Free  string `form:"free" json:"free"`
	Total string `form:"total" json:"total"`
}

type Dashboard_grpc_storage struct {
	Free  string `form:"free" json:"free"`
	Total string `form:"total" json:"total"`
}

type Dashboard_grpc_hardware struct {
	Battery Dashboard_grpc_battery `json:"battery"`
	Cpu     Dashboard_grpc_cpu     `json:"cpu"`
	Memory  Dashboard_grpc_memory  `json:"memory"`
	storage Dashboard_grpc_storage `json:"storage"`
}

// [{"usb":"1d6b:0003","status":"connect"},{"usb":"0eef:c000","status":"connect"},{"usb":"14e1:6000","status":"connect"},{"usb":"1d6b:0002","status":"connect"}]
type Dashboard_grpc_monitor_usb struct {
	Status       string `form:"status" json:"status"`
	Usb          string `form:"usb" json:"usb"`
	Manufacturer string `form:"manufacturer" json:"manufacturer"`
	Peripheral   string `form:"peripheral" json:"peripheral"`
}

// [{"threadcmd":"java","threadpid":"3413","threadcpuloading":"0.22","threadmemloading":"7.97","threadstatus":"S","threadusername":"root","threadismonitor":"","threadcmdline":"","threadexe":""}]
type Dashboard_grpc_monitor_process struct {
	Cmd           string `form:"threadcmd" json:"threadcmd"`
	Id            string `form:"threadpid" json:"threadpid"`
	CpuLoading    string `form:"threadcpuloading" json:"threadcpuloading"`
	MemoryLoading string `form:"threadmemloading" json:"threadmemloading"`
	Status        string `form:"threadstatus" json:"threadstatus"`
}

// [{"command":"docker-entrypoint.sh postgres -c max_connections=500 -c shared_buffers=1024MB -c work_mem=131072","id":"29a3cb3c7dcd54d42a51d831df7b473fe8e48af8f96e3d24754b94508d648d0d","name":"m2m-postgresSQL","ports":" 5432:5432 5432:5432","state":"running","compose":"singlerepo","created":"2022-06-08 17:45:52","image":"edgesolution/apphub-postgres:v1.0","ismonitor":"false","cpuusage":"0.00","memusage":"0.99"},{"command":"/usr/bin/docker-entrypoint.sh server /data","id":"f841f998697358ab85bc7fbb773a392868422819955c929c43a39ee08904d538","name":"minio","ports":" 9000:9000 9000:9000","state":"running","compose":"singlerepo","created":"2022-06-08 17:45:52","image":"edgesolution/apphub-minio:v1.0","ismonitor":"false","cpuusage":"0.00","memusage":"0.14"}]
type Dashboard_grpc_monitor_docker struct {
	Cmd     string `form:"command" json:"command"`
	Id      string `form:"id" json:"id"`
	Name    string `form:"name" json:"name"`
	Image   string `form:"image" json:"image"`
	Ports   string `form:"ports" json:"ports"`
	Created string `form:"created" json:"created"`
	Status  string `form:"state" json:"state"`
}

// [{\"appName\":\"Sample\",\"isrunning\":\"false\",\"packagename\":\"com.suke.widget.sample\",\"versionName\":\"1.0\"}]
type Dashboard_grpc_monitor_app struct {
	AppName string `form:"appName" json:"appName"`
	Status  string `form:"isrunning" json:"isrunning"`
	Package string `form:"packagename" json:"packagename"`
	Version string `form:"versionName" json:"versionName"`
}
