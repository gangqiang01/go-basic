package utils

import (
	"encoding/json"

	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
	beehiveCtx "github.com/jwzl/beehive/pkg/core/context"
	"github.com/jwzl/wssocket/model"
	"k8s.io/klog/v2"
)

const (
// operation
)

func BuildModelMessage(source string, target string, operation string, resource string, content interface{}) *model.Message {
	msg := model.NewMessage("")

	//Router
	msg.BuildRouter(source, "", target, resource, operation)

	//content
	msg.Content = content

	return msg
}

func BuildTrans2ICoreMessage(content interface{}) *model.Message {
	return BuildModelMessage(global.IMODULE_TRANSPORT, global.IMODULE_CORE, "do_edge_msg", "", content)
}

func BuildEdgeSendMessage(content interface{}) *model.Message {
	return BuildModelMessage(global.IMODULE_CORE, global.IMODULE_TRANSPORT, "send_edge_msg", "", content)
}

func SendResponse2Edge(req *types.Request, code string, v interface{}) {
	var payload string

	if req != nil {
		switch v.(type) {
		case string:
			payload = v.(string)
		case []byte:
			d := v.([]byte)
			payload = string(d)
		default:
			d, err := json.Marshal(v)
			if err != nil {
				klog.Errorf("json Marshal with  err %v", err)
				return
			}

			payload = string(d)
		}

		resp := req.BuildResponse(code, payload)

		msg := &types.IMessage{
			Resp: resp,
		}

		modelMsg := BuildEdgeSendMessage(msg)
		//klog.Infof("Send responce msg: %s", payload)
		beehiveCtx.Send(global.IMODULE_TRANSPORT, modelMsg)
	}
}

func SendRequest2Edge(req *types.Request) {
	if req != nil {
		msg := &types.IMessage{
			Req: req,
		}

		modelMsg := BuildEdgeSendMessage(msg)
		//m, _ := json.Marshal(req)
		klog.Infof("Send request msg: %s", string(req.Payload.Content))
		beehiveCtx.Send(global.IMODULE_TRANSPORT, modelMsg)
	}
}
