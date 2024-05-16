package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"k8s.io/klog/v2"
	"reflect"
)

type excelHandler struct {
	Sheet string
}
type sheet string

const (
	CONFIGEXTENSION sheet = "configExtension"
	DEVICEINSTANCE  sheet = "subDevice"

	DEVICEMODEL sheet = "deviceModel"
	PROPERTY    sheet = "property"
	EVENT       sheet = "event"
)

func NewExecHandler(sheet string) *excelHandler {
	return &excelHandler{
		Sheet: sheet,
	}
}

type ExtensionConfigExcel struct {
	IType        string
	Service      string
	Name         string
	AccessConfig string
}

func GetExtensionConfigTitle() []string {
	return []string{"type", "service", "name", "accessConfig"}
}

type DeviceInstanceExcel struct {
	DeviceID string
	Name     string
	EdgeID   string
	//DeviceOS                 string
	//DeviceCategory           string
	//DeviceVersion            int
	//DeviceIdentificationCode string
	Description string
	//GroupName                string
	//GroupID                  string
	//Creator                  string
	//DeviceAuthType           string
	//Secret                   string
	DeviceType string
	//GatewayID                string
	//GatewayName              string
	//Tags                     string
	LifeTimeOfDesiredValue int64
	DeviceModelRef         string
	ProtocolType           string
	Protocol               string
}

func GetDeviceInstanceExcelTitle() []string {
	return []string{"deviceId", "name", "edgeId", "description", "deviceType", "liftTimeOfDesiredValue", "deviceModel", "protocolType", "protocol"}
}

// "name", "description", "manufacturer", "industry", "dataFormat"
type DeviceModelExcel struct {
	Name         string
	Description  string
	Manufacturer string
	Industry     string
	//TagNumber       int64
	//GroupID         string
	DataFormat string

	//Creator         string
}

func GetDeviceModelExcelTitle() []string {
	return []string{"name", "description", "manufacturer", "industry", "dataFormat"}
}

// "service", "name", "description", "dataType", "maxValue", "minValue", "writeAble", "unit"
type PropertyModelExcel struct {
	Service     string
	Name        string
	Description string
	DataType    string
	MaxValue    float64
	MinValue    float64
	WriteAble   string
	Unit        string
}

func GetPropertyModelTitle() []string {
	return []string{"service", "name", "description", "dataType", "maxValue", "minValue", "writeAble", "unit"}
}

type EventModelExcel struct {
	Service     string
	Name        string
	Description string
}

func GetEventModelTitle() []string {
	return []string{"service", "name", "description"}
}

// slice is a struct not pointer
func (excelHandler *excelHandler) CreateExcel(title []string, slice []interface{}, file *excelize.File) error {
	index := file.NewSheet("Sheet1")
	file.SetSheetName("Sheet1", excelHandler.Sheet)
	if err := file.SetSheetRow(excelHandler.Sheet, "A1", &title); err != nil {
		klog.Errorf("Create excel file error: %s", err.Error())
		return err
	}

	rowNum := 1
	for _, value := range slice {
		t := reflect.TypeOf(value)
		value := reflect.ValueOf(value)
		row := make([]interface{}, 0)
		for l := 0; l < t.NumField(); l++ {
			val := value.Field(l).Interface()
			row = append(row, val)
		}

		rowNum++
		if err := file.SetSheetRow(excelHandler.Sheet, fmt.Sprintf("%s%v", "A", rowNum), &row); err != nil {
			return err
		}
	}

	file.SetActiveSheet(index)
	return nil
}

func (excelHandler *excelHandler) CreateExcelByMap(title []string, slice []map[string]string, file *excelize.File) error {
	index := file.NewSheet(excelHandler.Sheet)
	if err := file.SetSheetRow(excelHandler.Sheet, "A1", &title); err != nil {
		klog.Errorf("Create excel file error: %s", err.Error())
		return err
	}
	rowNum := 1
	for _, value := range slice {
		row := make([]interface{}, 0)
		for _, v := range value {
			row = append(row, v)
		}
		rowNum++
		if err := file.SetSheetRow(excelHandler.Sheet, fmt.Sprintf("%s%v", "A", rowNum), &row); err != nil {
			return err
		}
	}

	file.SetActiveSheet(index)
	return nil
}
