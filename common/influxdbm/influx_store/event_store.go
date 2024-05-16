package influx_store

import (
	"encoding/json"
	"fmt"
	"github.com/edgehook/ithings/common/types/v1"
	"k8s.io/klog/v2"
	"strings"
)

var (
	eventMeasurement = "report_event"
)

func StoreEvent(eventMsg *v1.ReportEventMsg) error {
	var points []*Point
	tags := map[string]string{"service": eventMsg.ServiceName, "event": eventMsg.EventName, "deviceId": eventMsg.DeviceID}
	fields := map[string]interface{}{"details": eventMsg.Details, "errMsg": eventMsg.ErrorMessage, "ts": eventMsg.Timestamp}
	point := &Point{
		Tags:   tags,
		Fields: fields,
	}
	points = append(points, point)

	xClient := GetInfluxClient()
	if xClient == nil {
		return nil
	}
	err := xClient.WritesPoints(eventMeasurement, points)
	if err != nil {
		return err
	}
	return nil
}

// dashboard simpleJson
func QueryTableEvent(deviceId, service, event, startTs, endTs string, count *int64) []*v1.InfluxEventData {

	var eventTables []*v1.InfluxEventData
	sql := fmt.Sprintf("%s %s", "select time, deviceId, service, event, details, ts, errMsg from", eventMeasurement)
	if deviceId != "" || service != "" || event != "" || startTs != "" || endTs != "" {
		sql = fmt.Sprintf("%s where", sql)
		var sqlList []string
		if deviceId != "" {
			sqlList = append(sqlList, fmt.Sprintf("deviceId = '%s'", deviceId))
		}
		if service != "" {
			sqlList = append(sqlList, fmt.Sprintf("service = '%s'", service))
		}

		if event != "" {
			sqlList = append(sqlList, fmt.Sprintf("event = '%s'", service))
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
	xClient := GetInfluxClient()
	if xClient == nil {
		return nil
	}
	responces, err := xClient.QueryDB(sql)
	if err != nil {
		klog.Errorf("Query db error: %v", err)
		return nil
	}
	klog.Infof("reponce: %v", responces)
	if len(responces) > 0 && len(responces[0].Series) > 0 {
		for _, row := range responces[0].Series[0].Values {
			for j, value := range row {
				klog.Infof("j:%d value:%v", j, value)

			}
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
				klog.Errorf("convert eventName to string error: %v", row[2])
			}
			details, ok := row[4].(string)
			if !ok {
				klog.Errorf("convert details to string error: %v", row[2])
			}
			ts, ok := row[5].(json.Number)
			if !ok {
				klog.Errorf("convert ts to json.Number error: %v", row[5])
			}
			temp, err := ts.Int64()
			if err != nil {
				klog.Errorf("convert ts to int64 error: %v", row[5])
			}

			errorMsg, ok := row[6].(string)
			if !ok {
				klog.Errorf("convert errorMsg to string error: %v", row[6])
			}
			eventMsg := v1.ReportEventMsg{
				ServiceName:  service,
				EventName:    propertyName,
				DeviceID:     deviceId,
				Timestamp:    temp,
				Details:      details,
				ErrorMessage: errorMsg,
			}
			influxEvent := &v1.InfluxEventData{
				ReportEventMsg: eventMsg,
				Time:           time,
			}
			eventTables = append(eventTables, influxEvent)
		}
	}
	return eventTables
}

func QueryEvent(deviceId, service, event string, startTs, endTs, count *int64) []*v1.ReportEventMsg {

	var eventTables []*v1.ReportEventMsg
	sql := fmt.Sprintf("%s %s", "select deviceId, service, event, details, ts, errMsg from", eventMeasurement)
	if deviceId != "" || service != "" || event != "" || startTs != nil || endTs != nil {
		sql = fmt.Sprintf("%s where", sql)
		var sqlList []string
		if deviceId != "" {
			sqlList = append(sqlList, fmt.Sprintf("deviceId = '%s'", deviceId))
		}
		if service != "" {
			sqlList = append(sqlList, fmt.Sprintf("service = '%s'", service))
		}

		if event != "" {
			sqlList = append(sqlList, fmt.Sprintf("event = '%s'", service))
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
	xClient := GetInfluxClient()
	if xClient == nil {
		return nil
	}
	responces, err := xClient.QueryDB(sql)
	if err != nil {
		klog.Errorf("Query db error: %v", err)
		return nil
	}
	klog.Infof("reponce: %v", responces)
	if len(responces) > 0 && len(responces[0].Series) > 0 {
		for _, row := range responces[0].Series[0].Values {
			for j, value := range row {
				klog.Infof("j:%d value:%v", j, value)

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
				klog.Errorf("convert eventName to string error: %v", row[2])
			}
			details, ok := row[4].(string)
			if !ok {
				klog.Errorf("convert details to string error: %v", row[2])
			}
			ts, ok := row[5].(json.Number)
			if !ok {
				klog.Errorf("convert ts to json.Number error: %v", row[5])
			}
			temp, err := ts.Int64()
			if err != nil {
				klog.Errorf("convert ts to int64 error: %v", row[5])
			}

			errorMsg, ok := row[6].(string)
			if !ok {
				klog.Errorf("convert errorMsg to string error: %v", row[6])
			}
			eventMsg := &v1.ReportEventMsg{
				ServiceName:  service,
				EventName:    propertyName,
				DeviceID:     deviceId,
				Timestamp:    temp,
				Details:      details,
				ErrorMessage: errorMsg,
			}
			eventTables = append(eventTables, eventMsg)
		}
	}
	return eventTables
}
