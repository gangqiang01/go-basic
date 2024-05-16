package transport

import (
	"context"
	"strings"

	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
	"github.com/edgehook/ithings/common/utils"
	"github.com/jwzl/beehive/pkg/core"
	beehiveCtx "github.com/jwzl/beehive/pkg/core/context"
	"github.com/jwzl/wssocket/model"
	"k8s.io/klog/v2"
)

type Transport struct {
	iSyncorClient *ISyncor
	iTransport    *ithingsTransport
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

// Register this module.
func Register() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	tp := &Transport{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	core.Register(tp)
}

// Name
func (tp *Transport) Name() string {
	return global.IMODULE_TRANSPORT
}

// Group
func (tp *Transport) Group() string {
	return global.IMODULE_TRANSPORT
}

// Enable indicates whether this module is enabled
func (tp *Transport) Enable() bool {
	//The module is always enabled!
	return true
}

// Start
func (tp *Transport) Start() {
	tp.iSyncorClient = NewISyncor(tp.ctx)
	if tp.iSyncorClient == nil {
		return
	}

	go tp.iSyncorClient.Run()
	go tp.iSyncorClient.ReceiveAppHubMsg()

	err, iTrans := NewithingsTransport(tp.ctx)
	if err != nil {
		klog.Errorf("new ithings mqtt transport with err %v", err)
		panic(err.Error())
	}

	tp.iTransport = iTrans

	//run ithings mqtt transport.
	go tp.iTransport.Run()
	//dispatch the messgae from edge.
	go tp.edgeMessageDispatch()

	for {
		select {
		case <-beehiveCtx.Done():
			klog.Warning("transport stopped since beehive context canceled!")
			tp.Cleanup()
			return
		default:
			tp.modelMessageRecvive()
		}
	}
}

// Cleanup
func (tp *Transport) Cleanup() {
}

func (tp *Transport) modelMessageRecvive() {
	v, err := beehiveCtx.Receive(global.IMODULE_TRANSPORT)
	if err != nil {
		klog.Errorf("behive channel with err %v", err)
		return
	}

	modelMessage, isThisType := v.(*model.Message)
	if !isThisType || modelMessage == nil {
		//invalid message type or msg == nil, Ignored.
		return
	}

	ops := modelMessage.GetOperation()
	if strings.Contains(ops, "send_edge_msg") {
		msg, isVaild := modelMessage.GetContent().(*types.IMessage)
		if isVaild && msg != nil {
			tp.iTransport.PushIMessageToSender(msg)
		}
	}
}

/*
* dispatch the message from edge.
 */
func (tp *Transport) edgeMessageDispatch() {
	for {
		select {
		case <-beehiveCtx.Done():
			klog.Warning("Transport stopped!")
			tp.Cleanup()
			return
		case msg, ok := <-tp.iTransport.GetRxQueueCh():
			if ok && msg != nil {
				modelMsg := utils.BuildTrans2ICoreMessage(msg)
				beehiveCtx.Send(global.IMODULE_CORE, modelMsg)
			}
		}
	}
}
