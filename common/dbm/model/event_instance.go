package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type EventInstance struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; type:varchar(64); not null" json:"name"`
	ServiceID       string `gorm:"column:service_id; not null" json:"serviceId"`
	AccessConfig    string `gorm:"column:access_config; type:varchar(4096); not null" json:"accessConfig,omitempty"`
	UpdateTimeStamp int64  `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (EventInstance) TableName() string {
	return "event_instance"
}
func NewEventInstance(name string, serviceID string, accessConfig string) *EventInstance {
	return &EventInstance{
		Name:         name,
		ServiceID:    serviceID,
		AccessConfig: accessConfig,
	}
}

func AddEventInstance(docs []*EventInstance) error {
	err := global.DBAccess.Create(&docs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func UpdateEventAccessConfig(serviceID string, name string, accessConfig string) error {
	err := global.DBAccess.Model(&EventInstance{}).Where("service_id = ? and name = ?", serviceID, name).Update("access_config", accessConfig).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func GetEventInstanceByServiceIdAndName(serviceID string, name string) (*EventInstance, error) {
	eventInstance := &EventInstance{}
	err := global.DBAccess.Model(&EventInstance{}).Where("service_id = ? and name = ?", serviceID, name).First(eventInstance).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return eventInstance, err
	}
	return eventInstance, nil
}
func IsExistEventInstance(serviceID string, name string) bool {
	var count int64
	err := global.DBAccess.Model(&EventInstance{}).Where(EventInstance{ServiceID: serviceID, Name: name}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func GetEventInstanceByServiceId(serviceID string) ([]*EventInstance, error) {
	var eventInstances []*EventInstance
	err := global.DBAccess.Model(&EventInstance{}).Where("service_id = ?", serviceID).Order("update_time_stamp desc").Find(&eventInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return eventInstances, err
	}
	return eventInstances, nil
}

func GetEventInstanceByDeviceId(deviceID string) ([]*EventInstance, error) {
	var eventInstances []*EventInstance
	serviceInstances, err := GetServiceInstanceByDeviceID(deviceID)
	if err != nil {
		return eventInstances, err
	}
	for _, serviceInstance := range serviceInstances {
		if events, err := GetEventInstanceByServiceId(serviceInstance.ID); err == nil {
			eventInstances = append(eventInstances, events...)
		}
	}
	return eventInstances, nil
}

func DeleteEventInstance(name string, serviceID string) error {
	err := global.DBAccess.Where("name = ? and service_id = ?", name, serviceID).Delete(&EventInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	if eventInstance, err := GetEventInstanceByServiceIdAndName(serviceID, name); err == nil {
		DeleteEventRuleRelationByEventId(eventInstance.ID)
	}
	return nil
}

func DeleteEventInstanceByServiceId(serviceID string) error {

	err := global.DBAccess.Where("service_id = ?", serviceID).Delete(&EventInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	if eventInstances, err := GetEventInstanceByServiceId(serviceID); err == nil {
		for _, eventInstance := range eventInstances {
			DeleteEventRuleRelationByEventId(eventInstance.ID)
		}
	}
	return nil
}
