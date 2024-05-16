package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type PropertyModel struct {
	ID              int64   `gorm:"primary_key; auto_increment" json:"id"`
	Name            string  `gorm:"column:name; not null; type:varchar(256);" json:"name,omitempty"`
	WriteAble       bool    `gorm:"column:writeAble; type:bool;" json:"writeAble"`
	Report          bool    `gorm:"column:report; type:bool;" json:"report"`
	MaxValue        float64 `gorm:"column:maxValue; type:float;" json:"maxValue"`
	MinValue        float64 `gorm:"column:minVale; type:float;" json:"minValue"`
	Unit            string  `gorm:"column:unit; type:varchar(256);" json:"unit"`
	DataType        string  `gorm:"column:data_type; type:varchar(256);" json:"dataType"`
	UpdateTimeStamp int64   `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
	Description     string  `gorm:"column:description; type:varchar(256);" json:"description,omitempty"`
	ServiceModelId  int64   `gorm:"column:service_model_id" json:"serviceModelId,omitempty"`
}

func (PropertyModel) TableName() string {
	return "property_model"
}

func GetPropertyModelByServiceId(serviceId int64) ([]*PropertyModel, error) {
	var propertyModel []*PropertyModel
	err := global.DBAccess.Where("service_model_id=?", serviceId).Order("update_time_stamp desc").Find(&propertyModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return propertyModel, err
}

func GetPropertyModelByPageAndServiceId(serviceId int64, page int, limit int) ([]*PropertyModel, error) {
	var propertyModel []*PropertyModel
	err := global.DBAccess.Where("service_model_id=?", serviceId).Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&propertyModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return propertyModel, err
}

func GetPropertyModelByPropertyId(propertyId int64) (*PropertyModel, error) {
	var propertyModel *PropertyModel
	err := global.DBAccess.First(&propertyModel, propertyId).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return propertyModel, err
	}
	return propertyModel, err
}
func GetPropertyModelByServiceModelIdAndPropertyName(serviceModelId int64, propertyName string) (*PropertyModel, error) {
	var propertyModel *PropertyModel
	err := global.DBAccess.Where(&PropertyModel{ServiceModelId: serviceModelId, Name: propertyName}).First(&propertyModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return propertyModel, err
	}
	return propertyModel, err
}

func GetPropertyModelByDeviceModelIdAndServiceNameAndPropertyName(deviceModelId int64, serviceName string, propertyNmae string) (*PropertyModel, error) {
	var propertyModel *PropertyModel
	serviceModel, err := GetServiceModelByDeviceModelIdAndServiceModelName(deviceModelId, serviceName)
	if err != nil {
		klog.Errorf("err: %v", err)
		return propertyModel, err
	}

	propertyModel, err = GetPropertyModelByServiceModelIdAndPropertyName(serviceModel.ID, propertyNmae)
	if err != nil {
		klog.Errorf("err: %v", err)
		return propertyModel, err
	}
	return propertyModel, err
}

func IsExistPropertyModel(serviceModelId int64, propertyName string) bool {
	var count int64
	err := global.DBAccess.Model(&PropertyModel{}).Where(&PropertyModel{ServiceModelId: serviceModelId, Name: propertyName}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func AddPropertyModel(propertyModel *PropertyModel) error {
	err := global.DBAccess.Create(&propertyModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SavePropertyModel(id int64, propertyModel *PropertyModel) error {
	err := global.DBAccess.Model(&PropertyModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		//"Name":        &propertyModel.Name,
		"WriteAble":   propertyModel.WriteAble,
		"MaxValue":    propertyModel.MaxValue,
		"MinValue":    propertyModel.MinValue,
		"Unit":        propertyModel.Unit,
		"DataType":    propertyModel.DataType,
		"Description": propertyModel.Description,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeletePropertyModel(id int64) error {
	//global.DBAccess.Begin()
	err := global.DBAccess.Delete(&PropertyModel{}, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		//global.DBAccess.Rollback()
		return err
	}
	//global.DBAccess.Commit()
	return nil
}

func IsExistPropertyInstanceByPropertyModelId(id int64) bool {
	propertyModel, err := GetPropertyModelByPropertyId(id)
	if err != nil {
		return false
	}
	serviceModel, err := GetServiceModelByServiceModelId(propertyModel.ServiceModelId)
	if err != nil {
		return false
	}

	deviceInstances, err := GetDeviceInstancesByDeviceModelId(serviceModel.DeviceModelId)

	if err != nil || len(deviceInstances) == 0 {
		return false
	}
	for _, deviceInstance := range deviceInstances {
		serviceInstances, err := GetServiceInstanceByDeviceID(deviceInstance.DeviceID)
		if err != nil {
			continue
		}
		if len(serviceInstances) > 0 {
			for _, serviceInstance := range serviceInstances {
				if res := IsExistPropertyInstance(serviceInstance.ID, propertyModel.Name); res {
					return true
				}
			}
		}
	}
	return false
}
