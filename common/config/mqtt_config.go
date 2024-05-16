package config

import (
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
	"os"
)

// MqttConfig
type MqttConfig struct {
	Broker             string
	TSLEnable          bool
	InsecureSkipVerify bool

	User, Passwd string
	ClientID     string
	// QOS indicates mqtt qos
	// 0: QOSAtMostOnce, 1: QOSAtLeastOnce, 2: QOSExactlyOnce
	// default 0
	// Note: Can not use "omitempty" option,  It will affect the output of the default configuration file
	QOS          int
	CaFilePath   string
	CertFilePath string
	KeyFilePath  string
	MaxGoRoutine int
}

func GetMqttConfig() *MqttConfig {

	transConfig := &MqttConfig{}

	broker := ITHINGS_CONFIG.GetString("transport.mqtt.broker")
	if broker == "" {
		klog.Infof("Broker url for mqtt client is nil")
		return nil
	}
	transConfig.Broker = broker

	transConfig.User = ITHINGS_CONFIG.GetString("transport.mqtt.usr")
	transConfig.Passwd = ITHINGS_CONFIG.GetString("transport.mqtt.passwd")
	transConfig.TSLEnable = ITHINGS_CONFIG.GetBool("transport.mqtt.ssl")
	transConfig.InsecureSkipVerify = ITHINGS_CONFIG.GetBool("transport.mqtt.insecure_skip_verify")

	if os.Getenv("MQTT_USERNAME") != "" {
		transConfig.User = os.Getenv("MQTT_USERNAME")
	}
	if os.Getenv("MQTT_PASSWORD") != "" {
		transConfig.Passwd = os.Getenv("MQTT_PASSWORD")
	}
	transConfig.ClientID = ITHINGS_CONFIG.GetString("transport.mqtt.clientid")
	if transConfig.ClientID == "" {
		macs := utils.GetLocalMACs()
		if macs == nil || len(macs) == 0 {
			//maybe, we will replace it as a dummy.
			klog.Warningf("no mac address, we generate a random value.")
			transConfig.ClientID = utils.NewUUID()[:12]
		} else {
			transConfig.ClientID = macs[0] + "ithings"
		}
	}

	//Read Qos
	transConfig.QOS = ITHINGS_CONFIG.GetInt("transport.mqtt.qos")
	transConfig.CaFilePath = ITHINGS_CONFIG.GetString("transport.mqtt.cafile")
	transConfig.CertFilePath = ITHINGS_CONFIG.GetString("transport.mqtt.certfile")
	transConfig.KeyFilePath = ITHINGS_CONFIG.GetString("transport.mqtt.keyfile")

	if os.Getenv("MQTT_TSL") != "" && os.Getenv("MQTT_TSL") == "true" {

		transConfig.TSLEnable = true
	} else {
		transConfig.TSLEnable = false
	}

	if os.Getenv("MQTT_CA_FILE") != "" {
		transConfig.CaFilePath = os.Getenv("MQTT_CA_FILE")
	}
	if os.Getenv("MQTT_CRT_FILE") != "" {
		transConfig.CertFilePath = os.Getenv("MQTT_CRT_FILE")
	}
	if os.Getenv("MQTT_KEY_FILE") != "" {
		transConfig.KeyFilePath = os.Getenv("MQTT_KEY_FILE")
	}
	if os.Getenv("MQTT_INSECURE_VERIFY") != "" && os.Getenv("MQTT_INSECURE_VERIFY") == "true" {

		transConfig.InsecureSkipVerify = true
	} else {
		transConfig.InsecureSkipVerify = false
	}

	transConfig.MaxGoRoutine = ITHINGS_CONFIG.GetInt("transport.mqtt.max-go-routines")
	if transConfig.MaxGoRoutine < 512 {
		transConfig.MaxGoRoutine = 1024
	}

	return transConfig
}
