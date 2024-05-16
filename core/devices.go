package core

import (
	"encoding/json"
	"fmt"
	"time"

	db "github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
	"github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

// API: GetDeviceSpecMeta
// get device spec meta and it's used for Fetch.
func GetDeviceSpecMeta(edgeID, protocType string) ([]*v1.DeviceSpecMeta, error) {
	// we will not send any message to edge since
	// edge is offline.
	if !EdgeIsOnline(edgeID) {
		return nil, global.ErrEdgeIsNotOnline
	}

	devices := make([]*v1.DeviceSpecMeta, 0)
	deviceInstances, err := db.GetAllDeviceInstancesV2(edgeID, protocType)
	if err != nil {
		return nil, err
	}

	for _, di := range deviceInstances {
		if di == nil {
			continue
		}

		dm, err := db.GetDeviceModelAllInfoByID(di.DeviceModelId)
		if err != nil || dm == nil {
			klog.Warningf("Device %s has no model, Skip it...", di.DeviceID)
			continue
		}

		devSpec := newDeviceSpecMeta(di, dm)
		if devSpec != nil && di.DeviceStatus != global.DeviceStatusInactive {
			devices = append(devices, devSpec)
		}
	}

	return devices, nil
}

// API: DoDeviceAction
// do device create/start/stop/delete/update.
func DoDeviceAction(deviceID string, action int) error {
	if defaultICore == nil {
		return global.ErrCoreNotReady
	}

	di, err := db.GetDeviceInstanceAllInfo(deviceID)
	if err != nil {
		return err
	}

	// we will not process any message from edge since
	// edge is offline.
	if !EdgeIsOnline(di.EdgeID) {
		return global.ErrEdgeIsNotOnline
	}

	dsm := defaultICore.iCore.devStatusMgr

	mutex := dsm.DeviceStatusLock(deviceID)
	defer dsm.DeviceStatusUnlock(mutex)

	devInfo := dsm.GetDeviceStatus(deviceID)
	if devInfo == nil {
		return global.ErrNoSuchDevice
	}

	switch action {
	case global.DeviceCreate:
		klog.Infof("create device (%s)", deviceID)
		//device has already active, we return directlly.
		if !devInfo.IsInactive() {
			return nil
		}

		dm, err := db.GetDeviceModelAllInfoByID(di.DeviceModelId)
		if err != nil {
			return err
		}

		if dm == nil {
			return global.ErrNoSuchDeviceModel
		}
		di.State = global.DeviceStateStarted
		devSpec := newDeviceSpecMeta(&di, dm)
		if devSpec == nil {
			return global.ErrNoSuchDevice
		}

		//build & send request to edge.
		req := types.BuildRequest(di.EdgeID, di.ProtocolType, "create", types.MSG_OPS_LIFE_CONTROL)
		req.SetContent(devSpec)

		resp, err := SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)
		if err != nil {
			return err
		}

		if resp.Payload.Code != global.IRespCodeOk {
			return fmt.Errorf("Error (code %s reason %s) from edge", resp.Payload.Code, resp.Payload.Content)
		}

		//store the device status.
		di.DeviceStatus = global.DeviceStatusActive
		return dsm.UpdateDeviceStatus(NewDeviceStatusDetails(&di))

	case global.DeviceStart:
		if devInfo.IsInactive() {
			return global.ErrDeviceIsOffline
		}

		if devInfo.IsOffline() {
			return global.ErrDeviceIsOffline
		}

		if devInfo.IsStarted() {
			return nil
		}

		klog.Infof("Start device (%s)", deviceID)
		devSpec := &v1.DeviceSpecMeta{
			DeviceID: di.DeviceID,
		}
		//build request to edge.
		req := types.BuildRequest(di.EdgeID, di.ProtocolType, "start", types.MSG_OPS_LIFE_CONTROL)
		req.SetContent(devSpec)

		resp, err := SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)
		if err != nil {
			return err
		}
		if resp.Payload.Code != global.IRespCodeOk {
			return fmt.Errorf("Error (code %s reason %s) from edge", resp.Payload.Code, resp.Payload.Content)
		}

		//update the start/stop state
		di.State = global.DeviceStateStarted
		di.DeviceStatus = ""
		return dsm.UpdateDeviceStatus(NewDeviceStatusDetails(&di))

	case global.DeviceStop:
		if devInfo.IsInactive() {
			return global.ErrDeviceIsOffline
		}

		if devInfo.IsOffline() {
			return global.ErrDeviceIsOffline
		}

		if !devInfo.IsStarted() {
			return nil
		}

		klog.Infof("Stop device (%s)", deviceID)

		devSpec := &v1.DeviceSpecMeta{
			DeviceID: di.DeviceID,
		}
		//build request to edge.
		req := types.BuildRequest(di.EdgeID, di.ProtocolType, "stop", types.MSG_OPS_LIFE_CONTROL)
		req.SetContent(devSpec)

		resp, err := SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)
		if err != nil {
			return err
		}

		if resp.Payload.Code != global.IRespCodeOk {
			return fmt.Errorf("Error (code %s reason %s) from edge", resp.Payload.Code, resp.Payload.Content)
		}

		//update the start/stop state
		di.State = global.DeviceStateStopped
		di.DeviceStatus = ""
		return dsm.UpdateDeviceStatus(NewDeviceStatusDetails(&di))

	case global.DeviceUpdate:
		if devInfo.IsInactive() {
			return global.ErrDeviceIsOffline
		}

		if devInfo.IsOffline() {
			return global.ErrDeviceIsOffline
		}

		klog.Infof("update device (%s)", deviceID)

		dm, err := db.GetDeviceModelAllInfoByID(di.DeviceModelId)
		if err != nil {
			return err
		}

		if dm == nil {
			return global.ErrNoSuchDeviceModel
		}

		devSpec := newDeviceSpecMeta(&di, dm)
		if devSpec == nil {
			return global.ErrNoSuchDevice
		}

		//build request to edge.
		devSpecList := []*v1.DeviceSpecMeta{devSpec}
		req := types.BuildRequest(di.EdgeID, di.ProtocolType, "update", types.MSG_OPS_LIFE_CONTROL)
		req.SetContent(devSpecList)

		resp, err := SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)
		if err != nil {
			return err
		}

		if resp.Payload.Code != global.IRespCodeOk {
			return fmt.Errorf("Error (code %s reason %s) from edge", resp.Payload.Code, resp.Payload.Content)
		}
		return nil

	case global.DeviceDelete:
		klog.Infof("Delete device (%s)", deviceID)
		//we delete it firstly.
		err := dsm.DeleteDeviceStatus(deviceID)
		if err != nil {
			klog.Warningf("DeleteDeviceStatus with err %v", err)
			return err
		}

		devSpec := &v1.DeviceSpecMeta{
			DeviceID: di.DeviceID,
		}

		//build request to edge and send.
		// we has no need to return the result of the edge, since
		// we will delete it later when device_status reportted.
		req := types.BuildRequest(di.EdgeID, di.ProtocolType, "delete", types.MSG_OPS_LIFE_CONTROL)
		req.SetContent(devSpec)
		SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)

		return nil
	}

	return global.ErrInvalidParms
}

func (ic *ICore) updateDeviceStatus(edgeID string, devicesStatus []*v1.DeviceStatusMessage) error {
	if edgeID == "" || len(devicesStatus) == 0 {
		return nil
	}

	dsm := ic.devStatusMgr
	for _, deviceStatus := range devicesStatus {
		deviceID := deviceStatus.DeviceID

		mutex := dsm.DeviceStatusLock(deviceID)

		devInfo := dsm.GetDeviceStatus(deviceID)
		if devInfo == nil {
			//No such device, we should delete it.
			dsm.DeviceStatusUnlock(mutex)
			klog.Warningf("No such device %s, we will delete it", deviceID)
			go DoDeviceAction(deviceID, global.DeviceDelete)
			continue
		}

		switch deviceStatus.Status {
		case global.DeviceStatusOnline:
			if devInfo.DeviceStatus != global.DeviceStatusOnline {
				if devInfo.DeviceStatus == global.DeviceStatusOffline {
					fetchHistory := dsm.GetFetchHistory(edgeID, devInfo.ProtocolType)
					/*
					* If the mapper has fetched when from offline to online, we consider
					* the mapper/edge is doing reboot/startup, and no need to sync all
					* devices to the mapper.
					* If the mapper has no fetch operation when offline to online, we consider
					* the edge is lost connect. then, we will update all syncs.
					 */
					if fetchHistory == nil {
						//we will update all device to the mapper.
						devSpecList, err := GetDeviceSpecMeta(edgeID, devInfo.ProtocolType)
						if err == nil {
							//build request to edge.
							klog.Infof("XXXXXXXXX Update Device(%s) XXXXXXXXXXXXXXXX", edgeID)
							req := types.BuildRequest(edgeID, devInfo.ProtocolType, "update", types.MSG_OPS_LIFE_CONTROL)
							req.SetContent(devSpecList)

							resp, err := SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)
							if err == nil {
								if resp.Payload.Code == global.IRespCodeOk {
									//mark the fetch flag.
									dsm.UpdateFetchHistory(edgeID, devInfo.ProtocolType, true)
								}
							}
						}
					}
				}

				devInfo.DeviceStatus = global.DeviceStatusOnline

				err := dsm.UpdateDeviceStatus(devInfo)
				if err != nil {
					klog.Warningf("UpdateDeviceStatus with err %s", err.Error())
				}

				/*
				* We sync the device desired property value to device's sides.
				* when device from offline to online.
				 */
				go ic.syncDeviceDesiredValuesToEdge(deviceID)
			}
		default:
			//we consider it as offline.
			devInfo.DeviceStatus = global.DeviceStatusOffline
			err := dsm.UpdateDeviceStatus(devInfo)
			if err != nil {
				klog.Warningf("UpdateDeviceStatus with err %s", err.Error())
			}

			db.InitAlertLogByDeviceId(devInfo.DeviceID, v1.AlertLogInvalid)
		}

		//unlock the device status.
		dsm.DeviceStatusUnlock(mutex)
	}

	return nil
}

func DecodeReportEventMessage(content string) (*v1.ReportEventMsg, error) {
	msg := &v1.ReportEventMsg{}

	//decode data.
	err := json.Unmarshal([]byte(content), msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func newDeviceSpecMeta(di *db.DeviceInstance, dm *db.DeviceModel) *v1.DeviceSpecMeta {
	if di == nil || dm == nil {
		return nil
	}

	dev := &v1.DeviceSpecMeta{}

	dev.DeviceID = di.DeviceID
	dev.DeviceOS = di.DeviceOS
	dev.DeviceCatagory = di.DeviceCategory
	dev.DeviceIdentificationCode = di.DeviceIdentificationCode
	dev.Protocol = *di.Protocol
	dev.State = di.State

	//decode the tags.
	json.Unmarshal([]byte(di.Tags), &dev.Tags)

	// Required: List of device services.
	dev.Services = make([]*v1.DeviceServiceSpec, 0)

	for _, si := range di.ServiceInstances {
		if si == nil {
			continue
		}

		dss := newDeviceServiceSpec(si, dm)
		if dss != nil {
			dev.Services = append(dev.Services, dss)
		}
	}

	if len(dev.Services) == 0 {
		return nil
	}

	return dev
}

func newDeviceServiceSpec(si *db.ServiceInstance, dm *db.DeviceModel) *v1.DeviceServiceSpec {
	var sm *db.ServiceModel

	for _, sm = range dm.ServiceModels {
		if sm == nil {
			continue
		}

		if sm.Name == si.Name {
			break
		}
	}

	if sm == nil {
		klog.Warningf("%s has no service model", si.Name)
		return nil
	}

	dss := &v1.DeviceServiceSpec{
		Name:       si.Name,
		Properties: make([]*v1.DevicePropertySpec, 0),
		Commands:   make([]*v1.DeviceCommandSpec, 0),
		Events:     make([]*v1.DeviceEventSpec, 0),
	}

	//property.
	for _, pi := range si.PropertyInstances {
		if pi == nil {
			continue
		}

		propertySpec := newDevicePropertySpec(pi, sm)
		if propertySpec != nil {
			dss.Properties = append(dss.Properties, propertySpec)
		}
	}

	//command
	for _, ci := range si.CommandInstances {
		if ci == nil {
			continue
		}

		commandSpec := newDeviceCommandSpec(ci, sm)
		if commandSpec != nil {
			dss.Commands = append(dss.Commands, commandSpec)
		}
	}

	//event.
	for _, ei := range si.EventInstances {
		if ei == nil {
			continue
		}

		eventSpec := newDeviceEventSpec(ei, sm)
		if eventSpec != nil {
			dss.Events = append(dss.Events, eventSpec)
		}
	}

	//If no config, we no need to generate dss.
	if len(dss.Events)+len(dss.Commands)+len(dss.Properties) <= 0 {
		return nil
	}

	return dss
}

func newDevicePropertySpec(pi *db.PropertyInstance, sm *db.ServiceModel) *v1.DevicePropertySpec {
	var pm *db.PropertyModel

	for _, pm = range sm.PropertyModels {
		if pm == nil {
			continue
		}

		if pm.Name == pi.Name {
			break
		}
	}

	if pm == nil || pi.AccessConfig == "" {
		klog.Warningf("property %s has no model or accessConfig is empty", pi.Name)
		return nil
	}

	dpm := &v1.DevicePropertyModel{
		Name:      pm.Name,
		WriteAble: pm.WriteAble,
		MaxValue:  pm.MaxValue,
		MinValue:  pm.MinValue,
		Unit:      pm.Unit,
		DataType:  pm.DataType,
	}

	return &v1.DevicePropertySpec{
		DevicePropertyModel: dpm,
		AccessConfig:        pi.AccessConfig,
	}
}

func newDeviceCommandSpec(ci *db.CommandInstance, sm *db.ServiceModel) *v1.DeviceCommandSpec {
	var cm *db.CommandModel

	for _, cm = range sm.CommandModels {
		if cm == nil {
			continue
		}

		if cm.Name == ci.Name {
			break
		}
	}

	if cm == nil || ci.AccessConfig == "" {
		klog.Warningf("command %s has no model or accessConfig is empty", ci.Name)
		return nil
	}

	dcm := &v1.DeviceCommandModel{
		Name: cm.Name,
	}

	json.Unmarshal([]byte(cm.RequestParam), &dcm.RequestParam)

	return &v1.DeviceCommandSpec{
		DeviceCommandModel: dcm,
		AccessConfig:       ci.AccessConfig,
	}
}

func newDeviceEventSpec(ei *db.EventInstance, sm *db.ServiceModel) *v1.DeviceEventSpec {
	var em *db.EventModel

	for _, em = range sm.EventModels {
		if em == nil {
			continue
		}

		if em.Name == ei.Name {
			break
		}
	}

	if em == nil || ei.AccessConfig == "" {
		klog.Warningf("event %s has no model or accessConfig is empty", ei.Name)
		return nil
	}

	dem := &v1.DeviceEventModel{
		Name:      em.Name,
		EventType: em.EventType,
	}

	return &v1.DeviceEventSpec{
		DeviceEventModel: dem,
		AccessConfig:     ei.AccessConfig,
	}
}

// API: StoreDevice.
// store all the device info by device spec.
func StoreDevice(spec *v1.DeviceSpec) (*db.DeviceInstance, error) {
	if spec.Name == "" {
		return nil, global.ErrInvalidParms
	}

	if spec.EdgeID == "" {
		return nil, global.ErrInvalidParms
	}
	//TODO: Check EdgeID.

	if spec.DeviceModelRef == "" {
		return nil, global.ErrInvalidParms
	}

	dm, err := db.GetDeviceModelAllInfoByName(spec.DeviceModelRef)
	if err != nil {
		return nil, err
	}

	if spec.ProtocolType == "" {
		return nil, global.ErrInvalidParms
	}

	protocol := ""
	if spec.Protocol != "" {
		protocol = spec.Protocol
	} else {
		p, err := json.Marshal(&v1.DefaultProtocol{
			IntervalUnit: "s",
			Interval:     10,
		})
		if err != nil {
			klog.Errorf("json Marshal with  err %v", err)
		}
		protocol = string(p)
	}

	//default is 5s for desired value life time.
	ltodv := global.DefaultLifeTimeOfDesiredValue
	if spec.LifeTimeOfDesiredValue >= 0 {
		ltodv = time.Millisecond * time.Duration(spec.LifeTimeOfDesiredValue)
	}

	doc := &db.DeviceInstance{
		DeviceID:                 utils.NewUUID(),
		Name:                     spec.Name,
		EdgeID:                   spec.EdgeID,
		DeviceOS:                 spec.DeviceOS,
		DeviceCategory:           spec.DeviceOS,
		DeviceVersion:            spec.DeviceVersion,
		DeviceIdentificationCode: spec.DeviceIdentificationCode,
		Description:              &spec.Description,
		GroupName:                spec.GroupName,
		GroupID:                  "",
		Health:                   0,
		Creator:                  spec.Creator,
		DeviceAuthType:           spec.DeviceAuthType,
		Secret:                   spec.DeviceAuthType,
		DeviceType:               spec.DeviceType,
		GatewayID:                spec.GatewayID,
		GatewayName:              spec.GatewayName,
		DeviceModelRef:           spec.DeviceModelRef,
		ProtocolType:             spec.ProtocolType,
		Protocol:                 &protocol,
		LifeTimeOfDesiredValue:   int64(ltodv),
		CreateTimeStamp:          utils.GetNowTimeStamp(),
		UpdateTimeStamp:          utils.GetNowTimeStamp(),
		DeviceStatus:             global.DeviceStatusInactive,
		State:                    global.DeviceStateStopped,
		DeviceModelId:            dm.ID,
	}

	if spec.Tags != nil {
		bytes, err := json.Marshal(spec.Tags)
		if err != nil {
			klog.Errorf("Marshal Tags with error %v", err)
			return nil, err
		} else {
			doc.Tags = string(bytes)
		}
	}

	err = db.AddDeviceInstance(doc)
	if err != nil {
		return nil, err
	}

	//TODO: ProtocolType is exist ?
	if spec.ExtensionConfig == nil {
		spec.ExtensionConfig = &v1.ExtensionConfig{}
	}

	err = spec.ExtensionConfig.StoreExtensionConfig(dm, doc.DeviceID)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
