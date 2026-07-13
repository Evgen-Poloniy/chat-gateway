package grpcserver

import (
	"fmt"
	"net"

	"github.com/Evgen-Poloniy/chat-gateway/internal/interceptor"
	grpcv1 "github.com/Evgen-Poloniy/chat-gateway/internal/transport/grpc/v1"
	dispatchv1 "github.com/Evgen-Poloniy/messenger-contracts/gen/go/dispatch/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(port int, handler *grpcv1.Handler, logger *logrus.Logger) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %v", port, err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Logger(logger)),
	)

	dispatchv1.RegisterDispatchServiceServer(s, handler)

	return &Server{
		grpcServer: s,
		listener:   lis,
	}, nil
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
