package isync

import (
	"context"
	"errors"
	"github.com/edgehook/ithings/common/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"k8s.io/klog/v2"
	"time"
)

type ISyncGRPCClient struct {
	ctx context.Context
	//grpc client
	Client ISyncClient
	Conn   *grpc.ClientConn
	//report streams
	SyncAgentStatusStream     ISync_SyncAgentStateClient
	SyncSendMsgToIthingStream ISync_SendMsgToIthingClient
}

var (
	// BackoffMaxDelay provided maximum delay when backing off after failed connection attempts.
	BackoffMaxDelay  = 1 * time.Second
	KeepAliveTime    = time.Duration(10) * time.Second
	KeepAliveTimeout = time.Duration(3) * time.Second

	IsyncClient *ISyncGRPCClient
)

func NewISyncGRPCClient(ctx context.Context, serverAddr, certFile string) (*ISyncGRPCClient, error) {
	var opts = make([]grpc.DialOption, 0)

	//set grpc debug logger.
	if false {
		grpclog.SetLoggerV2(&utils.LoggerWrapper{})
	}

	if certFile != "" {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			klog.Errorf("fail to NewClientTLSFromFile: %v", err)
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// WithBlock returns a DialOption which makes callers of Dial block until the
	// underlying connection is up. Without this, Dial returns immediately and
	// connecting the server happens in background.
	opts = append(opts, grpc.WithBlock())

	//add the autoconnect when connection is lost.
	opts = append(opts, grpc.WithBackoffMaxDelay(BackoffMaxDelay))
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// After a duration of this time if the client doesn't see any activity it
		// pings the server to see if the transport is still alive.
		// If set below 10s, a minimum value of 10s will be used instead.
		Time: KeepAliveTime,
		// After having pinged for keepalive check, the client waits for a duration
		// of Timeout and if no activity is seen even after that the connection is
		// closed.
		Timeout: KeepAliveTimeout,
		// If true, client sends keepalive pings even with no active RPCs. If false,
		// when there are no active RPCs, Time and Timeout will be ignored and no
		// keepalive pings will be sent.
		PermitWithoutStream: true,
	}))

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		klog.Errorf("fail to dial: %s", err.Error())
		return nil, err
	}

	klog.Infof("Connecting the grpc server %s successfully!", serverAddr)

	client := NewISyncClient(conn)
	IsyncClient = &ISyncGRPCClient{
		ctx:    ctx,
		Conn:   conn,
		Client: client,
	}
	return IsyncClient, nil
}

/*
* Get Agent status.
 */
func (c *ISyncGRPCClient) GetAgentStatus() ([]*ReportAgentState, error) {
	req := &AgentStateRequest{
		Tag: "GET_AGENT_STATUS",
	}
	if c.Client == nil {
		return nil, errors.New("Grpc client is nil")
	}
	resp, err := c.Client.GetAgentState(c.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.AgentInfo, nil
}

// send msg to AppHub from ithings
func (c *ISyncGRPCClient) SendMsgToAppHub(tag, body string) (*AppHubResponse, error) {
	req := &AppHubRequest{
		Tag:  tag,
		Body: body,
	}
	if c.Client == nil {
		return nil, errors.New("Grpc client is nil")
	}
	resp, err := c.Client.SendMsgToAppHub(c.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
* Get sync agent status stream
 */
func (c *ISyncGRPCClient) GetSyncAgentStatusStream() error {
	if c.Client == nil {
		return errors.New("Grpc client is nil")
	}
	stream, err := c.Client.SyncAgentState(c.ctx)
	if err != nil {
		return err
	}

	c.SyncAgentStatusStream = stream

	return nil
}

func (c *ISyncGRPCClient) GetSyncSendMsgToIthingsStream() error {
	if c.Client == nil {
		return errors.New("Grpc client is nil")
	}
	stream, err := c.Client.SendMsgToIthing(c.ctx)
	if err != nil {
		return err
	}

	c.SyncSendMsgToIthingStream = stream

	return nil
}

func (c *ISyncGRPCClient) DoSyncAgentStatus(doSync func(*ReportAgentState) (string, string)) error {
	if c.SyncAgentStatusStream == nil {
		err := c.GetSyncAgentStatusStream()
		if err != nil {
			return err
		}
	}

	stream := c.SyncAgentStatusStream

	syncStatus, err := stream.Recv()
	if err != nil {
		c.SyncAgentStatusStream = nil
		return err
	}

	code, reason := doSync(syncStatus)
	resp := &Response{
		StatusCode: code,
		Msg:        reason,
	}

	err = c.SyncAgentStatusStream.Send(resp)
	if err != nil {
		c.SyncAgentStatusStream = nil
		return err
	}

	return nil
}

// send msg to AppHub from ithings
func SendMsgToAppHub(tag, body string) (*AppHubResponse, error) {
	if IsyncClient == nil {
		return nil, errors.New("IsyncClient is nil")
	}
	resp, err := IsyncClient.SendMsgToAppHub(tag, body)
	if err != nil {
		klog.Errorf("Send Msg to AppHub error: %s", err.Error())
	}

	return resp, err
}

// send msg to ithings from AppHub
func (c *ISyncGRPCClient) DoSendMsgToIthings(doSync func(*AppHubRequest) (string, string)) error {
	if c.SyncSendMsgToIthingStream == nil {
		err := c.GetSyncSendMsgToIthingsStream()
		if err != nil {
			return err
		}
	}

	stream := c.SyncSendMsgToIthingStream

	syncMsg, err := stream.Recv()
	if err != nil {
		c.SyncSendMsgToIthingStream = nil
		return err
	}

	code, reason := doSync(syncMsg)
	resp := &Response{
		StatusCode: code,
		Msg:        reason,
	}

	err = c.SyncSendMsgToIthingStream.Send(resp)
	if err != nil {
		c.SyncSendMsgToIthingStream = nil
		return err
	}
	return nil
}

// Close the grpc connection.
func (c *ISyncGRPCClient) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
