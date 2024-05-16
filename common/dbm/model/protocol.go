package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type ProtocolTypes struct {
	ID           int64  `gorm:"primary_key; auto_increment" json:"id,omitempty"`
	ProtocolType string `gorm:"column:protocol_type; type:varchar(64); unique; not null" json:"type,omitempty"`
	Spec         string `gorm:"column:spec; type:varchar(36)" json:"spec,omitempty"`
}

func (*ProtocolTypes) TableName() string {
	return "protocol_types"
}

func GetProtocolType(protocolType string) (*ProtocolTypes, error) {
	var protoType = ProtocolTypes{}

	err := global.DBAccess.Where("protocol_type=?", protocolType).First(&protoType).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return &protoType, err
}

func GetAllProtocolType() ([]*ProtocolTypes, error) {
	var protoTypes []*ProtocolTypes
	err := global.DBAccess.Find(&protoTypes).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return protoTypes, err
}

func AddProtocolType(doc *ProtocolTypes) error {
	if doc == nil || doc.ProtocolType == "" {
		return global.ErrInvalidParms
	}

	protoType, _ := GetProtocolType(doc.ProtocolType)
	if protoType != nil && protoType.ProtocolType != "" {
		return nil
	}

	return global.DBAccess.Create(doc).Error
}

func DeleteProtocolType(protocolType string) error {
	protoType, _ := GetProtocolType(protocolType)
	if protoType == nil || protoType.ProtocolType == "" {
		return nil
	}

	return global.DBAccess.Delete(&ProtocolTypes{}, protocolType).Error
}
