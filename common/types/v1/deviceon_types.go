package v1

// device on amqp with event
type GrpcDeviceOnEventAmqp struct {
	DeviceName string `form:"deviceName" json:"deviceName,omitempty"`
	ModelName  string `form:"modelName" json:"modelName,omitempty"`
	EventName  string `form:"eventName" json:"eventName,omitempty"`
	Time       int64  `form:"time" json:"time,omitempty"`
}

type GrpcDeviceOnPeripheralAmqp struct {
	Version string `form:"version" json:"version,omitempty"`
	Status  *int   `form:"status" json:"status,omitempty"`
	MsgType *int   `form:"msgType" json:"msgType,omitempty"`
}

type GrpcDeviceOnAmqp struct {
	Type                        string `form:"type" json:"type,omitempty"`
	EdgeId                      string `form:"edgeId" json:"edgeId,omitempty"`
	*GrpcDeviceOnEventAmqp      `json:",inline"`
	*GrpcDeviceOnPeripheralAmqp `json:",inline"`
}

// periptheral
type GrpcDeviceOnPeripheralMsg struct {
	EdgeId  string                   `form:"edgeId" json:"edgeId,omitempty"`
	Content []map[string]interface{} `form:"content" json:"content,omitempty"`
}

type DeviceOnPeripheralMsg struct {
	EdgeIds   []string `form:"clientIds" json:"clientIds,omitempty"`
	Usbs      []string `form:"usbs" json:"usbs,omitempty"`
	Printers  []string `form:"printers" json:"printers,omitempty"`
	Batteries []string `form:"batteries" json:"batteries,omitempty"`
	Screens   int64    `form:"screens" json:"screens,omitempty"`
	Version   string   `form:"version" json:"version,omitempty"`
}

const (
	//peripheral type
	DeviceOnPeripheralUsb      string = "usb"
	DeviceOnPeripheralPrinters string = "printer"
	DeviceOnPeripheralBattery  string = "battery"
	DeviceOnPeripheralScreen   string = "screen"

	//peripheral model name
	DeviceOnUsbModelName     string = "usb"
	DeviceOnPrinterModelName string = "printer"
	DeviceOnBatteryModelName string = "battery"
	DeviceOnScreenModelName  string = "screen"

	//device model excel file path
	DeviceOnUsbModelPath     string = "./models/model_usb.xlsx"
	DeviceOnPrinterModelPath string = "./models/model_printer.xlsx"
	DeviceOnBatteryModelPath string = "./models/model_battery.xlsx"
	DeviceOnScreenModelPath  string = "./models/model_screen.xlsx"

	//sub device excel file path
	DeviceOnUsbDeviceInstancePath     string = "./subDevices/device_usb.xlsx"
	DeviceOnPrinterDeviceInstancePath string = "./subDevices/device_printer.xlsx"
	DeviceOnBatteryDeviceInstancePath string = "./subDevices/device_battery.xlsx"
	DeviceOnScreenDeviceInstancePath  string = "./subDevices/device_screen.xlsx"

	//deviceOn peripheral event name
	DeviceOnPrinterStatus string = "printer_status"
	DeviceOnPrinterPaper  string = "printer_paper"
	DeviceOnPrinterInk    string = "printer_ink"
	DeviceOnScreen        string = "screen_count"

	//deviceOn profile_setting_status queue status
	DeviceOnProfileSettingStatusSuccess = 1
	DeviceOnProfileSettingStatusFail    = 0
)
