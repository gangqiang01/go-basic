package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type CommandModel struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	RequestParam    string `gorm:"column:request_param; type:varchar(256);" json:"RequestParam"`
	ResponseParam   string `gorm:"column:response_param; type:varchar(256);" json:"ResponseParam"`
	Description     string `gorm:"column:description; type:varchar(256);" json:"description"`
	UpdateTimeStamp int64  `gorm:"autoUpdateTime:milli" json:"updateTimeStamp"`
	ServiceModelId  int64  `gorm:"column:service_model_id" json:"serviceModelId,omitempty"`
}

func (CommandModel) TableName() string {
	return "command_model"
}

func GetCommandModelByServiceModelId(serviceId int64) ([]*CommandModel, error) {
	var commandModel []*CommandModel
	err := global.DBAccess.Where("service_model_id=?", serviceId).Find(&commandModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return commandModel, err
}
func IsExistCommandModel(serviceModelId int64, commandName string) bool {
	var count int64
	err := global.DBAccess.Model(&CommandModel{}).Where(CommandModel{ServiceModelId: serviceModelId, Name: commandName}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}
func AddCommandModel(commandModel *CommandModel) error {
	err := global.DBAccess.Create(&commandModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveCommandModel(id int64, commandModel *CommandModel) error {
	err := global.DBAccess.Model(&CommandModel{}).Where("id = ?", id).Save(&commandModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteCommandModel(id int64) error {
	//global.DBAccess.Begin()
	err := global.DBAccess.Delete(&CommandModel{}, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		//global.DBAccess.Rollback()
		return err
	}
	//global.DBAccess.Commit()
	return nil
}
