package core

import (
	"context"
	"fmt"
	"time"

	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
	"github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

func GetDeviceTwin(deviceID string) *v1.DeviceTwinMessage {
	twinsData := defaultICore.iCore.twinMgr.GetDeviceTwins(deviceID)
	twinPropertys := defaultICore.iCore.GetDesiredTwins(deviceID)
	message := v1.NewDeviceTwinMessage(deviceID)
	for _, td := range twinsData {
		if td != nil {
			message.AddTwins(twinPropertys, td.Twins)
		}
	}
	return message
}

// API: UpdateDeviceDesiredTwins
// update device desired twins.
// we will update the desired twins into datastore firstly, then send sync update
// message to edge if device is online. Or, we will send all device desired
// twins to edge device when device is from offline to online on the next.
// If device is online, but sync failed, we will store the update to a golobal
// cache, and retry to sync the update message to edge device until the lifetime
// is reached!
func UpdateDeviceDesiredTwins(devTwinMsg *v1.DeviceTwinMessage) error {
	if devTwinMsg == nil {
		return global.ErrInvalidParms
	}
	if defaultICore == nil {
		return global.ErrCoreNotReady
	}

	twinMgr := defaultICore.iCore.twinMgr
	deviceID := devTwinMsg.DeviceID

	//store the twin into data store.
	err := twinMgr.UpdateDesiredTwins(devTwinMsg)
	if err != nil {
		return err
	}

	dsm := defaultICore.iCore.devStatusMgr
	devInfo := dsm.GetDeviceStatus(deviceID)
	if devInfo == nil {
		return global.ErrUnknown
	}

	//If device is offline, we return directlly. and wait
	// device online and send these update to edge.
	if !devInfo.IsOnline() {
		return nil
	}

	//run gorountine to process it.
	go func() {
		//If device is online, we send it to edge mapper.
		err := defaultICore.iCore.SendDesiredTwinsUpdate(deviceID, devTwinMsg.DesiredTwins)
		if err == nil {
			return
		}

		//If send failed, we will retry to send util reach
		//out the lifetime of the twin value.
		klog.Infof("Sync device %s faild, we add these into retry sync cache. ", deviceID)
		defaultICore.iCore.AddDesiredTwinIntoRetryCache(deviceID, devTwinMsg.DesiredTwins)
	}()

	return nil
}

// send desired twin update.
func (ic *ICore) SendDesiredTwinsUpdate(deviceID string, desiredTwins []*v1.TwinProperty) error {
	dsm := ic.devStatusMgr

	devInfo := dsm.GetDeviceStatus(deviceID)
	if devInfo == nil {
		return global.ErrNoSuchDevice
	}

	devTwinUpdateMsg := &v1.DeviceDesiredTwinsUpdateMessage{
		DeviceID:     deviceID,
		DesiredTwins: desiredTwins,
	}

	req := types.BuildRequest(devInfo.EdgeID, devInfo.ProtocolType, "desired_twins", types.MSG_OPS_SET_PROPERTY)
	req.SetContent(devTwinUpdateMsg)

	resp, err := ic.SendSyncRequestToEdge(req, global.DefaultEdgeMaxResponseTime)
	if err != nil {
		return err
	}

	if resp.Payload.Code != global.IRespCodeOk {
		return fmt.Errorf("Error (code %s reason %s) from edge", resp.Payload.Code, resp.Payload.Content)
	}

	return nil
}

// get desired property value not reached lifetime.
func (ic *ICore) GetDesiredTwins(deviceID string) []*v1.TwinProperty {
	dsm := ic.devStatusMgr
	twinMgr := ic.twinMgr

	devInfo := dsm.GetDeviceStatus(deviceID)
	if devInfo == nil {
		return nil
	}

	duration := devInfo.LifeTimeOfDesiredValue

	return twinMgr.GetDesiredTwins(deviceID, duration)
}

// sync device's desired value to edge
func (ic *ICore) syncDeviceDesiredValuesToEdge(deviceID string) error {
	desiredTwins := ic.GetDesiredTwins(deviceID)
	if desiredTwins == nil || len(desiredTwins) == 0 {
		klog.Errorf("GetDesiredTwins(deviceID=%s) failed", deviceID)
		return fmt.Errorf("GetDesiredTwins(deviceID=%s) failed", deviceID)
	}

	//delete desired twin update in retry cache.
	ic.deleteDesiredTwinInRetryCache(deviceID)

	//send the sync message to edge.
	err := ic.SendDesiredTwinsUpdate(deviceID, desiredTwins)
	if err != nil {
		klog.Errorf("SendDesiredTwinsUpdate with error %s", err)
		return err
	}

	return nil
}

func (ic *ICore) runRetrySyncDesiredTwinLoop(ctx context.Context) {
	dsm := ic.devStatusMgr

	for {
		select {
		case <-ctx.Done():
			klog.Infof("exit Icore since the context canceled....")
			return
		default:
			time.Sleep(time.Millisecond * 300)
			ic.twinRetrySyncCache.Range(func(key, value interface{}) bool {
				deviceID, isThisType := key.(string)
				desiredTwins, isThisType1 := value.([]*v1.TwinProperty)
				if !isThisType || !isThisType1 {
					return true
				}

				devInfo := dsm.GetDeviceStatus(deviceID)
				if devInfo == nil {
					return true
				}

				duration := devInfo.LifeTimeOfDesiredValue
				twins := getDeviceTwinNotExpired(desiredTwins, duration)
				if twins == nil || len(twins) == 0 {
					//all is expired, we should delete it.
					ic.deleteDesiredTwinInRetryCache(deviceID)
					return true
				}

				//sync device desired twin to edge.
				err := ic.SendDesiredTwinsUpdate(deviceID, twins)
				if err == nil {
					//send successful, we delete it in the sync cache.
					ic.deleteDesiredTwinInRetryCache(deviceID)
				}

				return true
			})
		}
	}
}

func getDeviceTwinNotExpired(desiredTwins []*v1.TwinProperty, duration int64) []*v1.TwinProperty {
	twins := make([]*v1.TwinProperty, 0)
	nowTimeStamp := utils.GetNowTimeStamp()

	for _, twin := range desiredTwins {
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

func (ic *ICore) AddDesiredTwinIntoRetryCache(deviceID string, desiredTwins []*v1.TwinProperty) {
	twins := ic.findDesiredTwinInRetryCache(deviceID)
	if twins == nil {
		twins = desiredTwins
	} else {
		for _, twin := range desiredTwins {
			var found bool

			if twin == nil {
				continue
			}

			found = false
			svc, prop := twin.Service, twin.PropertyName

			for i, t := range twins {
				if t == nil {
					continue
				}

				if t.Service == svc && t.PropertyName == prop {
					twins[i] = twin
					found = true
					break
				}
			}

			if !found {
				twins = append(twins, twin)
			}
		}
	}

	ic.twinRetrySyncCache.Store(deviceID, twins)
}

func (ic *ICore) findDesiredTwinInRetryCache(deviceID string) []*v1.TwinProperty {
	v, exist := ic.twinRetrySyncCache.Load(deviceID)
	if !exist {
		return nil
	}

	twins, isThisType := v.([]*v1.TwinProperty)
	if !isThisType {
		return nil
	}

	return twins
}

func (ic *ICore) deleteDesiredTwinInRetryCache(deviceID string) {
	if ic.findDesiredTwinInRetryCache(deviceID) != nil {
		ic.twinRetrySyncCache.Delete(deviceID)
	}
}
