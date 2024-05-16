package model

import (
	"gorm.io/gorm"
)

/*
* RegisterTables create all database tables in this function.
* Notice! you should create tables at here!
 */
func RegisterTables(db *gorm.DB) error {
	err := db.AutoMigrate(
		&DeviceModel{},
		&ServiceModel{},
		&PropertyModel{},
		&EventModel{},
		&CommandModel{},
		&DeviceInstance{},
		&ServiceInstance{},
		&RuleLinkage{},
		&RuleLinkageLog{},
		&PropertyInstance{},
		&EventInstance{},
		&CommandInstance{},
		&ProtocolTypes{},
		&AlertConfig{},
		&AlertLog{},
		&AlertHistory{},
		&EventRuleRelation{},
		&DataForward{},
		&DataForwardLog{},
		&DeviceDataForwardRelation{})

	if err != nil {
		return err
	}

	return nil
}
