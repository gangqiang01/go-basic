package isync

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/keepalive"
	"k8s.io/klog/v2"
)

type ISyncGRPCServer struct {
	ctx context.Context
	//grpc server
	Server *grpc.Server
	Ln     net.Listener
}

func NewISyncGRPCServer(ctx context.Context, serverAddr, certFile, keyFile string) (*ISyncGRPCServer, error) {
	var opts = make([]grpc.ServerOption, 0)

	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		klog.Errorf("failed to listen: %v", err)
		return nil, err
	}

	if certFile != "" && keyFile != "" {
		creds, err := credentials.NewClientTLSFromFile(certFile, keyFile)
		if err != nil {
			klog.Errorf("fail to NewClientTLSFromFile: %v", err)
			return nil, err
		}
		opts = append(opts, grpc.Creds(creds))
	}
	grpcServer := grpc.NewServer(opts...)

	return &ISyncGRPCServer{
		Ln:     ln,
		Server: grpcServer,
	}, nil
}

func (is *ISyncGRPCServer) Run() {
	RegisterISyncServer(is.Server, is)
	is.Server.Serve(is.Ln)
}

func (is *ISyncGRPCServer) SyncAgentState(s ISync_SyncAgentStateServer) error {
	for {
		err := s.Send(&ReportAgentState{AgentId: "12212", Status: "ONLINE"})
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		_, err = s.Recv()
		if err != nil {
			return err
		}
	}

	return nil
}
func (is *ISyncGRPCServer) SendMsgToIthing(s ISync_SendMsgToIthingServer) error {
	for {
		err := s.Send(&AppHubRequest{Tag: "alert", Body: "test"})
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
		_, err = s.Recv()
		if err != nil {
			return err
		}
	}

	return nil
}
func (is *ISyncGRPCServer) GetAgentState(context.Context, *AgentStateRequest) (*AgentStateResponse, error) {
	return &AgentStateResponse{}, nil
}

func (is *ISyncGRPCServer) SendMsgToAppHub(context.Context, *AppHubRequest) (*AppHubResponse, error) {
	return &AppHubResponse{}, nil
}
func (is *ISyncGRPCServer) mustEmbedUnimplementedISyncServer() {
}
