package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
	"time"
)

type RuleLinkage struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256); unique" json:"name"`
	Description     string `gorm:"column:description; type:varchar(256);" json:"description"`
	Status          string `gorm:"column:status;" json:"status"`
	Trigger         string `gorm:"column:trigger; type:varchar(1024);" json:"trigger"`
	Filter          string `gorm:"column:filter; type:varchar(1024);" json:"filter"`
	Action          string `gorm:"column:action; type:varchar(1024);" json:"action"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
	DeviceModelName string `gorm:"column:device_model_name" json:"deviceModelName"`
}

func (RuleLinkage) TableName() string {
	return "rule_linkage"
}

func GetRuleLinkage() ([]*RuleLinkage, error) {
	var rules []*RuleLinkage
	err := global.DBAccess.Order("create_time_stamp desc").Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}

func GetRuleLinkageById(id int64) (*RuleLinkage, error) {
	var rule *RuleLinkage
	err := global.DBAccess.Where(&RuleLinkage{ID: id}).First(&rule).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rule, err
}

func GetRuleLinkageByPage(page int, limit int) ([]*RuleLinkage, error) {
	var rules []*RuleLinkage
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}

func GetRuleLinkageCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&RuleLinkage{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetRuleLinkageByPageAndKeywords(page int, limit int, keywords string) ([]*RuleLinkage, error) {
	var rules []*RuleLinkage
	err := global.DBAccess.Where("name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}
func GetRuleLinkageCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&RuleLinkage{}).Where("name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetRuleLinkageByName(name string) (*RuleLinkage, error) {
	rule := &RuleLinkage{}
	err := global.DBAccess.Where("name = ?", name).First(rule).Error
	if err != nil {
		return nil, err
	}

	return rule, err
}
func AddRuleLinkage(ruleLinkage *RuleLinkage) error {
	ruleLinkage.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&ruleLinkage).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveRuleLinkage(id int64, name string, description *string, trigger, filter, action string) error {
	ruleMap := make(map[string]interface{}, 4)
	if name != "" {
		ruleMap["Name"] = name
	}
	if description != nil {
		ruleMap["Description"] = description
	}
	if trigger != "" {
		ruleMap["Trigger"] = trigger
	}
	if filter != "" {
		ruleMap["Filter"] = filter
	}
	if action != "" {
		ruleMap["Action"] = action
	}
	err := global.DBAccess.Model(&RuleLinkage{}).Where("id = ?", id).Updates(ruleMap).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func IsExistRuleLinkageByName(name string) bool {
	var count int64
	err := global.DBAccess.Model(&RuleLinkage{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func SaveRuleLinkageStatus(id int64, status string) error {
	err := global.DBAccess.Model(&RuleLinkage{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteRuleLinkage(id int64) error {
	if err := global.DBAccess.Delete(&RuleLinkage{}, id).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	DeleteEventRuleRelationByRuleId(id)
	return nil

}

func BatchDeleteRuleLinkage(ids []int64) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&RuleLinkage{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	for _, id := range ids {
		DeleteEventRuleRelationByRuleId(id)
	}
	return nil
}
