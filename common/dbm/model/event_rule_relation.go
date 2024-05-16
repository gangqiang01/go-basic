package model

import (
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type EventRuleRelation struct {
	ID           int64 `gorm:"primary_key; auto_increment" json:"id"`
	RuleId       int64 `gorm:"column:rule_id;" json:"ruleId"`
	EventId      int64 `gorm:"column:event_id;" json:"eventId"`
	EventModelId int64 `gorm:"column:event_model_id;" json:"eventModelId"`
}

func (EventRuleRelation) TableName() string {
	return "event_rule_relation"
}

func GetEventRuleRelation() ([]*EventRuleRelation, error) {
	var eventRuleRelations []*EventRuleRelation
	err := global.DBAccess.Find(&eventRuleRelations).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return eventRuleRelations, err
}

func GetEventRuleRelationByEventId(eventId int64) ([]*EventRuleRelation, error) {
	var eventRuleRelations []*EventRuleRelation
	err := global.DBAccess.Where("event_id = ?", eventId).Find(&eventRuleRelations).Error
	if err != nil {
		return nil, err
	}

	return eventRuleRelations, err
}

func GetEventRuleRelationByEventModelId(eventId int64) ([]*EventRuleRelation, error) {
	var eventRuleRelations []*EventRuleRelation
	err := global.DBAccess.Where("event_model_id = ?", eventId).Find(&eventRuleRelations).Error
	if err != nil {
		return nil, err
	}

	return eventRuleRelations, err
}

func GetEventRuleRelationByRuleId(ruleId int64) ([]*EventRuleRelation, error) {
	var eventRuleRelations []*EventRuleRelation
	err := global.DBAccess.Where("rule_id = ?", ruleId).Find(&eventRuleRelations).Error
	if err != nil {
		return nil, err
	}

	return eventRuleRelations, err
}

func IsExistEventRuleRelationByRuleIdAndEventId(ruleId, eventId int64) (bool, error) {
	var count int64
	err := global.DBAccess.Model(&EventRuleRelation{}).Where("rule_id = ? and event_id = ?", ruleId, eventId).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, err
	}

	return true, err
}

func AddEventRuleRelation(ruleId, eventId, modelEventId int64) error {
	eventRuleRelation := &EventRuleRelation{
		RuleId:       ruleId,
		EventId:      eventId,
		EventModelId: modelEventId,
	}
	err := global.DBAccess.Create(eventRuleRelation).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteEventRuleRelationByRuleId(ruleId int64) error {
	err := global.DBAccess.Where("rule_id = ?", ruleId).Delete(&EventRuleRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}

func DeleteEventRuleRelationByRuleIdAndEventId(ruleId, eventId int64) error {
	err := global.DBAccess.Where("rule_id = ? and event_id = ?", ruleId, eventId).Delete(&EventRuleRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}
func DeleteEventRuleRelationByRuleIdAndEventModelId(ruleId, modelEventId int64) error {
	err := global.DBAccess.Where("rule_id = ? and event_model_id = ?", ruleId, modelEventId).Delete(&EventRuleRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}

func BatchDeleteEventRuleRelationByRuleIds(ruleIds []int64) error {
	err := global.DBAccess.Where("rule_id in ?", ruleIds).Delete(&EventRuleRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteEventRuleRelationByEventId(eventId int64) error {
	err := global.DBAccess.Where("event_id = ?", eventId).Delete(&EventRuleRelation{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}
