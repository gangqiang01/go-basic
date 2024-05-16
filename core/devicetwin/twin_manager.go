package devicetwin

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/influxdbm/influx_store"
	"github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/core/devicetwin/store"
	"github.com/edgehook/ithings/dataforward"
	"k8s.io/klog/v2"
)

/*
* DeviceTwinManager
* This is a device twin/shadow manager to
* manage all device twin/shadow.
 */
type DeviceTwinManager struct {
	//data store interface.
	DataStore store.DBStore
}

func NewDeviceTwinManager() *DeviceTwinManager {
	return &DeviceTwinManager{
		DataStore: store.NewDBStore("memory_store"),
	}
}

func (dtm *DeviceTwinManager) Initialize() error {
	return dtm.DataStore.Initialize()
}

// manager the device twin data from edge report.
func (dtm *DeviceTwinManager) AddDevicesTwinsUpdateData(msg *v1.ReportDevicesMessage, edgeId string) error {
	if msg == nil {
		return nil
	}

	for _, devMsg := range msg.Devices {
		if devMsg == nil {
			continue
		}

		deviceID := devMsg.DeviceID
		reportedTwins := devMsg.Services

		err := dtm.DataStore.StoreReportedTwins(deviceID, reportedTwins)
		if err != nil {
			klog.Warningf("Store %s twins failed with err %v", deviceID, err)
			continue
		}

		//store influxDB
		if err := influx_store.StoreTwin(deviceID, reportedTwins); err != nil {
			klog.Warningf("Store influxDB %s twins failed with err %v", deviceID, err)
			continue
		}

		//data forward
		if err := dataforward.HandleDataForward(devMsg.DeviceID, edgeId, "property", devMsg); err != nil {
			klog.Errorf("Data forward with err %v", err)
			continue
		}
	}
	return nil
}

func (dtm *DeviceTwinManager) GetDeviceTwins(deviceID string) []*store.TwinsData {
	twinsData := dtm.DataStore.LoadReportedTwins(deviceID, 20)
	return twinsData
}

func (dtm *DeviceTwinManager) UpdateDesiredTwins(devTwinMsg *v1.DeviceTwinMessage) error {
	if devTwinMsg == nil {
		return global.ErrInvalidParms
	}

	deviceID := devTwinMsg.DeviceID
	desiredTwins := devTwinMsg.DesiredTwins

	return dtm.DataStore.UpdateDesiredTwins(deviceID, desiredTwins)
}

func (dtm *DeviceTwinManager) GetDesiredTwins(deviceID string, duration int64) []*v1.TwinProperty {
	twins := make([]*v1.TwinProperty, 0)
	nowTimeStamp := utils.GetNowTimeStamp()

	twinData := dtm.DataStore.GetDesiredTwins(deviceID)
	if twinData == nil || twinData.Twins == nil {
		return twins
	}

	for _, twin := range twinData.Twins {
		if twin == nil || twin.Value == nil {
			continue
		}

		//reached out the twin's lifetime, we skip it.
		if nowTimeStamp-twin.Timestamp > duration {
			continue
		}

		twins = append(twins, twin)
	}

	return twins
}
