package v1

import (
	"fmt"
	"github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog"
	"time"
)

// Represents the type and data validation of a property.
// Only one of its members may be specified.
type PropertyType struct {
	// +optional
	Int *PropertyTypeInt64 `json:"int,omitempty"`
	// +optional
	String *PropertyTypeString `json:"string,omitempty"`
	// +optional
	Double *PropertyTypeDouble `json:"double,omitempty"`
	// +optional
	Float *PropertyTypeFloat `json:"float,omitempty"`
	// +optional
	Boolean *PropertyTypeBoolean `json:"boolean,omitempty"`
	// +optional
	Bytes *PropertyTypeBytes `json:"bytes,omitempty"`
}

type PropertyTypeInt64 struct {
	// +optional
	DefaultValue int64 `json:"default_value,omitempty"`
	// +optional
	Minimum int64 `json:"min,omitempty"`
	// +optional
	Maximum int64 `json:"max,omitempty"`
	// The unit of the property
	// +optional
	Unit string `json:"unit,omitempty"`
}

type PropertyTypeString struct {
	// +optional
	DefaultValue string `json:"default_value,omitempty"`
}

type PropertyTypeDouble struct {
	// +optional
	DefaultValue float64 `json:"default_value,omitempty"`
	// +optional
	Minimum float64 `json:"min,omitempty"`
	// +optional
	Maximum float64 `json:"max,omitempty"`
	// The unit of the property
	// +optional
	Unit string `json:"unit,omitempty"`
}

type PropertyTypeFloat struct {
	// +optional
	DefaultValue float32 `json:"default_value,omitempty"`
	// +optional
	Minimum float32 `json:"min,omitempty"`
	// +optional
	Maximum float32 `json:"max,omitempty"`
	// The unit of the property
	// +optional
	Unit string `json:"unit,omitempty"`
}

type PropertyTypeBoolean struct {
	// +optional
	DefaultValue bool `json:"default_value,omitempty"`
}

type PropertyTypeBytes struct {
}

// device command param defination
type RequestParamDef struct {
	DataType PropertyType `json:"dt,omitempty"`
}

// The access mode for  a device property.
// +kubebuilder:validation:Enum=ReadWrite;ReadOnly
type PropertyAccessMode string

// Access mode constants for a device property.
const (
	ReadWrite PropertyAccessMode = "ReadWrite"
	ReadOnly  PropertyAccessMode = "ReadOnly"
)

// device property defination
type DevicePropertyModel struct {
	// Required: The device property name.
	Name        string  `form:"pn" json:"pn" binding:"required"`
	Report      bool    `form:"report" json:"report"`
	WriteAble   bool    `form:"rw" json:"rw"`
	MaxValue    float64 `form:"max" json:"max"`
	MinValue    float64 `form:"min" json:"min"`
	Unit        string  `form:"un" json:"un"`
	DataType    string  `form:"dt" json:"dt"`
	Description string  `form:"desc" json:"desc"`
}

type DeviceCommandModel struct {
	Name        string `json:"cn"`
	Description string `json:"desc,omitempty"`
	// params
	RequestParam map[string]string `json:"req_param,omitempty"`
}

// device event defination
type DeviceEventModel struct {
	Name        string  `form:"en" json:"en"  binding:"required"`
	EventType   string  `form:"et" json:"et,omitempty"`
	MaxValue    float64 `form:"max" json:"max,omitempty"`
	MinValue    float64 `form:"min" json:"min,omitempty"`
	Unit        string  `form:"un" json:"un,omitempty"`
	DataType    string  `form:"dt" json:"dt,omitempty"`
	Description string  `form:"desc" json:"desc,omitempty"`
}

// service describe a based module contains some properties, events
// and commands.
type DeviceServiceModel struct {
	Name           string                 `form:"name" json:"name" binding:"required"`
	Description    string                 `form:"desc" json:"desc,omitempty"`
	PropertyModels []*DevicePropertyModel `json:"prop_models,omitempty"`
	EventModels    []*DeviceEventModel    `json:"event_models,omitempty"`
	CommandModels  []*DeviceCommandModel  `json:"cmd_models,omitempty"`
}

// DeviceModelSpec defines the model / template for a device.It is a blueprint which describes the device
// capabilities.
type DeviceModelSpec struct {
	// Required: List of device services.
	ServiceModels []*DeviceServiceModel `json:"svc_models,omitempty"`
}

// DeviceModel is the Schema for the device model API
type DeviceModel struct {
	//update accessconfig and update protocol
	Name         string `form:"name" json:"name" binding:"required"`
	Description  string `form:"description" json:"description,omitempty"`
	Manufacturer string `form:"manufacturer" json:"manufacturer,omitempty"`
	Industry     string `form:"industry" json:"industry,omitempty"`
	DataFormat   string `form:"data_format" json:"data_format,omitempty"`
	DeviceNumber int64  `json:"device_number,omitempty"`
	TagNumber    int64  `json:"tag_number,omitempty"`
	GroupID      string `json:"group_id,omitempty"`
	//who create the device by ID.
	Creator         string `form:"creator" json:"creator"`
	CreateTimeStamp int64  `json:"create_timestamp,omitempty"`
	UpdateTimeStamp int64  `json:"update_timestamp,omitempty"`
	//Spec.
	//Spec DeviceModelSpec `json:"spec,omitempty"`
	*DeviceModelSpec `json:",inline,omitempty"`
}

// DeviceModelList contains a list of DeviceModel
type DeviceModelList struct {
	Items []DeviceModel `json:"items"`
}

func AddAllDeviceModel(deviceModel *DeviceModel) error {
	var (
		index      = 1
		deviceName = ""
	)
Loop:
	if isExist := model.IsExistDeviceModelByName(deviceModel.Name); isExist {
		deviceName = fmt.Sprintf("%s%d", deviceModel.Name, index)
		klog.Infof("deviceName:%s", deviceName)
		index++
		goto Loop
	} else {
		deviceName = deviceModel.Name
	}
	deviceModel.CreateTimeStamp = time.Now().UnixNano() / 1e6
	if err := global.DBAccess.Create(&model.DeviceModel{
		Name:            deviceName,
		Manufacturer:    deviceModel.Manufacturer,
		Industry:        deviceModel.Industry,
		Description:     deviceModel.Description,
		DataFormat:      deviceModel.DataFormat,
		Creator:         deviceModel.Creator,
		CreateTimeStamp: deviceModel.CreateTimeStamp,
	}).Error; err != nil {
		klog.Errorf("Create deviceModel err: %v", err)
		return err
	}

	dmodel, err := model.GetDeviceModelByName(deviceModel.Name)
	if err != nil {
		klog.Errorf("Check deviceModel err: %v", err)
		return err
	}

	if len(deviceModel.ServiceModels) > 0 {
		for _, service := range deviceModel.ServiceModels {
			if isExist := model.IsExistServiceModel(dmodel.ID, service.Name); !isExist {
				if err := global.DBAccess.Create(&model.ServiceModel{
					Name:          service.Name,
					Description:   service.Description,
					DeviceModelId: dmodel.ID,
				}).Error; err != nil {
					klog.Errorf("Create serviceModel err: %v", err)
					return err
				}
			}

			smodel, err := model.GetServiceModelByDeviceModelIdAndServiceModelName(dmodel.ID, service.Name)
			if err != nil {
				klog.Errorf("Check deviceModel err: %v", err)
				return err
			}
			tx := global.DBAccess.Begin()

			if len(service.PropertyModels) > 0 {
				for _, property := range service.PropertyModels {
					if isExist := model.IsExistPropertyModel(smodel.ID, property.Name); !isExist {
						if err := tx.Create(&model.PropertyModel{
							Name:           property.Name,
							Description:    property.Description,
							ServiceModelId: smodel.ID,
						}).Error; err != nil {
							klog.Errorf("Create propertyModel err: %v", err)
							tx.Rollback()
							return err
						}
					}

				}
			}

			if len(service.EventModels) > 0 {
				for _, event := range service.EventModels {
					if isExist := model.IsExistEventModel(smodel.ID, event.Name); !isExist {
						if err := tx.Create(&model.EventModel{
							Name:           event.Name,
							Description:    event.Description,
							ServiceModelId: smodel.ID,
						}).Error; err != nil {
							klog.Errorf("Create eventModel err: %v", err)
							tx.Rollback()
							return err
						}
					}

				}
			}

			if len(service.CommandModels) > 0 {
				for _, command := range service.CommandModels {
					if isExist := model.IsExistCommandModel(smodel.ID, command.Name); !isExist {
						if err := tx.Create(&model.CommandModel{
							Name:           command.Name,
							Description:    command.Description,
							ServiceModelId: smodel.ID,
						}).Error; err != nil {
							klog.Errorf("Create commandModel err: %v", err)
							tx.Rollback()
							return err
						}
					}

				}
			}
			tx.Commit()
		}

	}
	return nil
}
