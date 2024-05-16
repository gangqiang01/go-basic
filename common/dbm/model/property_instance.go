package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type PropertyInstance struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; type:varchar(64); not null" json:"name"`
	ServiceID       string `gorm:"column:service_id; not null" json:"serviceId"`
	AccessConfig    string `gorm:"column:access_config; type:varchar(4096); not null" json:"accessConfig,omitempty"`
	UpdateTimeStamp int64  `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (PropertyInstance) TableName() string {
	return "property_instance"
}

func NewPropertyInstance(name string, serviceID string, accessConfig string) *PropertyInstance {
	return &PropertyInstance{
		Name:         name,
		ServiceID:    serviceID,
		AccessConfig: accessConfig,
	}
}

func AddPropertyInstance(docs []*PropertyInstance) error {
	err := global.DBAccess.Create(&docs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}

	return nil
}
func UpdatePropertyAccessConfig(serviceID string, name string, accessConfig string) error {
	err := global.DBAccess.Model(&PropertyInstance{}).Where("service_id = ? and name = ?", serviceID, name).Update("access_config", accessConfig).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func GetPropertyInstanceByServiceIdAndName(serviceID string, name string) (*PropertyInstance, error) {
	propertyInstance := &PropertyInstance{}
	err := global.DBAccess.Model(&PropertyInstance{}).Where("service_id = ? and name = ?", serviceID, name).Find(propertyInstance).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return propertyInstance, err
	}
	return propertyInstance, nil
}

func GetPropertyInstanceByServiceId(serviceID string) ([]*PropertyInstance, error) {
	var propertyInstances []*PropertyInstance
	err := global.DBAccess.Model(&PropertyInstance{}).Where("service_id = ?", serviceID).Order("update_time_stamp desc").Find(&propertyInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return propertyInstances, err
	}
	return propertyInstances, nil
}

func IsExistPropertyInstance(serviceID string, name string) bool {
	var count int64
	err := global.DBAccess.Model(&PropertyInstance{}).Where(PropertyInstance{ServiceID: serviceID, Name: name}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true

}

func DeletePropertyInstance(name string, serviceID string) error {
	err := global.DBAccess.Where("name = ? and service_id = ?", name, serviceID).Delete(&PropertyInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeletePropertyInstanceByServiceId(serviceID string) error {
	err := global.DBAccess.Where("service_id = ?", serviceID).Delete(&PropertyInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
