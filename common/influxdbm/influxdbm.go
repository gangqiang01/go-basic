package influxdbm

import (
	"fmt"
	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/influxdbm/influx_store"
	_ "github.com/influxdata/influxdb1-client"
	"github.com/influxdata/influxdb1-client/v2"
	"k8s.io/klog/v2"
	"time"
)

var (
	policy = "twin_policy"
)

func connInfluxDB(conf *config.InfluxDBConfig) client.Client {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     conf.Address,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  time.Second,
	})
	if err != nil {
		klog.Errorln(err)
		return nil
	}

	return client
}

func InitInfluxDb() error {
	klog.Infof("init influxDB")
	config := config.GetInfluxDbConfig()
	if config == nil {
		klog.Errorf("InfluxDBConfig is missing")
		return fmt.Errorf("InfluxDB config is missing")
	}
	if !config.Enable {
		klog.Infof("InfluxDb is disable")
		return nil
	}
	flag := 0
retry:
	xClient := connInfluxDB(config)
	if xClient == nil && flag < 3 {
		flag++
		time.Sleep(3 * time.Second)
		goto retry
		klog.Errorf("Connect influxDb server fail")
		return fmt.Errorf("Connect influxDb server fail")
	}
	createDbSQL := client.NewQuery(fmt.Sprintf("create database %s", config.DbName), "", "")
	if _, err := xClient.Query(createDbSQL); err != nil {
		klog.Errorf("create database failed: %v", err)
		return err
	}
	influxClient := influx_store.RegisterInfluxDb(config.DbName, policy, xClient)
	influxClient.CreateRetentionPolicy(config.Duration)
	return nil
}
