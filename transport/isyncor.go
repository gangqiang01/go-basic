package transport

import (
	"context"
	"strings"
	"time"

	"encoding/json"
	"errors"
	"fmt"

	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	v1 "github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/core"
	"github.com/edgehook/ithings/transport/isync"
	"k8s.io/klog/v2"
)

// ISyncor is a inner syncor for
// mqtt connection status of apphub agent
type ISyncor struct {
	ctx    context.Context
	client *isync.ISyncGRPCClient
}

func NewISyncor(ctx context.Context) *ISyncor {
	//get local config.
	cfg := config.GetISyncorConfig()

retry_to_new:
	c, err := isync.NewISyncGRPCClient(ctx, cfg.ServerAddr, cfg.CertFile)
	if err != nil {
		select {
		case <-ctx.Done():
			klog.Infof("exit NewISyncor since the context canceled....")
			return nil
		default:
			klog.Errorf("NewISyncGRPCClient with err %s", err.Error())
			time.Sleep(3 * time.Second)
			goto retry_to_new
		}
	}

	return &ISyncor{
		ctx:    ctx,
		client: c,
	}
}

func (i *ISyncor) Run() {
retry_to_get:
	agentStatus, err := i.client.GetAgentStatus()
	if err != nil {
		select {
		case <-i.ctx.Done():
			klog.Infof("exit ISyncor Run since the context canceled....")
			return
		default:
			klog.Warningf("GetAgentStatus failed with err %s", err.Error())
			time.Sleep(3 * time.Second)
			goto retry_to_get
		}
	}

	klog.Infof("agentStatus = %v", agentStatus)
	StoreAgentStatus(agentStatus)

	//get agent status sync info
	for {
		select {
		case <-i.ctx.Done():
			klog.Infof("exit ISyncor Run since the context canceled....")
			return
		default:
			if err := i.client.DoSyncAgentStatus(func(agentStatus *isync.ReportAgentState) (string, string) {
				klog.Infof("XXX updated agentStatus = %v", agentStatus)
				StoreAgentStatus([]*isync.ReportAgentState{agentStatus})
				return "200", "success"
			}); err != nil {
				klog.Errorf("DoSyncAgentStatus failed with err %s", err.Error())

				//make all edge offline.
				for edgeID, edgeInfo := range v1.AgentStatusCache {
					if edgeInfo.Status != global.DeviceStatusOffline {
						//edgeInfo.Status = global.DeviceStatusOffline
						v1.AgentStatusCache[edgeID] = &v1.GrpcSyncStatus{Status: global.DeviceStatusOffline, AgentName: edgeInfo.AgentName}
					}
				}

				goto retry_to_get
			}
		}
	}
}

func (i *ISyncor) ReceiveAppHubMsg() {
	for {
		select {
		case <-i.ctx.Done():
			klog.Infof("exit ISyncor Run since the context canceled....")
			return
		default:
			if err := i.client.DoSendMsgToIthings(func(request *isync.AppHubRequest) (string, string) {
				//klog.Infof("receive apphub msg = %v", request)
				handleAppHubMsg(request)
				return "200", "success"
			}); err != nil {
				klog.Errorf("DoSendMsgToIthings failed with err %s", err.Error())
			}
		}
	}
}

func StoreAgentStatus(agentStatus []*isync.ReportAgentState) {
	if agentStatus == nil {
		return
	}

	for _, status := range agentStatus {
		if status == nil {
			continue
		}
		if status.AgentId == "" {
			continue
		}

		edgeInfo, ok := v1.AgentStatusCache[status.AgentId]
		if strings.Contains(strings.ToLower(status.Status), global.DeviceStatusOnline) {
			/*
			* edge is from offline to online, we consider some case as below:
			* 1. edge is startup or restart.
			* 2. the connection (edge <-> server) is lost and retry connect is successful.
			*
			* we will reset all mapper feth state. and let agent to re-fetch or server do
			* all update for mapper to sync device infomation correctlly.
			 */
			if !ok || edgeInfo.Status != global.DeviceStatusOnline || edgeInfo.AgentName != status.AgentName {
				core.ResetAllMapperFetchStateInThisEdge(status.AgentId)
				v1.AgentStatusCache[status.AgentId] = &v1.GrpcSyncStatus{Status: global.DeviceStatusOnline, AgentName: status.AgentName}
			}
		} else if strings.Contains(strings.ToLower(status.Status), global.DeviceStatusOffline) {
			/*
			* edge is from online to offline, we consider some case as below:
			* 1. edge is stopped, crashed, or stop by apphub server.
			* 2. the connection (edge <-> server) is lost.
			*
			* we should offline all device.
			 */
			if !ok || edgeInfo.Status != global.DeviceStatusOffline {
				core.MakeAllDeviceOfflineInThisEdge(status.AgentId)
				v1.AgentStatusCache[status.AgentId] = &v1.GrpcSyncStatus{Status: global.DeviceStatusOffline, AgentName: status.AgentName}
			}
		}

	}
}

func handleAppHubMsg(request *isync.AppHubRequest) error {
	tag := request.GetTag()
	switch tag {
	case v1.AlertTag:
		body := request.GetBody()
		msg := &v1.GrpcAlertInfo{}

		//decode data.
		err := json.Unmarshal([]byte(body), msg)
		if err != nil {
			klog.Errorf("Unmarshal grpc msg failed with err %s", err.Error())
			return err
		}

		//add alert list
		status := v1.AlertLogUnsolved
		if msg.Level == 0 {
			status = v1.AlertLogInvalid
		}
		processMonitorAlert := "Process"
		dockerMonitorAlert := "Docker"
		usbMonitorAlert := "USB"
		appMonitorAlert := "App"

		if msg.Name == processMonitorAlert {
			processList := strings.Split(msg.Description, " ")
			processes := processList[1 : len(processList)-2]
			processName := strings.Join(processes, " ")
			msg.Name = fmt.Sprintf("%s %s", msg.Name, processName)
		} else if msg.Name == dockerMonitorAlert {
			dockerList := strings.Split(msg.Description, " ")
			dockers := dockerList[0 : len(dockerList)-2]
			dockerName := strings.Join(dockers, " ")
			msg.Name = fmt.Sprintf("%s %s", msg.Name, dockerName)
		} else if msg.Name == appMonitorAlert {
			appList := strings.Split(msg.Description, " ")
			apps := appList[0 : len(appList)-2]
			appName := strings.Join(apps, " ")
			msg.Name = fmt.Sprintf("%s %s", msg.Name, appName)
		} else if msg.Name == usbMonitorAlert {
			index := strings.LastIndex(msg.Description, " ")
			msg.Name = fmt.Sprintf("%s %s", usbMonitorAlert, msg.Description[:index])
		}

		if isExist := model.IsExistAlertLogByNameAndDevice(msg.Name, msg.EdgeId, ""); isExist {
			alertLog, err := model.GetAlertLogByNameAndDevice(msg.Name, msg.EdgeId, "")
			if err != nil {
				klog.Errorf("GetAlertLogByNameAndDevice failed with err %s", err.Error())
				return err
			}
			if msg.Level == 0 && (alertLog.Level == 1 || alertLog.Level == 2) {
				status = v1.AlertLogResolved
			}

			if err := model.SaveAlertLog(alertLog.ID, msg.EdgeName, "", &status, &msg.Level, msg.Description); err != nil {
				klog.Errorf("SaveAlertLog failed with err %s", err.Error())
				return err
			}
			//rulelinkage.UpdateDeviceInstanceHealth(alertLog.DeviceId)
		} else {
			if err := model.AddAlertLog(&model.AlertLog{
				Name:        msg.Name,
				Description: msg.Description,
				Level:       msg.Level,
				EdgeName:    msg.EdgeName,
				EdgeId:      msg.EdgeId,
				DeviceName:  "",
				DeviceId:    "",
				Status:      status,
				LogType:     v1.AlertLogTypeMonitor,
			}); err != nil {
				klog.Errorf("AddAlertHistory failed with err %s", err.Error())
				return err
			}
		}

		//add alert history
		if err := model.AddAlertHistory(&model.AlertHistory{
			Name:        msg.Name,
			Description: msg.Description,
			Level:       msg.Level,
			EdgeName:    msg.EdgeName,
			EdgeId:      msg.EdgeId,
			DeviceName:  "",
			DeviceId:    "",
		}); err != nil {
			klog.Errorf("AddAlertHistory failed with err %s", err.Error())
			return err
		}
	case v1.DeviceOnPeripheral:
		body := request.GetBody()
		msg := &v1.GrpcDeviceOnPeripheralMsg{}

		//decode data.
		if err := json.Unmarshal([]byte(body), msg); err != nil {
			klog.Errorf("Unmarshal grpc msg failed with err %s", err.Error())
			return err
		}

		deviceInstances, err := model.GetAllDeviceInstancesV2ByEdgeId(msg.EdgeId)
		if err != nil {
			klog.Errorf("Get device instance err; deviceId: %s", msg.EdgeId)
			return errors.New(fmt.Sprintf("Get device instance err; deviceId: %s", msg.EdgeId))
		}

		for _, peripheral := range msg.Content {
			klog.Infof("peripheral: %v", peripheral)
			var (
				isStart interface{}
				normal  interface{}
			)
			eventName := ""
			isSet := true
			value := ""
			if value, ok := peripheral[v1.DeviceOnPrinterStatus]; ok {
				isStart = value
				normal = peripheral["normal"]
				eventName = v1.DeviceOnPrinterStatus
			}
			if value, ok := peripheral[v1.DeviceOnPrinterPaper]; ok {
				isStart = value
				normal = peripheral["normal"]
				eventName = v1.DeviceOnPrinterPaper
			}
			if value, ok := peripheral[v1.DeviceOnPrinterInk]; ok {
				isStart = value
				normal = peripheral["normal"]
				eventName = v1.DeviceOnPrinterInk
			}
			if value, ok := peripheral[v1.DeviceOnScreen]; ok {
				isStart = value
				normal = peripheral["normal"]
				eventName = v1.DeviceOnScreen
			}

			switch isStart.(type) {
			case bool:
				isSet = isStart.(bool)
			default:
				klog.Errorf("Peripheral isSet type error")
				continue
			}
			value = utils.ToString(normal)
			for _, deviceInstance := range deviceInstances {
				serviceInstances, err := model.GetServiceInstanceByDeviceID(deviceInstance.DeviceID)
				if err != nil {
					klog.Errorf("Get service instance err; deviceId: %s", deviceInstance.DeviceID)
					continue
				}
				for _, serviceInstance := range serviceInstances {
					eventInstance, err := model.GetEventInstanceByServiceIdAndName(serviceInstance.ID, eventName)
					if err != nil {
						klog.Errorf("Get event instance err;  name: %s", eventName)
						continue
					}
					if eventInstance.Name == eventName {
						//accessConfig := eventInstance.AccessConfig
						eventAccessConfig := &v1.EventsAccessConfig{}
						klog.Infof("before update accessConfig: %s", eventInstance.AccessConfig)
						//decode data.
						if err := json.Unmarshal([]byte(eventInstance.AccessConfig), eventAccessConfig); err != nil {
							klog.Errorf("Unmarshal event accessConfig failed with err %s", err.Error())
							return err
						}
						notEqual := int(3)
						eventAccessConfig.Rules[0].Value = value
						eventAccessConfig.Rules[0].Relation = &notEqual
						accessConfig, err := json.Marshal(eventAccessConfig)
						klog.Infof("After update accessconfig: %s", accessConfig)
						if err != nil {
							klog.Errorf("Marshal event accessConfig failed with err %s", err.Error())
							return err
						}

						if err := model.UpdateEventAccessConfig(serviceInstance.ID, eventName, string(accessConfig)); err != nil {
							klog.Errorf("Update event accessConfig err: %s", eventName)
							continue
						}
						if isSet {
							if deviceInstance.DeviceStatus == global.DeviceStatusInactive {
								if err := core.DoDeviceAction(deviceInstance.DeviceID, global.DeviceCreate); err != nil {
									klog.Errorf(fmt.Sprintf("%s: error: %s", deviceInstance.Name, err.Error()))
									continue
								}
							} else {
								if err := core.DoDeviceAction(deviceInstance.DeviceID, global.DeviceUpdate); err != nil {
									klog.Errorf(fmt.Sprintf("%s: error: %s", deviceInstance.Name, err.Error()))
									continue
								}
							}

						} else {
							if deviceInstance.DeviceStatus == global.DeviceStatusOnline {
								if err := core.DoDeviceAction(deviceInstance.DeviceID, global.DeviceStop); err != nil {
									klog.Errorf(fmt.Sprintf("%s: error: %s", deviceInstance.Name, err.Error()))
									continue
								}
							}
						}
					}

				}

			}
		}
	case v1.StopDeviceOnPeripheral:
		body := request.GetBody()
		msg := &v1.GrpcDeviceOnPeripheralMsg{}

		//decode data.
		if err := json.Unmarshal([]byte(body), msg); err != nil {
			klog.Errorf("Unmarshal grpc msg failed with err %s", err.Error())
			return err
		}

		deviceInstances, err := model.GetAllDeviceInstancesV2ByEdgeId(msg.EdgeId)
		if err != nil {
			klog.Errorf("Get device instance err; deviceId: %s", msg.EdgeId)
			return errors.New(fmt.Sprintf("Get device instance err; deviceId: %s", msg.EdgeId))
		}

		for _, deviceInstance := range deviceInstances {
			serviceInstances, err := model.GetServiceInstanceByDeviceID(deviceInstance.DeviceID)
			if err != nil {
				klog.Errorf("Get service instance err; deviceId: %s", deviceInstance.DeviceID)
				continue
			}
			for _, serviceInstance := range serviceInstances {
				eventInstances, err := model.GetEventInstanceByServiceId(serviceInstance.ID)
				if err != nil {
					klog.Errorf("Get event instance err;")
					continue
				}
				for _, eventInstance := range eventInstances {
					eventAccessConfig := &v1.EventsAccessConfig{}
					klog.Infof("before update accessConfig: %s", eventInstance.AccessConfig)
					//decode data.
					if err := json.Unmarshal([]byte(eventInstance.AccessConfig), eventAccessConfig); err != nil {
						klog.Errorf("Unmarshal event accessConfig failed with err %s", err.Error())
						return err
					}
					equal := int(2)
					eventAccessConfig.Rules[0].Value = "@@@"
					eventAccessConfig.Rules[0].Relation = &equal

					accessConfig, err := json.Marshal(eventAccessConfig)
					klog.Infof("After update accessconfig: %s", accessConfig)
					if err != nil {
						klog.Errorf("Marshal event accessConfig failed with err %s", err.Error())
						return err
					}

					if err := model.UpdateEventAccessConfig(serviceInstance.ID, eventInstance.Name, string(accessConfig)); err != nil {
						klog.Errorf("Update event accessConfig err: %s", eventInstance.Name)
						continue
					}

					if deviceInstance.DeviceStatus == global.DeviceStatusInactive {
						if err := core.DoDeviceAction(deviceInstance.DeviceID, global.DeviceCreate); err != nil {
							klog.Errorf(fmt.Sprintf("%s: error: %s", deviceInstance.Name, err.Error()))
							continue
						}
					} else {
						if err := core.DoDeviceAction(deviceInstance.DeviceID, global.DeviceUpdate); err != nil {
							klog.Errorf(fmt.Sprintf("%s: error: %s", deviceInstance.Name, err.Error()))
							continue
						}
					}
				}
			}
		}

	case v1.CancelDeviceOnPeripheral:
		body := request.GetBody()
		msg := &v1.GrpcDeviceOnPeripheralMsg{}

		//decode data.
		if err := json.Unmarshal([]byte(body), msg); err != nil {
			klog.Errorf("Unmarshal grpc msg failed with err %s", err.Error())
			return err
		}
		deviceInstances, err := model.GetAllDeviceInstancesV2ByEdgeId(msg.EdgeId)
		if err != nil {
			klog.Errorf("Get device instance err; deviceId: %s", msg.EdgeId)
			return errors.New(fmt.Sprintf("Get device instance err; deviceId: %s", msg.EdgeId))
		}
		for _, deviceInstance := range deviceInstances {
			if deviceInstance.DeviceStatus == global.DeviceStatusOnline {
				if err := core.DoDeviceAction(deviceInstance.DeviceID, global.DeviceStop); err != nil {
					klog.Errorf(fmt.Sprintf("%s: error: %s", deviceInstance.Name, err.Error()))
					continue
				}
			}
		}

	case v1.DeleteEdge:
		body := request.GetBody()
		msg := &v1.EdgeInfo{}

		//decode data.
		err := json.Unmarshal([]byte(body), msg)
		if err != nil {
			klog.Errorf("Unmarshal grpc msg failed with err %s", err.Error())
			return err
		}

		deviceInstances, err := model.GetDeviceInstanceByEdgeId(msg.EdgeId)
		if err != nil {
			klog.Errorf("GetDeviceInstanceByEdgeId with err %s", err.Error())
			return err
		}
		for _, deviceInstance := range deviceInstances {
			if err := model.DeleteDeviceInstance(deviceInstance.DeviceID); err != nil {
				klog.Errorf("DeleteDeviceInstance with err %s", err.Error())
				continue
			}
		}

		model.DeleteAlertHistoryByEdgeId(msg.EdgeId)
		model.DeleteAlertLogByEdgeId(msg.EdgeId)
	case v1.EdgeOffline:
		body := request.GetBody()
		msg := &v1.EdgeInfo{}

		//decode data.
		err := json.Unmarshal([]byte(body), msg)
		if err != nil {
			klog.Errorf("Unmarshal grpc msg failed with err %s", err.Error())
			return err
		}
		//init alertLog
		model.InitAlertLogByEdgeId(msg.EdgeId, v1.AlertLogInvalid)

	}
	return nil
}
