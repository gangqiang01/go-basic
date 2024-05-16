package core

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	db "github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/grp"
	"github.com/edgehook/ithings/common/influxdbm/influx_store"
	"github.com/edgehook/ithings/common/types"
	"github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/core/devicetwin"
	"github.com/edgehook/ithings/core/eventlistener"
	"github.com/edgehook/ithings/dataforward"
	"github.com/edgehook/ithings/rulelinkage"
	"k8s.io/klog/v2"
)

type ICore struct {
	mutex *sync.Mutex
	//go routine pool.
	Pool *grp.GoRoutinePool
	//event Listener Manager
	elm *eventlistener.EventListenerManager
	//device twin manager
	twinMgr *devicetwin.DeviceTwinManager
	//device status manager.
	devStatusMgr *DeviceStatusManager
	//device twin sync cache to edge.
	twinRetrySyncCache *sync.Map
}

func NewICore(maxGoRoutine int) *ICore {
	var coreMutex sync.Mutex
	var tSyncCache sync.Map

	pool := grp.NewGoRoutinePool(maxGoRoutine)

	return &ICore{
		Pool:               pool,
		mutex:              &coreMutex,
		elm:                eventlistener.NewEventListenerManager(),
		twinMgr:            devicetwin.NewDeviceTwinManager(),
		devStatusMgr:       NewDeviceStatusManager(),
		twinRetrySyncCache: &tSyncCache,
	}
}

func (ic *ICore) Start(ctx context.Context) error {
	//run retry sync desired twins loop.
	go ic.runRetrySyncDesiredTwinLoop(ctx)

	return ic.twinMgr.Initialize()
}

/*
* process the message from edge.
 */
func (ic *ICore) DoProcessEdgeMsg(msg *types.IMessage) error {
	if msg.Req != nil {
		//request.
		req := msg.Req

		ic.mutex.Lock()
		defer ic.mutex.Unlock()

		err := ic.Pool.Run(func() {
			switch req.Operation {
			case types.MSG_OPS_REGISTER:
				ic.doProtocolRegister(req)
			case types.MSG_OPS_FETCH:
				ic.doFetchDeviceMeta(req)
			case types.MSG_OPS_REPORT:
				ic.doDevicesReport(req)
			default:
			}
		})
		if err != nil {
			klog.Errorf("New goroutine pool failed with err %v", err)
			return err
		}
	} else {
		if msg.Resp != nil {
			//response message.
			resp := msg.Resp

			ic.mutex.Lock()
			defer ic.mutex.Unlock()

			err := ic.Pool.Run(func() {
				ic.doProcessResponse(resp)
			})
			if err != nil {
				klog.Errorf("New goroutine pool failed with err %v", err)
				return err
			}
		}
	}

	return nil
}

func (ic *ICore) doProtocolRegister(req *types.Request) {
	if req == nil {
		return
	}

	regParms := db.ProtocolTypes{}
	content := req.GetContent()

	err := json.Unmarshal([]byte(content), &regParms)
	if err != nil {
		klog.Errorf("json Unmarshal with err %v", err)
		utils.SendResponse2Edge(req, global.IRespCodeError, err.Error())
		return
	}

	err = db.AddProtocolType(&regParms)
	if err != nil {
		klog.Errorf("AddProtocolType with err %v", err)
		utils.SendResponse2Edge(req, global.IRespCodeError, err.Error())
		return
	}

	klog.Infof("register/update the protocol %s", regParms.ProtocolType)

	utils.SendResponse2Edge(req, global.IRespCodeOk, global.IRespOkString)
}

func (ic *ICore) doFetchDeviceMeta(req *types.Request) {
	if req == nil {
		return
	}

	edgeID := req.EdgeID
	protocType := req.MapperID

	devSpecList, err := GetDeviceSpecMeta(edgeID, protocType)
	if err != nil {
		klog.Errorf("GetDeviceSpecMeta with err %v", err)
		utils.SendResponse2Edge(req, global.IRespCodeError, "GetDeviceSpecMeta with err"+err.Error())
		return
	}

	utils.SendResponse2Edge(req, global.IRespCodeOk, devSpecList)

	//record the fetch action from mapper.
	ic.devStatusMgr.UpdateFetchHistory(edgeID, protocType, true)
}

/*
* doDevicesReport:
* process the device status.
 */
func (ic *ICore) doDevicesReport(req *types.Request) {
	if req == nil {
		return
	}

	rsc := req.Resource
	edgeID := req.EdgeID
	//protocType := req.MapperID
	content := req.Payload.Content

	// we will not process any message to edge since
	// edge is offline.
	if !EdgeIsOnline(edgeID) {
		klog.Warningf("edge(%s) is not already", edgeID)
		utils.SendResponse2Edge(req, global.IRespCodeError, global.ErrEdgeIsNotOnline.Error())
		return
	}

	switch {
	case strings.Contains(rsc, "device_status"):
		msg, err := decodeReportDevicesStatusMessage(content)
		//msgjson, _ := json.Marshal(msg)
		//klog.Infof("Receive device_status msg: %v##edgeId: %s", string(msgjson), edgeID)
		if err != nil {
			klog.Errorf("Decode edge report message with err %v", err)
			utils.SendResponse2Edge(req, global.IRespCodeInvalidMsg, global.IRespInvalidMsgString)
			return
		}

		//update device status.
		if err = ic.updateDeviceStatus(edgeID, msg.DevicesStatus); err != nil {
			klog.Errorf("Update device status with err %v", err)
			utils.SendResponse2Edge(req, global.IRespCodeInvalidMsg, global.IRespInvalidMsgString)
			return
		}
	case strings.Contains(rsc, "device_data"):
		msg, err := decodeReportDevicesMessage(content)
		//log
		//msgjson, _ := json.Marshal(msg)
		//klog.Infof("Receive device_data msg: %v##edgeId: %s", string(msgjson), edgeID)

		if err != nil {
			klog.Errorf("Decode edge report message with err %v", err)
			utils.SendResponse2Edge(req, global.IRespCodeInvalidMsg, global.IRespInvalidMsgString)
			return
		}
		//manager the device reported twins
		if err := ic.Pool.Run(func() {
			err = ic.twinMgr.AddDevicesTwinsUpdateData(msg, edgeID)
			if err != nil {
				klog.Errorf("AddDevicesTwinsUpdateData with err %v", err)
				utils.SendResponse2Edge(req, global.IRespCodeInternalError, global.IRespInternalErrString)
				return
			}
		}); err != nil {
			klog.Errorf("[device_data] New goroutine pool failed with err %v", err)
		}

	case strings.Contains(rsc, "device_event"):
		msg, err := DecodeReportEventMessage(content)
		//log
		msgjson, _ := json.Marshal(msg)
		klog.Infof("Receive device_event msg: %v", string(msgjson))
		if err != nil {
			klog.Errorf("Decode edge report event message with err %v", err)
			utils.SendResponse2Edge(req, "201", "invalid message format")
			return
		}

		//store influxDB
		if err := ic.Pool.Run(func() {
			if err := influx_store.StoreEvent(msg); err != nil {
				klog.Errorf("Store influxDB %s twins failed with err %v", msg.DeviceID, err)
			}
		}); err != nil {
			klog.Errorf("[device_event][StoreEvent] New goroutine pool failed with err %v", err)
		}

		//data forward
		if err := ic.Pool.Run(func() {
			if err := dataforward.HandleDataForward(msg.DeviceID, edgeID, "event", msg); err != nil {
				klog.Errorf("Data forward with err %v", err)
			}
		}); err != nil {
			klog.Errorf("[device_event][HandleDataForward]New goroutine pool failed with err %v", err)
		}

		//handle rule
		if err := ic.Pool.Run(func() {
			err = rulelinkage.HandleRule(msg, edgeID, UpdateDeviceDesiredTwins, v1.RuleLinkageHandleTypeEvent)
			if err != nil {
				klog.Errorf("handleRule with err %v", err)
			}
		}); err != nil {
			klog.Errorf("[device_event][HandleRule]New goroutine pool failed with err %v", err)
		}

	case strings.Contains(rsc, "event_recover"):
		msg, err := DecodeReportEventMessage(content)
		//log
		msgjson, _ := json.Marshal(msg)
		klog.Infof("Receive event_recover msg: %v", string(msgjson))
		if err != nil {
			klog.Errorf("Decode edge report event message with err %v", err)
			utils.SendResponse2Edge(req, "201", "invalid message format")
			return
		}

		//recover rule
		err = rulelinkage.HandleRule(msg, edgeID, UpdateDeviceDesiredTwins, v1.RuleLinkageHandleTypeRecover)
		if err != nil {
			klog.Errorf("handleRule with err %v", err)
		}
	}

	utils.SendResponse2Edge(req, global.IRespCodeOk, global.IRespOkString)
}

func (ic *ICore) doProcessResponse(resp *types.Response) {
	if resp == nil {
		return
	}

	msgID := resp.GetMsgParentID()
	ic.elm.MatchEventAndDispatch(msgID, resp)
}

/*
* SendSyncRequestToEdge
* send the request and wait response.
* timeOut = 0, we will wait it forver.
* this function is thread safe.
 */
func (ic *ICore) SendSyncRequestToEdge(req *types.Request, timeOut time.Duration) (*types.Response, error) {
	if req == nil {
		return nil, global.ErrInvalidParms
	}

	//send request to edge
	utils.SendRequest2Edge(req)

	//wait the message response.
	msgID := req.GetMessageID()
	v, err := ic.elm.WatchEvent(msgID, timeOut)
	if err != nil {
		return nil, err
	}

	resp, isThisType := v.(*types.Response)
	if !isThisType || resp == nil {
		return nil, global.ErrInvalidResponseStruct
	}

	return resp, nil
}
