package model

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

type ServiceInstance struct {
	ID                string              `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name              string              `gorm:"column:name; type:varchar(64); not null" json:"name"`
	DeviceID          string              `gorm:"column:device_id; type:varchar(36); not null" json:"deviceId"`
	UpdateTimeStamp   int64               `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
	PropertyInstances []*PropertyInstance `gorm:"foreignKey: ServiceID; references:id"`
	EventInstances    []*EventInstance    `gorm:"foreignKey: ServiceID; references:id"`
	CommandInstances  []*CommandInstance  `gorm:"foreignKey: ServiceID; references:id"`
}

func (*ServiceInstance) TableName() string {
	return "service_instance"
}

func NewServiceInstance(name, deviceID string) *ServiceInstance {
	return &ServiceInstance{
		ID:       utils.NewUUID(),
		Name:     name,
		DeviceID: deviceID,
	}
}

func AddServiceInstance(docs []*ServiceInstance) error {
	err := global.DBAccess.Create(&docs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func GetServiceInstanceByDeviceID(deviceID string) ([]*ServiceInstance, error) {
	var serviceInstances []*ServiceInstance

	err := global.DBAccess.Preload("PropertyInstances").Preload("EventInstances").Preload("CommandInstances").Order("update_time_stamp desc").Where(ServiceInstance{DeviceID: deviceID}).Find(&serviceInstances).Error
	if err != nil {
		klog.Errorf("GetServiceInstanceByDeviceID err: %v", err)
		return serviceInstances, err
	}

	return serviceInstances, err
}

func GetServiceInstanceByDeviceIDAndName(deviceID string, name string) (*ServiceInstance, error) {
	var serviceInstances *ServiceInstance

	err := global.DBAccess.Where("device_id = ? and name = ?", deviceID, name).First(&serviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return serviceInstances, err
}

func IsExistServiceInstance(deviceID string, name string) bool {
	var count int64
	err := global.DBAccess.Model(&ServiceInstance{}).Where("device_id = ? and name = ?", deviceID, name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func DeleteServiceInstance(deviceID string) error {
	err := global.DBAccess.Where("device_id = ?", deviceID).Delete(&ServiceInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		//global.DBAccess.Rollback()
		return err
	}
	return nil
}

func GetServiceInstanceByDeviceNameAndServiceName(deviceId string, name string) (*ServiceInstance, error) {
	var serviceInstances *ServiceInstance
	deviceInstance, err := GetDeviceInstanceByDeviceId(deviceId)
	if err != nil {
		return serviceInstances, err
	}
	if err := global.DBAccess.Where("device_id = ? and name = ?", deviceInstance.DeviceID, name).First(&serviceInstances).Error; err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return serviceInstances, err
}
