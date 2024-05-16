package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
	"time"
)

type RuleLinkageLog struct {
	//ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	ID              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	Trigger         string `gorm:"column:trigger; type:text;" json:"trigger"`
	TriggerDetails  string `gorm:"column:trigger_details; type:text;" json:"triggerDetails"`
	Status          int32  `gorm:"column:status;" json:"status"`
	Error           string `gorm:"column:error; type:text;" json:"error"`
	Action          string `gorm:"column:action; type:text;" json:"action"`
	ActionDetails   string `gorm:"column:action_details; type:text;" json:"actionDetails"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (RuleLinkageLog) TableName() string {
	return "rule_linkage_log"
}

func GetRuleLinkageLog() ([]*RuleLinkageLog, error) {
	var ruleLinkageLogs []*RuleLinkageLog
	err := global.DBAccess.Order("create_time_stamp desc").Find(&ruleLinkageLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ruleLinkageLogs, err
}

func GetRuleLinkageLogById(id string) (*RuleLinkageLog, error) {
	var ruleLinkageLog *RuleLinkageLog
	err := global.DBAccess.Where(&RuleLinkageLog{ID: id}).Find(&ruleLinkageLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ruleLinkageLog, err
}

func GetRuleLinkageLogByPage(page int, limit int) ([]*RuleLinkageLog, error) {
	var ruleLinkageLogs []*RuleLinkageLog
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&ruleLinkageLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ruleLinkageLogs, err
}

func GetRuleLinkageLogCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&RuleLinkageLog{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetRuleLinkageLogByPageAndName(page int, limit int, name string) ([]*RuleLinkageLog, error) {
	var ruleLinkageLog []*RuleLinkageLog
	err := global.DBAccess.Where("name = ?", name).Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&ruleLinkageLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ruleLinkageLog, err
}

func GetRuleLinkageLogCountByName(name string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&RuleLinkageLog{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetRuleLinkageLogByPageAndKeywords(page int, limit int, keywords string) ([]*RuleLinkageLog, error) {
	var ruleLinkageLogs []*RuleLinkageLog
	err := global.DBAccess.Where("trigger LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&ruleLinkageLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ruleLinkageLogs, err
}
func GetRuleLinkageLogCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&RuleLinkageLog{}).Where("trigger LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetRuleLinkageLogByPageAndCondition(page int, limit int, trigger, name string, beginTs *int64, endTs *int64) ([]*RuleLinkageLog, error) {
	var rules []*RuleLinkageLog
	tx := global.DBAccess.Model(&RuleLinkageLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if trigger != "" {
		tx = tx.Where("trigger LIKE ?", "%"+trigger+"%")
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

func GetRuleLinkageLogCountByCondition(trigger, name string, beginTs *int64, endTs *int64) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&RuleLinkageLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if trigger != "" {
		tx = tx.Where("trigger LIKE ?", "%"+trigger+"%")
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

func AddRuleLinkageLog(ruleLinkageLog *RuleLinkageLog) error {
	ruleLinkageLog.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&ruleLinkageLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveRuleLinkageLogStatus(id, error string, status int32) error {
	err := global.DBAccess.Model(&RuleLinkageLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Status": status,
		"Error":  error,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteRuleLinkageLog(id string) error {
	if err := global.DBAccess.Where("id = ?", id).Delete(&RuleLinkageLog{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAllRuleLinkageLog() error {
	if err := global.DBAccess.Exec("DELETE FROM rule_linkage_log").Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func BatchDeleteRuleLinkageLog(ids []string) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&RuleLinkageLog{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
