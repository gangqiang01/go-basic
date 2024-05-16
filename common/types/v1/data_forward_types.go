package v1

const (
	//log status
	SOURCETYPEMODEL string = "model"
	SOURCETYPEEDGE  string = "edge"
	SOURCETYPETAG   string = "tag"

	DestinationTypeAMQP  string = "amqp"
	DestinationTypeKafka string = "kafka"
	DestinationTypeMqtt  string = "mqtt"
)

// dataForward
type DataForwardSource struct {
	Type     string `form:"type" json:"type"`
	EdgeId   string `form:"edgeId" json:"edgeId,omitempty"`
	ModelId  *int64 `form:"modelId" json:"modelId,omitempty"`
	DeviceId string `form:"deviceId" json:"deviceId,omitempty"`
}

type DestinationAmqp struct {
	//username, password, host string, queueName
	QueueName string `form:"queueName" json:"queueName,omitempty"`
}

type DestinationKafka struct {
	Topic string `json:"topic,omitempty"`
}

type DestinationMqtt struct {
}

type DataForwardDestination struct {
	Type     string `form:"type" json:"type"`
	Host     string `form:"host" json:"host,omitempty"`
	UserName string `form:"username" json:"username,omitempty"`
	Password string `form:"password" json:"password,omitempty"`
	//amqp
	*DestinationAmqp `json:",inline,omitempty"`
	//kafka
	*DestinationKafka `json:",inline,omitempty"`
	//mqtt
	*DestinationMqtt `json:",inline,omitempty"`
}
