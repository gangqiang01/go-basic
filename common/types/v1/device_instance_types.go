package v1

import (
	db "github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

// DeviceTwin provides a logical representation of A series of control properties on a device(writable
// properties in the device model). These properties can have a Desired state and a Reported state.
// The cloud configures the `Desired`state of a device's properties and this configuration update is pushed
// to the edge node. The mapper sends a command to the device to change these properties value as per
// the desired state . It receives the `Reported` state of the property once the previous operation is
// complete and sends the reported state to the cloud.
// Offline device interaction in the edge is possible via twin properties for control/command operations.
type DeviceTwin struct {
	// Required: the desired property list.
	Desired []*TwinProperty `json:"desired,omitempty"`
	// Required: the reported property list.
	Reported []*TwinProperty `json:"reported,omitempty"`
}

func NewDeviceTwin(desired, reported []*TwinProperty) *DeviceTwin {
	return &DeviceTwin{
		Desired:  desired,
		Reported: reported,
	}
}

// TwinProperty represents the device property for which an Expected/Actual state can be defined.
type TwinProperty struct {
	Service string `form:"service" json:"svc"`
	// Required: The property name for which the desired/reported values are specified.
	// This property should be present in the device model.
	PropertyName string `form:"property_name" json:"pn"`
	// Required: The value for this property.
	Value interface{} `form:"value" json:"val,omitempty"`
	//the timestamp of the property collecting value
	// +optional
	Timestamp int64 `json:"ts,omitempty"`
	// some error message info when collecting this property.
	// +optional, just for reported.
	ErrorMessage string `json:"err_msg,omitempty"`
}

// New device twin property.
func NewTwinProperty(svc, propName, errMsg string, val interface{}) *TwinProperty {
	now := utils.GetNowTimeStamp()

	return &TwinProperty{
		Service:      svc,
		PropertyName: propName,
		Value:        val,
		Timestamp:    now,
		ErrorMessage: errMsg,
	}
}

// DeviceSpec represents a single device instance. It is an instantation of a device model.
type DeviceSpec struct {
	//Required: device name, changed by user
	Name string `json:"Name"`

	//Required: the device belong to which edge.
	EdgeID                   string `form:"edge_id" json:"edge_id"`
	DeviceOS                 string `form:"os" json:"os,omitempty"`
	DeviceCatagory           string `form:"catagory" json:"catagory,omitempty"`
	DeviceVersion            int    `form:"version" json:"version,omitempty"`
	DeviceIdentificationCode string `form:"id_code" json:"id_code,omitempty"`
	Description              string `form:"desc" json:"desc"`

	//group.
	//TODO: reserved in future
	GroupName string `form:"group_name" json:"group_name,omitempty"`
	//who create the device by ID.
	Creator string `form:"creator" json:"creator,omitempty"`

	// +optional
	//TODO: reserved in future
	DeviceAuthType string `form:"auth_type" json:"auth_type,omitempty"`
	Secret         string `form:"secret" json:"secret,omitempty"`

	//Device Type
	//Kind: Direct, GateWay, SubDevcie.
	DeviceType  string `form:"dev_type" json:"dev_type,omitempty"`
	GatewayID   string `form:"gateway_id" json:"gateway_id,omitempty"`
	GatewayName string `form:"gateway_name" json:"gateway_name,omitempty"`
	// Additional metadata like tags.
	// +optional
	Tags map[string]string `form:"tags" json:"tags,omitempty"`

	// Required: DeviceModelRef is reference to the device model used as a template
	// to create the device instance.
	//this should be the Name of device model.
	DeviceModelRef string `form:"model_ref" json:"model_ref,omitempty"`

	ProtocolType string `form:"proto_type" json:"proto_type,omitempty"`
	// Required: The protocol configuration used to connect to the device.
	//this should be a json string
	Protocol string `form:"proto_config" json:"proto_config,omitempty"`
	// The Life time of the desired value.
	LifeTimeOfDesiredValue int64 `form:"ltodv" json:"ltodv,omitempty"`
	// ExtensionConfig which describe how to access the device properties,command, and events.
	// +optional
	ExtensionConfig *ExtensionConfig `form:"ext_config" json:"ext_config,omitempty"`
}

// millisecond: ms  second: s  minute: m  hour: h
type DefaultProtocol struct {
	IntervalUnit string `json:"it_u"`
	Interval     int64  `json:"it"`
}

type ExtensionConfig struct {
	// Required: List of device services.
	Services []*DeviceServiceSpec `json:"services,omitempty"`
}

func (ec *ExtensionConfig) FindDeviceServiceSpec(name string) *DeviceServiceSpec {
	for _, dss := range ec.Services {
		if dss == nil {
			continue
		}

		if dss.Name == name {
			return dss
		}
	}

	return nil
}

func (ec *ExtensionConfig) StoreExtensionConfig(dm *db.DeviceModel, deviceID string) error {
	if ec.Services == nil {
		ec.Services = make([]*DeviceServiceSpec, 0)
	}

	//save all service instance.
	docs := make([]*db.ServiceInstance, 0)
	for _, sm := range dm.ServiceModels {
		if sm == nil {
			continue
		}

		doc := db.NewServiceInstance(sm.Name, deviceID)
		docs = append(docs, doc)
		dss := ec.FindDeviceServiceSpec(sm.Name)
		if dss == nil {
			dss = &DeviceServiceSpec{
				Name: sm.Name,
			}
			ec.Services = append(ec.Services, dss)
		}

		dss.model = sm
		dss.doc = doc
	}

	err := db.AddServiceInstance(docs)
	if err != nil {
		return err
	}

	for _, dss := range ec.Services {
		if dss == nil {
			continue
		}

		err = dss.StoreServiceSpec()
		if err != nil {
			return err
		}
	}

	return nil
}

// DeviceServiceSpec is the  an instantation of a DeviceServiceModel.
type DeviceServiceSpec struct {
	Name       string                `json:"name"`
	Properties []*DevicePropertySpec `json:"props,omitempty"`
	Events     []*DeviceEventSpec    `json:"events,omitempty"`
	Commands   []*DeviceCommandSpec  `json:"cmds,omitempty"`
	model      *db.ServiceModel
	doc        *db.ServiceInstance
}

func (dss *DeviceServiceSpec) FindDevicePropertySpec(name string) *DevicePropertySpec {
	for _, p := range dss.Properties {
		if p.Name == name {
			return p
		}
	}

	return nil
}

func (dss *DeviceServiceSpec) FindDeviceEventSpec(name string) *DeviceEventSpec {
	for _, e := range dss.Events {
		if e.Name == name {
			return e
		}
	}

	return nil
}

func (dss *DeviceServiceSpec) FindDeviceCommandSpec(name string) *DeviceCommandSpec {
	for _, c := range dss.Commands {
		if c.Name == name {
			return c
		}
	}

	return nil
}

func (dss *DeviceServiceSpec) StoreServiceSpec() error {
	sm := dss.model
	serviceID := dss.doc.ID

	if dss.Name == "" {
		dss.Name = sm.Name
	}
	if dss.Properties == nil {
		dss.Properties = make([]*DevicePropertySpec, 0)
	}
	if dss.Events == nil {
		dss.Events = make([]*DeviceEventSpec, 0)
	}
	if dss.Commands == nil {
		dss.Commands = make([]*DeviceCommandSpec, 0)
	}

	//store property instance according to the property models.
	if sm.PropertyModels != nil && len(sm.PropertyModels) > 0 {
		docs := make([]*db.PropertyInstance, 0)

		for _, pm := range sm.PropertyModels {
			if pm == nil {
				continue
			}

			doc := db.NewPropertyInstance(pm.Name, serviceID, "")
			dps := dss.FindDevicePropertySpec(pm.Name)
			if dps != nil {
				doc.AccessConfig = dps.AccessConfig
			}

			docs = append(docs, doc)
		}

		//store all properties in this service.
		err := db.AddPropertyInstance(docs)
		if err != nil {
			return err
		}
	}

	//update events instance according to the event models.
	if sm.EventModels != nil && len(sm.EventModels) > 0 {
		docs := make([]*db.EventInstance, 0)
		for _, em := range sm.EventModels {
			if em == nil {
				continue
			}

			doc := db.NewEventInstance(em.Name, serviceID, "")
			des := dss.FindDeviceEventSpec(em.Name)
			if des != nil {
				doc.AccessConfig = des.AccessConfig
			}

			docs = append(docs, doc)
		}

		//store all events in this service.
		err := db.AddEventInstance(docs)
		if err != nil {
			return err
		}
	}

	//update command instance according to the command models.
	if sm.CommandModels != nil && len(sm.CommandModels) > 0 {
		docs := make([]*db.CommandInstance, 0)

		for _, cm := range sm.CommandModels {
			if cm == nil {
				continue
			}

			doc := db.NewCommandInstance(cm.Name, serviceID, "")
			dcs := dss.FindDeviceCommandSpec(cm.Name)
			if dcs != nil {
				doc.AccessConfig = dcs.AccessConfig
			}

			docs = append(docs, doc)
		}

		//store all events in this service.
		err := db.AddCommandInstance(docs)
		if err != nil {
			return err
		}
	}

	return nil
}

// DevicePropertySpec is an instantation of a DevicePropertyModel.
type DevicePropertySpec struct {
	*DevicePropertyModel `json:",inline"`
	// List of AccessConfig which describe how to access the device properties,command, and events.
	// AccessConfig must unique by AccessConfig.propertyName.
	// +optional
	//this should be a json string
	// AccessConfig must unique by AccessConfig.propertyName.
	AccessConfig string `json:"ac,omitempty"`
}

// DeviceEventSpec is an instantation of a DeviceEventModel.
type DeviceEventSpec struct {
	*DeviceEventModel `json:",inline"`
	// List of AccessConfig which describe how to access the device properties,command, and events.
	// AccessConfig must unique by AccessConfig.propertyName.
	// +optional
	//this should be a json string
	// AccessConfig must unique by AccessConfig.propertyName.
	AccessConfig string `json:"ac,omitempty"`
}

// DeviceCommandSpec is an instantation of a DeviceCommandModel.
type DeviceCommandSpec struct {
	*DeviceCommandModel `json:",inline"`
	// List of AccessConfig which describe how to access the device properties,command, and events.
	// AccessConfig must unique by AccessConfig.propertyName.
	// +optional
	//this should be a json string
	// AccessConfig must unique by AccessConfig.propertyName.
	AccessConfig string `json:"ac,omitempty"`
}

// DeviceStatus reports the device state and the desired/reported values of twin attributes.
type DeviceStatus struct {
	//device status
	// inactive, active, online, offline, error etc.
	DeviceStatus string `json:"dev_stat"`
	//start/stop states.
	State string `json:"state,omitempty"`
	// DeviceTwin provides a logical representation of A series of control properties on a device
	// Required:
	Twins *DeviceTwin `json:"twins"`
}

type EventRules struct {
	// such as: >, <, ==, !=, ...
	Relation *int   `json:"relation,omitempty"`
	Value    string `json:"value,omitempty"`
	// Service name to which the property in the rule belongs
	ServiceName string `json:"serviceName,omitempty"`
	// Property name to which the property in the rule belongs
	PropertyName string `json:"propertyName,omitempty"`
}

type EventsAccessConfig struct {
	// The relationship between rules in the event, such as and or or
	Condition string `json:"condition"`
	// Rules contained in the event
	Rules []*EventRules `json:"rules,omitempty"`
	// Event detection period
	DetectionPeriod *int `json:"detection,omitempty"`
	// Unit of detection period, such as: ms, s, m, h.
	DetectionPeriodUnit string `json:"detectionUnit,omitempty"`
	// How many consecutive detections are required to determine whether the event is valid
	Maintain *int `json:"maintain,omitempty"`
	// Do you only need to report once or continuously
	// once means only need to report once, continued means continuous Report
	ReportType string `json:"reportType,omitempty"`
}

func NewDeviceStatus() *DeviceStatus {
	return &DeviceStatus{
		DeviceStatus: "inactive",
		State:        "stopped",
	}
}

func (ds *DeviceStatus) Start() {
	ds.State = "started"
}

// Device is the Schema for the devices API
type Device struct {
	DeviceID        string `json:"device_id"`
	GroupID         string `json:"group_id,omitempty"`
	CreateTimeStamp int64  `json:"create_timestamp,omitempty"`
	UpdateTimeStamp int64  `json:"update_timestamp,omitempty"`

	Spec   *DeviceSpec   `json:"spec,omitempty"`
	Status *DeviceStatus `json:"status,omitempty"`
}

// DeviceList contains a list of Device
type DeviceList struct {
	Items []*Device `json:"items"`
}

// copy sub device to other edge
func CopyDeviceInstance(deviceInstance *db.DeviceInstance, edgeId string) error {
	deviceInstances := make([]*db.DeviceInstance, 0)
	propertyInstances := make([]*db.PropertyInstance, 0)
	eventInstances := make([]*db.EventInstance, 0)
	commandInstances := make([]*db.CommandInstance, 0)
	serviceInstances := make([]*db.ServiceInstance, 0)

	deviceId := utils.NewUUID()
	deviceInstances = append(deviceInstances, &db.DeviceInstance{
		DeviceID:                 deviceId,
		Name:                     deviceInstance.Name,
		EdgeID:                   edgeId,
		DeviceOS:                 deviceInstance.DeviceOS,
		DeviceCategory:           deviceInstance.DeviceOS,
		DeviceVersion:            deviceInstance.DeviceVersion,
		DeviceIdentificationCode: deviceInstance.DeviceIdentificationCode,
		Description:              deviceInstance.Description,
		GroupName:                deviceInstance.GroupName,
		GroupID:                  deviceInstance.GroupID,
		Creator:                  deviceInstance.Creator,
		DeviceAuthType:           deviceInstance.DeviceAuthType,
		Secret:                   deviceInstance.DeviceAuthType,
		DeviceType:               deviceInstance.DeviceType,
		GatewayID:                deviceInstance.GatewayID,
		GatewayName:              deviceInstance.GatewayName,
		DeviceModelRef:           deviceInstance.DeviceModelRef,
		ProtocolType:             deviceInstance.ProtocolType,
		Protocol:                 deviceInstance.Protocol,
		LifeTimeOfDesiredValue:   deviceInstance.LifeTimeOfDesiredValue,
		CreateTimeStamp:          utils.GetNowTimeStamp(),
		UpdateTimeStamp:          utils.GetNowTimeStamp(),
		DeviceStatus:             global.DeviceStatusInactive,
		State:                    global.DeviceStateStopped,
		DeviceModelId:            deviceInstance.DeviceModelId,
	})
	if len(deviceInstance.ServiceInstances) > 0 {
		for _, service := range deviceInstance.ServiceInstances {
			serviceInstance := db.NewServiceInstance(service.Name, deviceId)
			serviceInstances = append(serviceInstances, serviceInstance)

			if len(service.PropertyInstances) > 0 {
				for _, property := range service.PropertyInstances {
					propertyInstances = append(propertyInstances, &db.PropertyInstance{
						Name:         property.Name,
						ServiceID:    serviceInstance.ID,
						AccessConfig: property.AccessConfig,
					})
				}
			}

			if len(service.EventInstances) > 0 {
				for _, event := range service.EventInstances {
					eventInstances = append(eventInstances, &db.EventInstance{
						Name:         event.Name,
						ServiceID:    serviceInstance.ID,
						AccessConfig: event.AccessConfig,
					})
				}
			}

			if len(service.CommandInstances) > 0 {
				for _, command := range service.CommandInstances {
					commandInstances = append(commandInstances, &db.CommandInstance{
						Name:         command.Name,
						ServiceID:    serviceInstance.ID,
						AccessConfig: command.AccessConfig,
					})
				}
			}
		}

	}
	if err := db.AddDeviceInstanceAll(deviceInstances, serviceInstances, propertyInstances, eventInstances, commandInstances); err != nil {
		klog.Errorln("Add device instance error")
		return err
	}
	return nil
}
