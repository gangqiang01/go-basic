package model

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
	"time"
)

type DataForward struct {
	ID              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256); unique" json:"name"`
	Description     string `gorm:"column:description; type:varchar(256);" json:"description"`
	Status          string `gorm:"column:status;" json:"status"`
	Source          string `gorm:"column:source; type:varchar(1024);" json:"source"`
	Filter          string `gorm:"column:filter; type:varchar(1024);" json:"filter"`
	Destination     string `gorm:"column: destination; type:varchar(1024);" json:"destination"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (DataForward) TableName() string {
	return "data_forward"
}

func GetDataForward() ([]*DataForward, error) {
	var rules []*DataForward
	err := global.DBAccess.Order("create_time_stamp desc").Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}

func GetDataForwardById(id string) (*DataForward, error) {
	var rule *DataForward
	err := global.DBAccess.Where(&DataForward{ID: id}).First(&rule).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rule, err
}

func GetDataForwardByPage(page int, limit int) ([]*DataForward, error) {
	var rules []*DataForward
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}

func GetDataForwardCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DataForward{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDataForwardByPageAndKeywords(page int, limit int, keywords string) ([]*DataForward, error) {
	var rules []*DataForward
	err := global.DBAccess.Where("name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&rules).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return rules, err
}
func GetDataForwardCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DataForward{}).Where("name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDataForwardByName(name string) (*DataForward, error) {
	rule := &DataForward{}
	err := global.DBAccess.Where("name = ?", name).First(rule).Error
	if err != nil {
		return nil, err
	}

	return rule, err
}
func AddDataForward(ruleLinkage *DataForward) error {
	ruleLinkage.CreateTimeStamp = time.Now().UnixNano() / 1e6
	ruleLinkage.ID = utils.NewUUID()
	err := global.DBAccess.Create(&ruleLinkage).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveDataForward(id, name string, description *string, source, destination string) error {
	dataForwardMap := make(map[string]interface{}, 4)
	if name != "" {
		dataForwardMap["Name"] = name
	}
	if description != nil {
		dataForwardMap["Description"] = description
	}
	if source != "" {
		dataForwardMap["Source"] = source
	}
	if destination != "" {
		dataForwardMap["Destination"] = destination
	}
	err := global.DBAccess.Model(&DataForward{}).Where("id = ?", id).Updates(dataForwardMap).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func IsExistDataForwardByName(name string) bool {
	var count int64
	err := global.DBAccess.Model(&DataForward{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func SaveDataForwardStatus(id, status string) error {
	err := global.DBAccess.Model(&DataForward{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteDataForward(id string) error {
	if err := global.DBAccess.Where("id = ?", id).Delete(&DataForward{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	DeleteDeviceDataForwardRelationByDataForwardId(id)
	return nil

}

func BatchDeleteDataForward(ids []string) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&DataForward{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	for _, id := range ids {
		DeleteDeviceDataForwardRelationByDataForwardId(id)
	}
	return nil
}
