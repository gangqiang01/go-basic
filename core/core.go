package core

import (
	"context"
	"strings"

	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
	"github.com/jwzl/beehive/pkg/core"
	beehiveCtx "github.com/jwzl/beehive/pkg/core/context"
	"github.com/jwzl/wssocket/model"
	"k8s.io/klog/v2"
)

type ithingsCore struct {
	iCore      *ICore
	ctx        context.Context
	cancelFunc context.CancelFunc
}

var defaultICore *ithingsCore

// Register this module.
func Register() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ic := &ithingsCore{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	defaultICore = ic

	core.Register(ic)
}

// Name
func (ic *ithingsCore) Name() string {
	return global.IMODULE_CORE
}

// Group
func (ic *ithingsCore) Group() string {
	return global.IMODULE_CORE
}

// Enable indicates whether this module is enabled
func (ic *ithingsCore) Enable() bool {
	//The module is always enabled!
	return true
}

// Start
func (ic *ithingsCore) Start() {
	var maxGoRoutine int

	conf := config.GetMqttConfig()
	if conf == nil {
		maxGoRoutine = global.DefaultMaxGoRoutines
	} else {
		maxGoRoutine = conf.MaxGoRoutine
	}

	ic.iCore = NewICore(maxGoRoutine)

	//start the icore.
	err := ic.iCore.Start(ic.ctx)
	if err != nil {
		klog.Errorf("icore start with err %v", err)
		return
	}

	for {
		select {
		case <-beehiveCtx.Done():
			klog.Warning("core stopped since beehive context canceled!")
			ic.Cleanup()
			return
		default:
			ic.modelMessageRecvive()
		}
	}
}

func (ic *ithingsCore) Cleanup() {
	ic.cancelFunc()
}

func (ic *ithingsCore) modelMessageRecvive() {
	v, err := beehiveCtx.Receive(global.IMODULE_CORE)
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
	if strings.Contains(ops, "do_edge_msg") {
		msg, isVaild := modelMessage.GetContent().(*types.IMessage)
		if !isVaild || msg == nil {
			return
		}

		//process the message from edge on many go routines.
		ic.iCore.DoProcessEdgeMsg(msg)
	}
}
