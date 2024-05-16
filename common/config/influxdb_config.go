package config

import (
	"k8s.io/klog"
	"os"
)

type InfluxDBConfig struct {
	Address  string
	DbName   string
	Username string
	Password string
	Duration string
	Enable   bool
}

func GetInfluxDbConfig() *InfluxDBConfig {

	address := ITHINGS_CONFIG.GetString("db.influx.address")
	dbName := ITHINGS_CONFIG.GetString("db.influx.db_name")
	username := ITHINGS_CONFIG.GetString("db.influx.user")
	password := ITHINGS_CONFIG.GetString("db.influx.passwd")
	duration := ITHINGS_CONFIG.GetString("db.influx.duration")
	enable := ITHINGS_CONFIG.GetBool("db.influx.enable")
	klog.Infof("userName: %s,pwd: %s", username, password)
	if os.Getenv("INFLUX_DB") != "" {
		dbName = os.Getenv("INFLUX_DB")
	}
	if os.Getenv("INFLUX_USER") != "" {
		username = os.Getenv("INFLUX_USER")
	}
	if os.Getenv("INFLUX_PASSWORD") != "" {
		password = os.Getenv("INFLUX_PASSWORD")
	}
	if address == "" || dbName == "" || username == "" {
		return nil
	}
	return &InfluxDBConfig{
		Address:  address,
		DbName:   dbName,
		Username: username,
		Password: password,
		Duration: duration,
		Enable:   enable,
	}
}
