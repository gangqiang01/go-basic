package influx_store

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/v2"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

var (
	xClient *influxDbClient
	once    = &sync.Once{}
)

type influxDbClient struct {
	Client          client.Client
	DbName          string
	RetentionPolicy string
}

type Point struct {
	Tags   map[string]string
	Fields map[string]interface{}
}

func RegisterInfluxDb(dbName, policy string, client client.Client) *influxDbClient {
	if xClient == nil {
		xClient = &influxDbClient{
			DbName:          dbName,
			RetentionPolicy: policy,
			Client:          client,
		}
	}
	return xClient

}

func GetInfluxClient() *influxDbClient {
	return xClient
}

// query
func (influxClient influxDbClient) QueryDB(cmd string) (res []client.Result, err error) {
	if influxClient.Client == nil {
		klog.Errorf("InfluxDb client is not initialized")
		return nil, fmt.Errorf("InfluxDb client is nil")
	}
	klog.Infof("query sql: %s", cmd)
	q := client.Query{
		Command:         cmd,
		Database:        influxClient.DbName,
		RetentionPolicy: influxClient.RetentionPolicy,
	}
	if response, err := influxClient.Client.Query(q); err == nil {
		if response.Error() != nil {
			klog.Errorf("query %s error: %v", cmd, response.Error())
			return res, response.Error()
		}
		res = response.Results
	} else {
		klog.Errorf("query %s error: %v", cmd, err)
		return res, err
	}
	return res, nil
}

func (influxClient influxDbClient) CreateRetentionPolicy(duration string) error {
	//1h（1小时）、1d（1天）、1w（1周）
	createPolicySql := fmt.Sprintf("create retention policy %s on  %s duration %s replication 1 SHARD DURATION %s DEFAULT", influxClient.RetentionPolicy, influxClient.DbName, duration, duration)
	if _, err := influxClient.QueryDB(createPolicySql); err != nil {
		klog.Errorf("create retention policy failed: %v", err)
		alterPolicySql := fmt.Sprintf("alter retention policy %s on %s duration %s replication 1 SHARD DURATION %s DEFAULT", influxClient.RetentionPolicy, influxClient.DbName, duration, duration)
		if _, err := influxClient.QueryDB(alterPolicySql); err != nil {
			klog.Errorf("alter retention policy failed: %v", err)
		}
	}
	return nil
}

// insert
func (influxClient influxDbClient) WritePoint(measurement string, tags map[string]string, fields map[string]interface{}) error {
	if influxClient.Client == nil {
		klog.Errorf("InfluxDb client is not initialized")
		return nil
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        influxClient.DbName,
		Precision:       "ms", //精度，默认ns
		RetentionPolicy: influxClient.RetentionPolicy,
	})
	if err != nil {
		klog.Errorln(err)
		return err
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		klog.Errorln(err)
		return err
	}
	bp.AddPoint(pt)
	err = influxClient.Client.Write(bp)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	//klog.Infof("insert success")
	return nil
}

// batch insert
func (influxClient influxDbClient) WritesPoints(measurement string, points []*Point) error {
	if influxClient.Client == nil {
		klog.Errorf("InfluxDb client is not initialized")
		return nil
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        influxClient.DbName,
		Precision:       "ms", //精度，默认ns
		RetentionPolicy: influxClient.RetentionPolicy,
	})
	if err != nil {
		klog.Errorln(err)
		return err
	}
	for _, point := range points {
		pt, err := client.NewPoint(measurement, point.Tags, point.Fields, time.Now())
		if err != nil {
			klog.Errorln(err)
		}
		bp.AddPoint(pt)
	}

	err = influxClient.Client.Write(bp)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	//klog.Infof("insert success")
	return nil
}

func (influxClient influxDbClient) StopClient() error {
	err := influxClient.Client.Close()
	if err != nil {
		klog.Errorf("Stop influxDb client error: %v", err)
		return err
	}
	xClient = nil
	return nil
}
