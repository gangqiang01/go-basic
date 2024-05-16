package v1

type InfluxTwinData struct {
	TwinProperty `json:",inline"`
	Time         string `json:"time"`
	DeviceId     string `json:"deviceId"`
}

type InfluxEventData struct {
	ReportEventMsg `json:"inline"`
	Time           string `json:"time"`
}
