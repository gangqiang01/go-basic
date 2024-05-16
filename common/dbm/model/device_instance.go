package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
	"time"
)

type DeviceInstance struct {
	DeviceID                 string  `gorm:"column:device_id; type:varchar(36); primary_key;" json:"deviceId"`
	Name                     string  `gorm:"column:name; type:varchar(64); not null" json:"name"`
	EdgeID                   string  `gorm:"column:edge_id; type:varchar(36); not null" json:"edgeID"`
	DeviceOS                 string  `gorm:"column:device_os; type:varchar(36)" json:"deviceOS,omitempty"`
	DeviceCategory           string  `gorm:"column:device_category; type:varchar(64)" json:"deviceCategory,omitempty"`
	DeviceVersion            int     `gorm:"column:device_version; type:int" json:"deviceVersion,omitempty"`
	DeviceIdentificationCode string  `gorm:"column:device_identification_code; type:varchar(64)" json:"deviceIdentificationCode,omitempty"`
	Description              *string `gorm:"column:description; type:varchar(256)" json:"description"`
	GroupName                string  `gorm:"column:group_name; type:varchar(64)" json:"groupName,omitempty"`
	GroupID                  string  `gorm:"column:group_id; type:varchar(64)" json:"groupId,omitempty"`
	Creator                  string  `gorm:"column:creator; type:varchar(64)" json:"creator,omitempty"`
	DeviceAuthType           string  `gorm:"column:device_auth_type; type:varchar(64)" json:"deviceAuthType,omitempty"`
	Secret                   string  `gorm:"column:secret; type:varchar(64)" json:"secret,omitempty"`
	DeviceType               string  `gorm:"column:device_type; type:varchar(64)" json:"deviceType,omitempty"`
	GatewayID                string  `gorm:"column:gateway_id; type:varchar(64)" json:"gatewayId,omitempty"`
	GatewayName              string  `gorm:"column:gateway_name; type:varchar(64)" json:"gatewayName,omitempty"`
	Tags                     string  `gorm:"column:tags; type:text" json:"tags,omitempty"`
	//0: normal, 1: warning, 2: error
	Health                 int64   `gorm:"column:health;" json:"health,omitempty"`
	LifeTimeOfDesiredValue int64   `form:"column:ltodv" json:"ltodv,omitempty"`
	DeviceModelRef         string  `gorm:"column:device_model_ref; type:varchar(64)" json:"deviceModelRef,omitempty"`
	ProtocolType           string  `gorm:"column:protocol_type; type:varchar(36)" json:"protocolType"`
	Protocol               *string `gorm:"column:protocol; type:varchar(4096)" json:"protocol,omitempty"`
	CreateTimeStamp        int64   `gorm:"column:create_time_stamp" json:"createTimeStamp,omitempty"`
	UpdateTimeStamp        int64   `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp,omitempty"`
	//active, inActive, offline, online
	DeviceStatus string `gorm:"column:device_status" json:"deviceStatus,omitempty"`
	//started, stoped
	State            string             `gorm:"column:state" json:"state,omitempty"`
	DeviceModelId    int64              `gorm:"column:device_model_id" json:"deviceModelId"`
	ServiceInstances []*ServiceInstance `gorm:"foreignKey:DeviceID"`
}

func (DeviceInstance) TableName() string {
	return "device_instance"
}

func GetDeviceInstances() ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Order("create_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}

func GetDeviceInstanceByDeviceId(deviceID string) (DeviceInstance, error) {
	var deviceInstance DeviceInstance
	err := global.DBAccess.Where("device_id=?", deviceID).First(&deviceInstance).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstance, err
	}
	return deviceInstance, err
}

func GetDeviceInstance(deviceID string) *DeviceInstance {
	var deviceInstance DeviceInstance

	err := global.DBAccess.Where("device_id=?", deviceID).Find(&deviceInstance).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil
	}

	return &deviceInstance
}

func GetDeviceInstancesByDeviceModelId(deviceModelId int64) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("device_model_id=?", deviceModelId).Order("update_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstances, err
	}
	return deviceInstances, err
}

func GetDeviceInstanceByName(name string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("name=?", name).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstances, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceByEdgeId(edgeId string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("edge_id=?", edgeId).Order("update_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstances, err
	}
	return deviceInstances, err
}

func GetDeviceInstanceByEdgeIdAndDeviceStatus(edgeId, deviceStatus string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("edge_id=? and device_status=?", edgeId, deviceStatus).Order("update_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstances, err
	}
	return deviceInstances, err
}

func GetDeviceInstanceAllInfo(deviceID string) (DeviceInstance, error) {
	deviceInstance, err := GetDeviceInstanceByDeviceId(deviceID)
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstance, err
	}
	serviceInstances, err := GetServiceInstanceByDeviceID(deviceID)
	if err != nil {
		klog.Errorf("err: %v", err)
		return deviceInstance, err
	}
	deviceInstance.ServiceInstances = serviceInstances
	return deviceInstance, nil
}

func GetDeviceInstanceByPage(page int, limit int) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

// get all device instance by edgeID and protocol type.
func GetAllDeviceInstancesV2(edgeID, protocType string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance

	err := global.DBAccess.Where("edge_id = ? AND protocol_type = ?", edgeID, protocType).Order("create_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	//Get all service.
	for _, di := range deviceInstances {
		si, err := GetServiceInstanceByDeviceID(di.DeviceID)
		if err != nil {
			klog.Warningf("GetServiceInstanceByDeviceID with err: %v", err)
			continue
		}

		di.ServiceInstances = si
	}

	return deviceInstances, nil
}

func GetAllProtocolTypeInThisEdge(edgeID string) map[string]string {
	protoMap := make(map[string]string)
	var deviceInstances []*DeviceInstance

	err := global.DBAccess.Where("edge_id = ? ", edgeID).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("GetAllProtocolTypeInThisEdge with err: %v", err)
		return protoMap
	}

	for _, di := range deviceInstances {
		if di == nil {
			continue
		}

		protoMap[di.ProtocolType] = "t"
	}

	return protoMap
}

func GetAllDeviceInstancesV2ByEdgeId(edgeID string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance

	err := global.DBAccess.Where("edge_id = ? ", edgeID).Order("create_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	//Get all service.
	for _, di := range deviceInstances {
		si, err := GetServiceInstanceByDeviceID(di.DeviceID)
		if err != nil {
			klog.Warningf("GetServiceInstanceByDeviceID with err: %v", err)
			continue
		}

		di.ServiceInstances = si
	}

	return deviceInstances, nil
}

func GetDeviceInstanceByPageAndKeywords(page int, limit int, keywords string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDeviceInstanceByPageAndModelId(page int, limit int, modelId int64) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("device_model_id = ?", modelId).Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCountByModelId(modelId int64) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("device_model_id = ?", modelId).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceInstanceByPageAndKeywordsAndModelId(page int, limit int, keywords string, modelId int64) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("device_model_id = ? and name LIKE ?", modelId, "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCountByKeywordsAndModelId(keywords string, modelId int64) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("device_model_id = ? and name LIKE ?", modelId, "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceInstanceByPageAndProtocolType(page int, limit int, protocolType string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("protocol_type = ?", protocolType).Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCountByProtocolType(protocolType string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("protocol_type = ?", protocolType).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceInstanceByPageAndKeywordsAndProtocolType(page int, limit int, keywords, protocolType string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("protocol_type = ? and name LIKE ?", protocolType, "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCountByKeywordsAndProtocolType(keywords, protocolType string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("protocol_type = ? and name LIKE ?", protocolType, "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceInstanceByPageAndEdgeId(page int, limit int, edgeId string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("edge_id = ?", edgeId).Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}

func GetDeviceInstanceByPageAndCondition(page int, limit int, keywords, protocolType, edgeId string, modelId *int64, deviceStatus string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	tx := global.DBAccess.Model(&DeviceInstance{})

	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	if protocolType != "" {
		tx = tx.Where("protocol_type = ?", protocolType)
	}
	if edgeId != "" {
		tx = tx.Where("edge_id = ?", edgeId)
	}
	if modelId != nil {
		tx = tx.Where("device_model_id = ?", *modelId)
	}
	if deviceStatus != "" {
		tx = tx.Where("device_status = ?", deviceStatus)
	}

	err := tx.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}

func GetDeviceInstanceCountByCondition(keywords, protocolType, edgeId string, modelId *int64, deviceStatus string) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&DeviceInstance{})

	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	if protocolType != "" {
		tx = tx.Where("protocol_type = ?", protocolType)
	}
	if edgeId != "" {
		tx = tx.Where("edge_id = ?", edgeId)
	}
	if modelId != nil {
		tx = tx.Where("device_model_id = ?", *modelId)
	}
	if deviceStatus != "" {
		tx = tx.Where("device_status = ?", deviceStatus)
	}
	err := tx.Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDeviceInstanceCountByEdgeId(edgeId string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("edge_id = ?", edgeId).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetDeviceInstanceByPageAndKeywordsAndEdgeId(page int, limit int, keywords, edgeId string) ([]*DeviceInstance, error) {
	var deviceInstances []*DeviceInstance
	err := global.DBAccess.Where("edge_id = ? and name LIKE ?", edgeId, "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&deviceInstances).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return deviceInstances, err
}
func GetDeviceInstanceCountByKeywordsAndEdgeId(keywords, edgeId string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("edge_id = ? and name LIKE ?", edgeId, "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDeviceInstanceCountByStatusAndHealth(deviceStatus string, health int64) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("device_status = ? and health = ?", deviceStatus, health).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetDeviceInstanceCountByStatus(deviceStatus string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("device_status = ?", deviceStatus).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func IsExistDeviceInstanceByNameAndEdgeId(name, edgeId string) bool {
	var count int64
	err := global.DBAccess.Model(&DeviceInstance{}).Where("name = ? and edge_id = ?", name, edgeId).Count(&count).Error
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

func AddDeviceInstance(deviceInstance *DeviceInstance) error {
	deviceInstance.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&deviceInstance).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveDeviceInstance(deviceId string, deviceInstance *DeviceInstance) error {
	err := global.DBAccess.Model(&DeviceInstance{}).Where("device_id = ?", deviceId).Updates(map[string]interface{}{
		"Name": deviceInstance.Name,
	}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

// update
func UpdateDeviceInstanceProtocol(deviceId string, protocol string) error {
	err := global.DBAccess.Model(&DeviceInstance{}).Where(" device_id = ?", deviceId).Update("protocol", protocol).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func UpdateDeviceInstanceHealth(deviceID string, health int64) error {
	err := global.DBAccess.Model(&DeviceInstance{}).Where(" device_id = ?", deviceID).Update("health", health).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

// update all device status in this edge.
func UpdateAllDevInstStatusInThisEdge(edgeID, status string) error {
	err := global.DBAccess.Model(&DeviceInstance{}).Where(" edge_id = ?", edgeID).Update("device_status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateDeviceInstance(deviceID string, doc *DeviceInstance) error {
	vals := make(map[string]interface{})

	vals["update_time_stamp"] = time.Now().UnixNano() / 1e6
	if doc.State != "" {
		vals["state"] = doc.State
	}
	if doc.Protocol != nil {
		vals["protocol"] = doc.Protocol
	}
	if doc.DeviceStatus != "" {
		vals["device_status"] = doc.DeviceStatus
	}

	err := global.DBAccess.Model(&DeviceInstance{}).Where(" device_id = ?", deviceID).Updates(vals).Error
	if err != nil {
		klog.Errorf("UpdateDeviceInstance with err: %v", err)
		return err
	}

	return nil
}

func DeleteDeviceInstance(deviceId string) error {
	var serviceInstances []*ServiceInstance
	tx := global.DBAccess.Begin()
	serviceErr := tx.Where(&ServiceInstance{DeviceID: deviceId}).Find(&serviceInstances).Error
	if serviceErr != nil {
		tx.Rollback()
		klog.Errorf("err: %v", serviceErr)
		return serviceErr
	}

	for _, serviceInstance := range serviceInstances {
		sid := serviceInstance.ID

		eventInstances, err := GetEventInstanceByServiceId(sid)
		if err != nil {
			klog.Errorf("err: %v", err)
		}
		for _, eventInstance := range eventInstances {
			DeleteEventRuleRelationByEventId(eventInstance.ID)
		}
		if propertyErr := tx.Where(&PropertyInstance{ServiceID: sid}).Delete(&PropertyInstance{}).Error; propertyErr != nil {
			tx.Rollback()
			klog.Errorf("err: %v", propertyErr)
			return propertyErr
		}
		if eventErr := tx.Where(&EventInstance{ServiceID: sid}).Delete(&EventInstance{}).Error; eventErr != nil {
			tx.Rollback()
			klog.Errorf("err: %v", eventErr)
			return eventErr
		}
		if commandErr := tx.Where(&CommandInstance{ServiceID: sid}).Delete(&CommandInstance{}).Error; commandErr != nil {
			tx.Rollback()
			klog.Errorf("err: %v", commandErr)
			return commandErr
		}
	}

	if serviceErr := tx.Where(&ServiceInstance{DeviceID: deviceId}).Delete(&ServiceInstance{}).Error; serviceErr != nil {
		tx.Rollback()
		klog.Errorf("err: %v", serviceErr)
		return serviceErr
	}

	if err := tx.Where(&DeviceInstance{DeviceID: deviceId}).Delete(&DeviceInstance{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		tx.Rollback()
		return err
	}

	if err := tx.Where("device_id = ?", deviceId).Delete(&DeviceDataForwardRelation{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		tx.Rollback()
		return err
	}
	if err := tx.Where("device_id = ?", deviceId).Delete(&AlertHistory{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	if err := tx.Where("device_id = ?", deviceId).Delete(&AlertLog{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	tx.Commit()
	return nil
}

func AddDeviceInstanceAll(deviceInstances []*DeviceInstance, services []*ServiceInstance, propertys []*PropertyInstance, events []*EventInstance, commands []*CommandInstance) error {
	tx := global.DBAccess.Begin()
	if len(deviceInstances) > 0 {
		if err := tx.CreateInBatches(&deviceInstances, 1000).Error; err != nil {
			tx.Rollback()
			klog.Errorf("err: %v", err)
			return err
		}
	}

	if len(services) > 0 {
		if err := tx.CreateInBatches(&services, 1000).Error; err != nil {
			tx.Rollback()
			klog.Errorf("err: %v", err)
			return err
		}
	}

	if len(propertys) > 0 {
		if err := tx.CreateInBatches(&propertys, 1000).Error; err != nil {
			tx.Rollback()
			klog.Errorf("err: %v", err)
			return err
		}
	}

	if len(events) > 0 {
		if err := tx.CreateInBatches(&events, 1000).Error; err != nil {
			tx.Rollback()
			klog.Errorf("err: %v", err)
			return err
		}
	}

	if len(commands) > 0 {
		if err := tx.CreateInBatches(&commands, 1000).Error; err != nil {
			tx.Rollback()
			klog.Errorf("err: %v", err)
			return err
		}
	}
	tx.Commit()
	return nil
}
