package transport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
	"github.com/edgehook/ithings/transport/mqtt"
	"k8s.io/klog/v2"
)

type ithingsTransport struct {
	config  *config.MqttConfig
	ctx     context.Context
	client  *mqtt.Client
	rxQueue chan *types.IMessage
	txQueue chan *types.IMessage
}

func NewithingsTransport(ctx context.Context) (error, *ithingsTransport) {
	conf := config.GetMqttConfig()
	if conf == nil {
		return errors.New("mqtt config is invalid"), nil
	}

	transport := &ithingsTransport{
		config:  conf,
		ctx:     ctx,
		rxQueue: make(chan *types.IMessage, global.DefaultMqttRxQueueSize),
		txQueue: make(chan *types.IMessage, global.DefaultMqttTxQueueSize),
	}

	//TODO: willmessage.
	wm := &mqtt.WillMessage{
		Topic:    "/ithing/me/die",
		Payload:  "/ithing/me/die",
		Qos:      byte(conf.QOS),
		Retained: true,
	}

	//tls config
	tlsConfig, err := createTLSConfig(conf.CaFilePath,
		conf.CertFilePath, conf.KeyFilePath)
	if err != nil {
		tlsConfig = nil
	}

	brokerAddress := conf.Broker
	if !strings.Contains(brokerAddress, ":") {
		brokerAddress = brokerAddress + ":1883"
	}
	klog.Infof("Connect mqtt username=%s, password=%s", conf.User, conf.Passwd)
	//create the mqtt client.
	transport.client = mqtt.NewMQTTClient(brokerAddress, conf.ClientID,
		conf.User, conf.Passwd, tlsConfig, wm)

	transport.client.OnConnectFn = transport.OnConnect
	transport.client.OnLostFn = transport.OnLost
	transport.client.OnReconnectFn = transport.OnReconnect

	return nil, transport
}

func (it *ithingsTransport) connect() error {
	return it.client.Connect()
}

func (it *ithingsTransport) OnConnect(c *mqtt.Client) {
	klog.Infof("On Connected Broker!")

	subTopic := "adv/ithings/edge/#"
	err := it.client.Subscribe(subTopic, byte(it.config.QOS), it.MessageArrived)
	if err != nil {
		klog.Errorf("subscribe with err %s", err.Error())
		return
	}

	klog.Infof("subscribe topic: [%s]", subTopic)

}

func (it *ithingsTransport) Run() {
retry_connect:
	err := it.connect()
	if err != nil {
		select {
		case <-it.ctx.Done():
			klog.Infof("ithings transport will be close for context cancel.")
			it.Close()
			return
		default:
			klog.Errorf("connect with err %v", err)
			time.Sleep(3 * time.Second)
			goto retry_connect
		}
	}

	klog.Infof("connect to broker %s successfuly !", it.config.Broker)

	//read the iMessage from txQueue and send it sequentially
	for {
		select {
		case <-it.ctx.Done():
			klog.Infof("ithings transport will be close for context cancel.")
			it.Close()
			return
		case msg, ok := <-it.txQueue:
			if !ok {
				klog.Warningf("txQueue has been closed!")
				return
			}

			if msg == nil {
				continue
			}

			if msg.Req != nil {
				req := msg.Req
				it.publish(req.BuildTopic(), req.BuildPayload())
			} else {
				if msg.Resp != nil {
					resp := msg.Resp
					it.publish(resp.BuildTopic(), resp.BuildPayload())
				}
			}
		}
	}
}

func (it *ithingsTransport) publish(topic, payload string) error {
	var reties = int(0)
	qos := byte(it.config.QOS)

retry_to_send:
	err := it.client.Publish(topic, payload, qos, false)
	if err != nil {
		if reties > 3 {
			return err
		}

		reties++
		time.Sleep(100 * time.Millisecond)

		goto retry_to_send
	}

	return nil
}

func (it *ithingsTransport) MessageArrived(topic string, payload []byte) {
	msg := &types.IMessage{}

	levels := strings.Split(topic, "/")
	if len(levels) < 7 || payload == nil {
		klog.Warningf("%s is invalid format, we ignored!", topic)
		return
	}

	if strings.Contains(levels[6], types.MSG_OPS_REPLY) {
		//response message.
		msg.Resp = types.ParseResponse(levels, payload)
	} else {
		msg.Req = types.ParseRequest(levels, payload)
	}

	it.rxQueue <- msg
}

func (it *ithingsTransport) GetRxQueueCh() chan *types.IMessage {
	return it.rxQueue
}

func (it *ithingsTransport) PushIMessageToSender(msg *types.IMessage) {
	it.txQueue <- msg
}

func (it *ithingsTransport) OnLost(c *mqtt.Client, err error) {
	klog.Infof("Connect Broker is lost with err: [%s]!", err.Error())
}

func (it *ithingsTransport) OnReconnect(c *mqtt.Client, options *paho.ClientOptions) {
	klog.Infof("try to reconnect the broker!")
}

func (it *ithingsTransport) Close() {
	if it.client != nil {
		it.client.Close()
	}
}

// create tls config
func createTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	pool := x509.NewCertPool()
	rootCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	ok := pool.AppendCertsFromPEM(rootCA)
	if !ok {
		return nil, errors.New("fail to load ca content")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ClientCAs:    pool,
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	}

	return tlsConfig, nil
}
