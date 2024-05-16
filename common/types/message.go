package types

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

const (
	SERVER_TOPIC_PREFIX = "adv/ithings/server"
	EDGE_TOPIC_PREFIX   = "adv/ithings/edge"

	MSG_OPS_REGISTER     = "register"
	MSG_OPS_REPLY        = "reply"
	MSG_OPS_REPORT       = "report"
	MSG_OPS_FETCH        = "fetch"
	MSG_OPS_LIFE_CONTROL = "life_control"
	MSG_OPS_SET_PROPERTY = "set_property"
)

type RequestPayload struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Request struct {
	EdgeID    string
	MapperID  string
	Resource  string
	Operation string

	//message payload.
	Payload RequestPayload
}

// BuildRequest.
func BuildRequest(edgeID, mapperID, resource, operation string) *Request {
	req := &Request{
		EdgeID:    edgeID,
		MapperID:  mapperID,
		Resource:  resource,
		Operation: operation,
	}

	req.Payload.ID = newUUID()

	return req
}

/*
* Parse the request package.
 */
func ParseRequest(levels []string, payload []byte) *Request {
	if len(levels) < 7 {
		return nil
	}

	req := &Request{
		EdgeID:    levels[3],
		MapperID:  levels[5],
		Operation: levels[6],
	}

	if len(levels) > 7 {
		req.Resource = levels[7]
	}

	//parse the payload.
	reqPayload := &req.Payload
	err := json.Unmarshal(payload, reqPayload)
	if err != nil {
		klog.Errorf("ParseRequest err: %v", err)
		return nil
	}

	return req
}

func (r *Request) SetContent(v interface{}) {
	var content string

	switch v.(type) {
	case string:
		content = v.(string)
	case []byte:
		d := v.([]byte)
		content = string(d)
	default:
		d, err := json.Marshal(v)
		if err != nil {
			klog.Errorf("json Marshal with  err %v", err)
			return
		}

		content = string(d)
	}

	r.Payload.Content = content
}

func (r *Request) GetContent() string {
	return r.Payload.Content
}

func (r *Request) GetMessageID() string {
	return r.Payload.ID
}

func (r *Request) BuildTopic() string {
	topic := SERVER_TOPIC_PREFIX + "/" + r.EdgeID + "/mapper/" +
		r.MapperID + "/" + r.Operation
	if r.Operation != MSG_OPS_REPLY {
		topic += "/" + r.Resource
	}

	return topic
}

func (r *Request) BuildPayload() string {
	bytes, _ := json.Marshal(r.Payload)
	return string(bytes)
}

/*
* Build Response message.
 */
func (r *Request) BuildResponse(code string, payload string) *Response {
	resp := &Response{
		EdgeID:    r.EdgeID,
		MapperID:  r.MapperID,
		Operation: MSG_OPS_REPLY,
	}

	resp.Payload.ID = newUUID()
	resp.Payload.ParentID = r.Payload.ID
	resp.Payload.Code = code
	resp.Payload.Content = payload

	return resp
}

type ResponsePayload struct {
	ID       string `json:"id"`
	ParentID string `json:"pid"`
	Code     string `json:"code"`
	Content  string `json:"content"`
}

type Response struct {
	EdgeID    string
	MapperID  string
	Operation string

	//message body.
	Payload ResponsePayload
}

/*
* ParseResponse
* parse the async response message.
 */
func ParseResponse(levels []string, payload []byte) *Response {
	//parse topic
	if len(levels) < 7 {
		return nil
	}

	resp := &Response{
		EdgeID:    levels[3],
		MapperID:  levels[5],
		Operation: levels[6],
	}

	//parse the payload.
	respPayload := &resp.Payload
	err := json.Unmarshal(payload, respPayload)
	if err != nil {
		klog.Errorf("ParseResponse err: %v", err)
		return nil
	}

	return resp
}

func (r *Response) BuildTopic() string {
	return SERVER_TOPIC_PREFIX + "/" + r.EdgeID + "/mapper/" +
		r.MapperID + "/" + MSG_OPS_REPLY
}

func (r *Response) BuildPayload() string {
	bytes, _ := json.Marshal(r.Payload)
	return string(bytes)
}

func (r *Response) GetMsgParentID() string {
	return r.Payload.ParentID
}

/*
* Ithings Message for transport.
 */
type IMessage struct {
	Req  *Request
	Resp *Response
}

func newUUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}
