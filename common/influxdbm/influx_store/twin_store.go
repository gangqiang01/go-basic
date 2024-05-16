package influx_store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
	"strings"
)

var (
	twinMeasurement = "report_twin"
)

func StoreTwin(deviceID string, twinProperties []*v1.TwinProperty) error {
	xClient := GetInfluxClient()
	if xClient == nil {
		//klog.Error("InfluxClient is nil")
		return nil
	}
	var points []*Point
	for _, twinProperty := range twinProperties {
		if twinProperty.Value == nil {
			klog.Errorf("value is nil: %s", twinProperty.PropertyName)
			continue
		}
		val := utils.ToString(twinProperty.Value)
		tags := map[string]string{"service": twinProperty.Service, "property": twinProperty.PropertyName, "deviceId": deviceID}
		fields := map[string]interface{}{"value": val, "ts": twinProperty.Timestamp}
		point := &Point{
			Tags:   tags,
			Fields: fields,
		}
		points = append(points, point)
	}
	err := xClient.WritesPoints(twinMeasurement, points)
	if err != nil {
		return err
	}
	return nil
}

// dashboard simpleJson
func QueryTableTwin(deviceId, serviceName, propertyName, startTs, endTs string, count *int64) []*v1.InfluxTwinData {
	xClient := GetInfluxClient()
	if xClient == nil {
		return nil
	}
	var twinPropertys []*v1.InfluxTwinData
	sql := fmt.Sprintf("%s %s", "select deviceId, service, property, value, ts, errMsg from", twinMeasurement)
	if deviceId != "" || serviceName != "" || propertyName != "" || startTs != "" || endTs != "" {
		sql = fmt.Sprintf("%s where", sql)
		var sqlList []string
		if deviceId != "" {
			sqlList = append(sqlList, fmt.Sprintf("deviceId = '%s'", deviceId))
		}
		if serviceName != "" {
			sqlList = append(sqlList, fmt.Sprintf("service = '%s'", serviceName))
		}

		if propertyName != "" {
			sqlList = append(sqlList, fmt.Sprintf("property = '%s'", propertyName))
		}
		//utc
		if startTs != "" {
			sqlList = append(sqlList, fmt.Sprintf("time >= '%v'", startTs))
		}
		//utc
		if endTs != "" {
			sqlList = append(sqlList, fmt.Sprintf("time <= '%v'", endTs))
		}

		whereSql := strings.Join(sqlList, " AND ")
		sql = fmt.Sprintf("%s %s", sql, whereSql)

	}
	if count != nil {
		sql = fmt.Sprintf("%s limit %d", sql, *count)
	}

	klog.Infof("query sql: %s", sql)
	responces, _ := xClient.QueryDB(sql)
	//klog.Infof("reponce: %v", responces)
	if len(responces) > 0 && len(responces[0].Series) > 0 {
		for _, row := range responces[0].Series[0].Values {
			time, ok := row[0].(string)
			if !ok {
				klog.Errorf("convert time to string error: %v", row[0])
			}
			deviceId, ok := row[1].(string)
			if !ok {
				klog.Errorf("convert deviceId to string error: %v", row[1])
			}
			service, ok := row[2].(string)
			if !ok {
				klog.Errorf("convert service to string error: %v", row[2])
			}
			propertyName, ok := row[3].(string)
			if !ok {
				klog.Errorf("convert propertyName to string error: %v", row[2])
			}
			ts, ok := row[5].(json.Number)
			if !ok {
				klog.Errorf("convert ts to json.Number error: %v", row[5])
			}
			temp, err := ts.Int64()
			if err != nil {
				klog.Errorf("convert ts to int64 error: %v", row[5])
			}
			twinProperty := v1.TwinProperty{
				Service:      service,
				PropertyName: propertyName,
				Value:        row[4],
				Timestamp:    temp,
			}
			twinData := &v1.InfluxTwinData{
				DeviceId:     deviceId,
				TwinProperty: twinProperty,
				Time:         time,
			}
			twinPropertys = append(twinPropertys, twinData)
		}
	}
	return twinPropertys
}

func QueryTwin(deviceId, serviceName, propertyName string, startTs, endTs, count *int64) []*v1.TwinProperty {
	xClient := GetInfluxClient()
	if xClient == nil {
		return nil
	}
	var twinPropertys []*v1.TwinProperty
	sql := fmt.Sprintf("%s %s", "select deviceId, service, property, value, ts, errMsg from", twinMeasurement)
	if deviceId != "" || serviceName != "" || propertyName != "" || startTs != nil || endTs != nil {
		sql = fmt.Sprintf("%s where", sql)
		var sqlList []string
		if deviceId != "" {
			sqlList = append(sqlList, fmt.Sprintf("deviceId = '%s'", deviceId))
		}
		if serviceName != "" {
			sqlList = append(sqlList, fmt.Sprintf("service = '%s'", serviceName))
		}

		if propertyName != "" {
			sqlList = append(sqlList, fmt.Sprintf("property = '%s'", propertyName))
		}
		//utc
		if startTs != nil {
			sqlList = append(sqlList, fmt.Sprintf("ts >= %v", *startTs))
		}
		//utc
		if endTs != nil {
			sqlList = append(sqlList, fmt.Sprintf("ts <= %v", *endTs))
		}

		whereSql := strings.Join(sqlList, " AND ")
		sql = fmt.Sprintf("%s %s", sql, whereSql)

	}
	if count != nil {
		sql = fmt.Sprintf("%s limit %d", sql, *count)
	}

	klog.Infof("query sql: %s", sql)
	responces, _ := xClient.QueryDB(sql)
	//klog.Infof("reponce: %v", responces)
	if len(responces) > 0 && len(responces[0].Series) > 0 {
		for _, row := range responces[0].Series[0].Values {
			service, ok := row[2].(string)
			if !ok {
				klog.Errorf("convert service to string error: %v", row[2])
			}
			propertyName, ok := row[3].(string)
			if !ok {
				klog.Errorf("convert propertyName to string error: %v", row[2])
			}
			ts, ok := row[5].(json.Number)
			if !ok {
				klog.Errorf("convert ts to json.Number error: %v", row[5])
			}
			temp, err := ts.Int64()
			if err != nil {
				klog.Errorf("convert ts to int64 error: %v", row[5])
			}
			twinProperty := &v1.TwinProperty{
				Service:      service,
				PropertyName: propertyName,
				Value:        row[4],
				Timestamp:    temp,
			}
			twinPropertys = append(twinPropertys, twinProperty)
		}
	}
	return twinPropertys
}

func DeleteTwin(deviceId, serviceName, propertyName string) error {
	xClient := GetInfluxClient()
	if xClient == nil {
		return errors.New("Influx client is nul")
	}
	sql := fmt.Sprintf("delete from %s ", twinMeasurement)
	if deviceId != "" || serviceName != "" || propertyName != "" {
		sql = fmt.Sprintf("%s where", sql)
		var sqlList []string
		if deviceId != "" {
			sqlList = append(sqlList, fmt.Sprintf("deviceId = '%s'", deviceId))
		}
		if serviceName != "" {
			sqlList = append(sqlList, fmt.Sprintf("service = '%s'", serviceName))
		}

		if propertyName != "" {
			sqlList = append(sqlList, fmt.Sprintf("property = '%s'", propertyName))
		}

		whereSql := strings.Join(sqlList, " AND ")
		sql = fmt.Sprintf("%s %s", sql, whereSql)

	}

	klog.Infof("query sql: %s", sql)
	_, err := xClient.QueryDB(sql)
	return err
}
