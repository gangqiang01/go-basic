package model

import (
	"github.com/edgehook/ithings/common/global"

	"k8s.io/klog/v2"
)

type EventModel struct {
	ID              int64   `gorm:"primary_key; auto_increment" json:"id"`
	Name            string  `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	EventType       string  `gorm:"column:event_type; not null; type:varchar(256);" json:"eventType"`
	MaxValue        float64 `gorm:"column:max_value; type:float;" json:"maxValue"`
	MinValue        float64 `gorm:"column:min_vale; type:float;" json:"minValue"`
	Unit            string  `gorm:"column:unit; type:varchar(256);" json:"unit"`
	DataType        string  `gorm:"column:data_type; type:varchar(256);" json:"dataType"`
	Description     string  `gorm:"column:description; type:varchar(256);" json:"description"`
	UpdateTimeStamp int64   `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
	ServiceModelId  int64   `gorm:"column:service_model_id" json:"serviceModelId"`
}

func (EventModel) TableName() string {
	return "event_model"
}

func GetEventModelByServiceModelId(serviceId int64) ([]*EventModel, error) {
	var eventModel []*EventModel
	err := global.DBAccess.Where(EventModel{ServiceModelId: serviceId}).Find(&eventModel).Order("update_time_stamp desc").Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return eventModel, err
}

func GetEventModelByServiceIdAndName(serviceID int64, name string) (*EventModel, error) {
	var eventModel *EventModel
	err := global.DBAccess.Model(&EventModel{}).Where("service_model_id = ? and name = ?", serviceID, name).First(&eventModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return eventModel, err
	}
	return eventModel, nil
}

func GetEventModelByEventId(eventId int64) (*EventModel, error) {
	var eventModel *EventModel
	err := global.DBAccess.First(&eventModel, eventId).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return eventModel, err
	}
	return eventModel, err
}
func IsExistEventModel(serviceModelId int64, eventName string) bool {
	var count int64
	err := global.DBAccess.Model(&EventModel{}).Where(EventModel{ServiceModelId: serviceModelId, Name: eventName}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func AddEventModel(eventModel *EventModel) error {
	err := global.DBAccess.Create(&eventModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveEventModel(id int64, eventModel *EventModel) error {
	err := global.DBAccess.Model(&EventModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		//"Name":        eventModel.Name,
		"EventType":   eventModel.EventType,
		"MaxValue":    eventModel.MaxValue,
		"MinValue":    eventModel.MinValue,
		"Unit":        eventModel.Unit,
		"DataType":    eventModel.DataType,
		"Description": eventModel.Description,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteEventModel(id int64) error {
	//global.DBAccess.Begin()
	err := global.DBAccess.Delete(&EventModel{}, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		//global.DBAccess.Rollback()
		return err
	}
	//global.DBAccess.Commit()
	return nil
}

func IsExistEventInstanceByEventModelId(id int64) (bool, error) {
	eventModel, err := GetEventModelByEventId(id)
	if err != nil {
		return false, err
	}
	serviceModel, err := GetServiceModelByServiceModelId(eventModel.ServiceModelId)
	if err != nil {
		return false, err
	}

	deviceInstances, err := GetDeviceInstancesByDeviceModelId(serviceModel.DeviceModelId)

	if err != nil || len(deviceInstances) == 0 {
		return false, err
	}
	for _, deviceInstance := range deviceInstances {
		serviceInstances, err := GetServiceInstanceByDeviceID(deviceInstance.DeviceID)
		if err != nil {
			continue
		}
		if len(serviceInstances) > 0 {
			for _, serviceInstance := range serviceInstances {
				if res := IsExistEventInstance(serviceInstance.ID, eventModel.Name); res {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
