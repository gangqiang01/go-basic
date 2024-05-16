package core

import (
	"sync"
	"time"

	db "github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

/*
* this is device status manager
* manage all device's status. to process some
* down link issue between server and agent, we
* add device status manager to ensure that after down link
* and reconnect successful, we sync device info correctly.
 */
type DeviceStatusManager struct {
	FetchHistory         *sync.Map
	DeviceStatusCache    *sync.Map
	deviceStatusMutexMap *sync.Map
}

func NewDeviceStatusManager() *DeviceStatusManager {
	var historyMap sync.Map
	var deviceStatusCache sync.Map
	var devStatusMutexMap sync.Map

	return &DeviceStatusManager{
		FetchHistory:         &historyMap,
		DeviceStatusCache:    &deviceStatusCache,
		deviceStatusMutexMap: &devStatusMutexMap,
	}
}

type HistoryRecord struct{}

func (dsm *DeviceStatusManager) UpdateFetchHistory(edgeId, protocType string, val bool) {
	key := edgeId + "_" + protocType

	record := dsm.GetFetchHistory(edgeId, protocType)

	if val {
		if record == nil {
			record = &HistoryRecord{}
			dsm.FetchHistory.Store(key, record)
		}
	} else {
		//we delete it.
		if record != nil {
			dsm.FetchHistory.Delete(key)
		}
	}
}

func (dsm *DeviceStatusManager) GetFetchHistory(edgeId, protocType string) *HistoryRecord {
	key := edgeId + "_" + protocType

	v, exist := dsm.FetchHistory.Load(key)
	if !exist {
		return nil
	}

	record, isThisType := v.(*HistoryRecord)
	if !isThisType {
		return nil
	}

	return record
}

type DeviceStatusDetails struct {
	DeviceID               string
	EdgeID                 string
	ProtocolType           string
	UpdateTimeStamp        int64
	LifeTimeOfDesiredValue int64
	DeviceStatus           string
	State                  string
}

func NewDeviceStatusDetails(dev *db.DeviceInstance) *DeviceStatusDetails {
	return &DeviceStatusDetails{
		DeviceID:               dev.DeviceID,
		EdgeID:                 dev.EdgeID,
		ProtocolType:           dev.ProtocolType,
		UpdateTimeStamp:        dev.UpdateTimeStamp,
		LifeTimeOfDesiredValue: dev.LifeTimeOfDesiredValue,
		DeviceStatus:           dev.DeviceStatus,
		State:                  dev.State,
	}
}

func (dsd *DeviceStatusDetails) ToDBDeviceInstance() *db.DeviceInstance {
	return &db.DeviceInstance{
		DeviceID:               dsd.DeviceID,
		EdgeID:                 dsd.EdgeID,
		ProtocolType:           dsd.ProtocolType,
		UpdateTimeStamp:        dsd.UpdateTimeStamp,
		LifeTimeOfDesiredValue: dsd.LifeTimeOfDesiredValue,
		DeviceStatus:           dsd.DeviceStatus,
		State:                  dsd.State,
	}
}

func (dsd *DeviceStatusDetails) IsInactive() bool {
	return dsd.DeviceStatus == global.DeviceStatusInactive
}

func (dsd *DeviceStatusDetails) IsOffline() bool {
	return dsd.DeviceStatus == global.DeviceStatusOffline
}

func (dsd *DeviceStatusDetails) IsOnline() bool {
	return dsd.DeviceStatus == global.DeviceStatusOnline
}

func (dsd *DeviceStatusDetails) IsStarted() bool {
	return dsd.State == global.DeviceStateStarted
}

func (dsm *DeviceStatusManager) getDeviceStatusLock(deviceID string) *sync.Mutex {
	v, exist := dsm.deviceStatusMutexMap.Load(deviceID)
	if !exist {
		return nil
	}

	mutex, isMutex := v.(*sync.Mutex)
	if !isMutex {
		return nil
	}

	return mutex
}

func (dsm *DeviceStatusManager) createDeviceStatusMutex(deviceID string) *sync.Mutex {
	mutex := dsm.getDeviceStatusLock(deviceID)
	if mutex != nil {
		return mutex
	}

	var deviceStatusMutex sync.Mutex
	mutex = &deviceStatusMutex

	dsm.deviceStatusMutexMap.Store(deviceID, mutex)
	return mutex
}

func (dsm *DeviceStatusManager) DeviceStatusLock(deviceID string) *sync.Mutex {
	mutex := dsm.getDeviceStatusLock(deviceID)
	if mutex == nil {
		mutex = dsm.createDeviceStatusMutex(deviceID)
	}

	mutex.Lock()
	return mutex
}

func (dsm *DeviceStatusManager) DeviceStatusUnlock(mutex *sync.Mutex) {
	if mutex == nil {
		return
	}
	mutex.Unlock()
}

func (dsm *DeviceStatusManager) DeleteDeviceStatusMutex(deviceID string) {
	mutex := dsm.getDeviceStatusLock(deviceID)
	if mutex != nil {
		dsm.deviceStatusMutexMap.Delete(deviceID)
	}
}

func (dsm *DeviceStatusManager) UpdateDeviceStatus(dev *DeviceStatusDetails) error {
	var isChanged = bool(false)

	deviceID := dev.DeviceID

	device := dsm.GetDeviceStatus(deviceID)
	if device == nil {
		return global.ErrNoSuchDevice
	}

	if dev.DeviceStatus != "" {
		isChanged = true
		device.DeviceStatus = dev.DeviceStatus
	}
	if dev.State != "" {
		isChanged = true
		device.State = dev.State
	}

	//update device status.
	if isChanged {
		device.UpdateTimeStamp = time.Now().Unix()
		dsm.DeviceStatusCache.Store(deviceID, device)

		//update device status at database.
		doc := dev.ToDBDeviceInstance()
		return db.UpdateDeviceInstance(deviceID, doc)
	}

	return nil
}

func (dsm *DeviceStatusManager) UpdateAllDeviceStatusInThisEdge(edgeID, status string) error {
	//update device status In database.
	err := db.UpdateAllDevInstStatusInThisEdge(edgeID, status)
	if err != nil {
		klog.Errorf("UpdateAllDevInstStatusInThisEdge with err: %v", err)
		return err
	}

	//update device status In cache.
	dsm.DeviceStatusCache.Range(func(key, value interface{}) bool {
		devInfo, isThisType := value.(*DeviceStatusDetails)
		if !isThisType || devInfo == nil {
			return true
		}

		if devInfo.EdgeID != edgeID {
			return true
		}

		if devInfo.DeviceStatus != status {
			deviceID := devInfo.DeviceID

			mutex := dsm.DeviceStatusLock(deviceID)

			devInfo.DeviceStatus = status
			devInfo.UpdateTimeStamp = time.Now().Unix()

			dsm.DeviceStatusCache.Store(deviceID, devInfo)

			dsm.DeviceStatusUnlock(mutex)
		}

		return true
	})

	return nil
}

func (dsm *DeviceStatusManager) GetDeviceStatus(deviceID string) *DeviceStatusDetails {
	var devInfo DeviceStatusDetails

	v, exist := dsm.DeviceStatusCache.Load(deviceID)
	if !exist {
		//we lookup it on database.
		devInfo := db.GetDeviceInstance(deviceID)
		if devInfo == nil {
			return nil
		}

		//find it into database, we store it into cache.
		v = NewDeviceStatusDetails(devInfo)
		// store it into cache.
		dsm.DeviceStatusCache.Store(deviceID, v)
	}

	dev, isThisType := v.(*DeviceStatusDetails)
	if !isThisType {
		//we lookup it on database.
		devInfo := db.GetDeviceInstance(deviceID)
		if devInfo == nil {
			return nil
		}

		//find it into database, we store it into cache.
		dev = NewDeviceStatusDetails(devInfo)
		// store it into cache.
		dsm.DeviceStatusCache.Store(deviceID, dev)
	}

	devInfo = *dev
	return &devInfo
}

// DeleteDeviceStatus will delete device in db.
func (dsm *DeviceStatusManager) DeleteDeviceStatus(deviceID string) error {
	if dsm.GetDeviceStatus(deviceID) == nil {
		return nil
	}

	//delete device status
	dsm.DeviceStatusCache.Delete(deviceID)

	//then we delete the mutex.
	dsm.DeleteDeviceStatusMutex(deviceID)

	return db.DeleteDeviceInstance(deviceID)
}
