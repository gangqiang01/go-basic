package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
	"time"
)

type AlertConfig struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256)" json:"name"`
	Description     string `gorm:"column:description; type:varchar(256);" json:"description"`
	Level           int64  `gorm:"column:level;" json:"level"`
	Notification    string `gorm:"column:notification; type:varchar(1024);" json:"notification"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (AlertConfig) TableName() string {
	return "alert_config"
}

func GetAlert() ([]*AlertConfig, error) {
	var alerts []*AlertConfig
	err := global.DBAccess.Order("update_time_stamp desc").Find(&alerts).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alerts, err
}

func GetAlertByPage(page int, limit int) ([]*AlertConfig, error) {
	var alerts []*AlertConfig
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&alerts).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alerts, err
}

func GetAlertCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertConfig{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetAlertByPageAndKeywords(page int, limit int, keywords string) ([]*AlertConfig, error) {
	var alerts []*AlertConfig
	err := global.DBAccess.Where("name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("update_time_stamp desc").Limit(limit).Find(&alerts).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alerts, err
}
func GetAlertCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertConfig{}).Where("name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetAlertByName(name string) (*AlertConfig, error) {
	alert := &AlertConfig{}
	err := global.DBAccess.Where("name = ?", name).First(alert).Error
	if err != nil {
		return nil, err
	}

	return alert, err
}
func GetAlertById(id int64) (*AlertConfig, error) {
	alert := &AlertConfig{}
	err := global.DBAccess.First(alert, id).Error
	if err != nil {
		return nil, err
	}

	return alert, err
}
func AddAlert(alert *AlertConfig) error {
	alert.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&alert).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveAlert(id int64, alert *AlertConfig) error {
	err := global.DBAccess.Model(&AlertConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Name":         alert.Name,
		"Description":  alert.Description,
		"Level":        alert.Level,
		"Notification": alert.Notification,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func IsExistAlertByNameAndLevel(name string, level int64) bool {
	var count int64
	err := global.DBAccess.Debug().Model(&AlertConfig{}).Where("name = ? and level = ?", name, level).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func DeleteAlert(id int64) error {
	err := global.DBAccess.Delete(&AlertConfig{}, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}

func BatchDeleteAlert(ids []int64) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&AlertConfig{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
