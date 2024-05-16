package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type CommandInstance struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; type:varchar(64); not null" json:"name"`
	ServiceID       string `gorm:"column:service_id; not null" json:"serviceId"`
	UpdateTimeStamp int64  `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
	AccessConfig    string `gorm:"column:access_config; type:varchar(4096); not null" json:"accessConfig,omitempty"`
}

func (CommandInstance) TableName() string {
	return "command_instance"
}

func NewCommandInstance(name string, serviceID string, accessConfig string) *CommandInstance {
	return &CommandInstance{
		Name:         name,
		ServiceID:    serviceID,
		AccessConfig: accessConfig,
	}
}

func AddCommandInstance(docs []*CommandInstance) error {
	err := global.DBAccess.Create(&docs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}

	return nil
}
func UpdateCommandAccessConfig(serviceID string, name string, accessConfig string) error {
	err := global.DBAccess.Model(&CommandInstance{}).Where("service_id = ? and name = ?", serviceID, name).Update("access_config", accessConfig).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func GetCommandInstanceByServiceIdAndName(serviceID string, name string) (*CommandInstance, error) {
	commandInstance := &CommandInstance{}
	err := global.DBAccess.Model(&CommandInstance{}).Where("service_id = ? and name = ?", serviceID, name).Find(commandInstance).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return commandInstance, err
	}
	return commandInstance, nil
}
func IsExistCommandInstance(serviceID string, name string) bool {
	var count int64
	err := global.DBAccess.Model(&CommandInstance{}).Where(CommandInstance{ServiceID: serviceID, Name: name}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func GetCommandInstanceByServiceId(serviceID string) ([]*CommandInstance, error) {
	var commandInstances []*CommandInstance
	err := global.DBAccess.Model(&CommandInstance{}).Where("service_id = ?", serviceID).Find(&commandInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return commandInstances, err
	}
	return commandInstances, nil
}

func DeleteCommandInstance(name string, serviceID string) error {
	err := global.DBAccess.Where("name = ? and service_id = ?", name, serviceID).Delete(&CommandInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteCommandInstanceByServiceId(serviceID string) error {
	err := global.DBAccess.Where("service_id = ?", serviceID).Delete(&CommandInstance{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
