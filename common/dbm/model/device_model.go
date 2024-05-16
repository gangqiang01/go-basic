package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
	"time"
)

type DeviceModel struct {
	ID              int64             `gorm:"primary_key; auto_increment" json:"id"`
	Name            string            `gorm:"column:name; not null; unique; type:varchar(256);" json:"name"`
	Manufacturer    string            `gorm:"column:manufacturer; type:varchar(256);" json:"manufacturer"`
	Industry        string            `gorm:"column:industry; type:varchar(256);" json:"industry"`
	TagNumber       int64             `gorm:"column:tag_number;" json:"tagNumber"`
	GroupID         string            `gorm:"column:group_id; type:varchar(256);" json:"groupId"`
	DataFormat      string            `gorm:"column:data_format; type:varchar(256);" json:"dataType"`
	Description     string            `gorm:"column:description; type:varchar(256);" json:"description"`
	Creator         string            `gorm:"column:creator; type:varchar(256);" json:"creator"`
	CreateTimeStamp int64             `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64             `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
	DeviceNumber    int64             `gorm:"column:device_number;" json:"deviceNumber"`
	ServiceModels   []*ServiceModel   `gorm:"foreignKey:DeviceModelId"`
	DeviceInstances []*DeviceInstance `gorm:"foreignKey:DeviceModelId"`
}

func (DeviceModel) TableName() string {
	return "device_model"
}
func GetModels() ([]*DeviceModel, error) {
	var models []*DeviceModel
	err := global.DBAccess.Order("update_time_stamp desc").Find(&models).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return models, err
}

func GetModelByPage(page int, limit int) ([]*DeviceModel, error) {
	var models []*DeviceModel
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&models).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return models, err
}

func GetModelByPageAndKeywords(page int, limit int, keywords string) ([]*DeviceModel, error) {
	var deviceModels []*DeviceModel
	err := global.DBAccess.Where("name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("update_time_stamp desc").Limit(limit).Find(&deviceModels).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceModels, err
}
func GetDeviceModelCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceModel{}).Where("name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceModelByName(name string) (*DeviceModel, error) {
	deviceModel := &DeviceModel{}
	err := global.DBAccess.Where("name = ?", name).First(deviceModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return deviceModel, err
}
func IsExistDeviceModelByName(name string) bool {
	var count int64
	err := global.DBAccess.Model(&DeviceModel{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	//klog.Infof("count: %d", count)
	if count > 0 {
		return true
	}
	return false
}
func GetDeviceModelAllInfoByName(name string) (*DeviceModel, error) {
	dm, err := GetDeviceModelByName(name)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	dm.ServiceModels, err = GetServiceModelByDeviceModelId(dm.ID)
	if err != nil {
		return nil, err
	}

	return dm, err
}

func GetDeviceModelAllInfoByID(id int64) (*DeviceModel, error) {
	dm, err := GetDeviceModelById(id)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	dm.ServiceModels, err = GetServiceModelByDeviceModelId(dm.ID)
	if err != nil {
		return nil, err
	}

	return dm, nil
}

func GetDeviceModelCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceModel{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceModelById(id int64) (*DeviceModel, error) {
	var deviceModel *DeviceModel
	err := global.DBAccess.First(&deviceModel, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceModel, err
	}
	return deviceModel, err
}

func AddDeviceModel(deviceModel *DeviceModel) error {
	deviceModel.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&deviceModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func UpdateDeviceModelName(oname string, name string) error {
	err := global.DBAccess.Model(&DeviceModel{}).Where("name = ?", oname).Update("name", name).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveDeviceModel(id int64, deviceModel *DeviceModel) error {
	err := global.DBAccess.Model(&DeviceModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Name":         deviceModel.Name,
		"Manufacturer": deviceModel.Manufacturer,
		"Industry":     deviceModel.Industry,
		"TagNumber":    deviceModel.TagNumber,
		"GroupID":      deviceModel.GroupID,
		"DataFormat":   deviceModel.DataFormat,
		"Description":  deviceModel.Description,
		"Creator":      deviceModel.Creator,
		"DeviceNumber": deviceModel.DeviceNumber,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteDeviceModel(id int64) error {
	var serviceModels []*ServiceModel
	tx := global.DBAccess.Begin()
	serviceErr := global.DBAccess.Where(&ServiceModel{DeviceModelId: id}).Find(&serviceModels).Error
	if serviceErr != nil {
		tx.Rollback()
		klog.Errorf("err: %v", serviceErr)
		return serviceErr
	}

	for _, serviceModel := range serviceModels {
		sid := serviceModel.ID
		if propertyErr := tx.Where(&PropertyModel{ServiceModelId: sid}).Delete(&PropertyModel{}).Error; propertyErr != nil {
			tx.Rollback()
			klog.Errorf("err: %v", propertyErr)
			return propertyErr
		}
		if eventErr := tx.Where(&EventModel{ServiceModelId: sid}).Delete(&EventModel{}).Error; eventErr != nil {
			tx.Rollback()
			klog.Errorf("err: %v", eventErr)
			return eventErr
		}
		if commandErr := tx.Where(&CommandModel{ServiceModelId: sid}).Delete(&CommandModel{}).Error; commandErr != nil {
			tx.Rollback()
			klog.Errorf("err: %v", commandErr)
			return commandErr
		}
	}

	if serviceErr := tx.Where(&ServiceModel{DeviceModelId: id}).Delete(&ServiceModel{}).Error; serviceErr != nil {
		tx.Rollback()
		klog.Errorf("err: %v", serviceErr)
		return serviceErr
	}

	err := tx.Delete(&DeviceModel{}, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		tx.Rollback()
		//global.DBAccess.Rollback()
		return err
	}
	tx.Commit()
	//global.DBAccess.Commit()
	return nil

}

func (dm *DeviceModel) FindServiceModel(name string) *ServiceModel {
	for _, sm := range dm.ServiceModels {
		if sm == nil {
			continue
		}
		if sm.Name == name {
			return sm
		}
	}

	return nil
}
