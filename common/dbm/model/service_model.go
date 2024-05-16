package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type ServiceModel struct {
	ID              int64            `gorm:"primary_key; auto_increment" json:"id"`
	Name            string           `gorm:"column:name; not null; type:varchar(256)" json:"name"`
	Description     string           `gorm:"column:description; default:null; type:varchar(256);" json:"description"`
	UpdateTimeStamp int64            `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
	PropertyModels  []*PropertyModel `gorm:"foreignKey: ServiceModelId; references:id"`
	EventModels     []*EventModel    `gorm:"foreignKey: ServiceModelId; references:id"`
	CommandModels   []*CommandModel  `gorm:"foreignKey: ServiceModelId; references:id"`
	DeviceModelId   int64            `gorm:"column:device_model_id" json:"deviceModelId"`
}

func (ServiceModel) TableName() string {
	return "service_model"
}

func GetServiceModelByDeviceModelId(deviceModelId int64) ([]*ServiceModel, error) {
	var serviceModel []*ServiceModel

	err := global.DBAccess.Preload("PropertyModels").Preload("EventModels").Preload("CommandModels").Order("update_time_stamp desc").Where(ServiceModel{DeviceModelId: deviceModelId}).Find(&serviceModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return serviceModel, err
}

func GetServiceModelByServiceModelId(serviceModelId int64) (*ServiceModel, error) {
	var serviceModel *ServiceModel
	err := global.DBAccess.First(&serviceModel, serviceModelId).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return serviceModel, err
}
func GetServiceModelByDeviceModelIdAndServiceModelName(deviceModelId int64, serviceModelName string) (*ServiceModel, error) {
	var serviceModel *ServiceModel
	err := global.DBAccess.Where(&ServiceModel{DeviceModelId: deviceModelId, Name: serviceModelName}).First(&serviceModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return serviceModel, err
}

func IsExistServiceModel(deviceModelId int64, serviceModelName string) bool {
	var count int64
	err := global.DBAccess.Model(&ServiceModel{}).Where(&ServiceModel{DeviceModelId: deviceModelId, Name: serviceModelName}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true

}
func AddServiceModel(serviceModel *ServiceModel) error {
	err := global.DBAccess.Create(&serviceModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteServiceModel(id int64) error {
	//global.DBAccess.Begin()
	err := global.DBAccess.Delete(&ServiceModel{}, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		//global.DBAccess.Rollback()
		return err
	}
	//global.DBAccess.Commit()
	return nil
}
