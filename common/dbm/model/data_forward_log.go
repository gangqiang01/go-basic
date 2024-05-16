package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
	"time"
)

type DataForwardLog struct {
	//ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	ID              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	Source          string `gorm:"column:source; type:text;" json:"source"`
	SourceDetails   string `gorm:"column:source_details; type:text;" json:"sourceDetails"`
	Status          int32  `gorm:"column:status;" json:"status"`
	Error           string `gorm:"column:error; type:text;" json:"error"`
	Way             string `gorm:"column:way; type:text;" json:"way"`
	WayDetails      string `gorm:"column:way_details; type:text;" json:"wayDetails"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (DataForwardLog) TableName() string {
	return "data_forward_log"
}

func GetDataForwardLog() ([]*DataForwardLog, error) {
	var dataForwardLogs []*DataForwardLog
	err := global.DBAccess.Order("create_time_stamp desc").Find(&dataForwardLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return dataForwardLogs, err
}

func GetDataForwardLogById(id string) (*DataForwardLog, error) {
	var dataForwardLog *DataForwardLog
	err := global.DBAccess.Where(&DataForwardLog{ID: id}).Find(&dataForwardLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return dataForwardLog, err
}

func GetDataForwardLogByPage(page int, limit int) ([]*DataForwardLog, error) {
	var dataForwardLogs []*DataForwardLog
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&dataForwardLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return dataForwardLogs, err
}

func GetDataForwardLogCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DataForwardLog{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDataForwardLogByPageAndName(page int, limit int, name string) ([]*DataForwardLog, error) {
	var dataForwardLog []*DataForwardLog
	err := global.DBAccess.Where("name = ?", name).Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&dataForwardLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return dataForwardLog, err
}

func GetDataForwardLogCountByName(name string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DataForwardLog{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDataForwardLogByPageAndKeywords(page int, limit int, keywords string) ([]*DataForwardLog, error) {
	var dataForwardLogs []*DataForwardLog
	err := global.DBAccess.Where("source LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&dataForwardLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return dataForwardLogs, err
}
func GetDataForwardLogCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DataForwardLog{}).Where("source LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDataForwardLogByPageAndCondition(page int, limit int, trigger, name string, beginTs *int64, endTs *int64) ([]*DataForwardLog, error) {
	var rules []*DataForwardLog
	tx := global.DBAccess.Model(&DataForwardLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if trigger != "" {
		tx = tx.Where("source LIKE ?", "%"+trigger+"%")
	}

	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}

func GetDataForwardLogCountByCondition(trigger, name string, beginTs *int64, endTs *int64) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&DataForwardLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if trigger != "" {
		tx = tx.Where("source LIKE ?", "%"+trigger+"%")
	}

	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func AddDataForwardLog(dataForwardLog *DataForwardLog) error {
	dataForwardLog.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&dataForwardLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveDataForwardLogStatus(id, error string, status int32) error {
	err := global.DBAccess.Model(&DataForwardLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Status": status,
		"Error":  error,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteDataForwardLog(id string) error {
	if err := global.DBAccess.Where("id = ?", id).Delete(&DataForwardLog{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAllDataForwardLog() error {
	if err := global.DBAccess.Exec("DELETE FROM data_forward_log").Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func BatchDeleteDataForwardLog(ids []string) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&DataForwardLog{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
