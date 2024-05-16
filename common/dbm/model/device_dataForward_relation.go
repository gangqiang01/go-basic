package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type DeviceDataForwardRelation struct {
	ID            int64  `gorm:"primary_key; auto_increment" json:"id"`
	DeviceId      string `gorm:"column:device_id;" json:"deviceId"`
	DataForwardId string `gorm:"column:data_forward_id;" json:"dataForwardId"`
	DeviceModelId int64  `gotm:"column:device_model_id" json:"deviceModelId"`
}

func (DeviceDataForwardRelation) TableName() string {
	return "device_data_forward_relation"
}

func GetDeviceDataForwardRelation() ([]*DeviceDataForwardRelation, error) {
	var deviceDataForwardRelations []*DeviceDataForwardRelation
	err := global.DBAccess.Find(&deviceDataForwardRelations).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceDataForwardRelations, err
}

func GetDeviceDataForwardRelationByDeviceId(deviceId string) ([]*DeviceDataForwardRelation, error) {
	var deviceDataForwardRelations []*DeviceDataForwardRelation
	err := global.DBAccess.Where("device_id = ?", deviceId).Find(&deviceDataForwardRelations).Error
	if err != nil {
		return nil, err
	}
	return deviceDataForwardRelations, err
}
func GetDeviceDataForwardRelationByDeviceModelId(deviceModelId int64) ([]*DeviceDataForwardRelation, error) {
	var deviceDataForwardRelations []*DeviceDataForwardRelation
	err := global.DBAccess.Where("device_model_id = ?", deviceModelId).Find(&deviceDataForwardRelations).Error
	if err != nil {
		return nil, err
	}
	return deviceDataForwardRelations, err
}
func GetDeviceDataForwardRelationByDataForwardId(dataForwardId string) ([]*DeviceDataForwardRelation, error) {
	var deviceDataForwardRelations []*DeviceDataForwardRelation
	err := global.DBAccess.Where("data_forward_id = ?", dataForwardId).Find(&deviceDataForwardRelations).Error
	if err != nil {
		return nil, err
	}

	return deviceDataForwardRelations, err
}

func IsExistDeviceDataForwardRelationByDeviceIdAndDataForwardId(deviceId, dataForwardId string) (bool, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceDataForwardRelation{}).Where("device_id = ? and data_forward_id = ?", deviceId, dataForwardId).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, err
	}

	return true, err
}

func AddDeviceDataForwardRelation(deviceId, dataForwardId string, modelId int64) error {
	relation := &DeviceDataForwardRelation{
		DeviceId:      deviceId,
		DataForwardId: dataForwardId,
		DeviceModelId: modelId,
	}
	err := global.DBAccess.Create(relation).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteDeviceDataForwardRelationByDeviceId(deviceId string) error {
	err := global.DBAccess.Where("device_id = ?", deviceId).Delete(&DeviceDataForwardRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}

func DeleteDeviceDataForwardRelationByDeviceIdAndDataForwardId(deviceId, dataForwardId string) error {
	err := global.DBAccess.Where("device_id = ? and data_forward_id = ?", deviceId, dataForwardId).Delete(&DeviceDataForwardRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil
}
func DeleteDeviceDataForwardRelationByDeviceModelIdAndDataForwardId(deviceModelId int64, dataForwardId string) error {
	err := global.DBAccess.Where("device_model_id = ? and data_forward_id = ?", deviceModelId, dataForwardId).Delete(&DeviceDataForwardRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil
}
func BatchDeleteDeviceDataForwardRelationByDeviceIds(deviceIds []string) error {
	err := global.DBAccess.Where("device_id in ?", deviceIds).Delete(&DeviceDataForwardRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteDeviceDataForwardRelationByDataForwardId(dataForwardId string) error {
	err := global.DBAccess.Where("data_forward_id = ?", dataForwardId).Delete(&DeviceDataForwardRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil
}
