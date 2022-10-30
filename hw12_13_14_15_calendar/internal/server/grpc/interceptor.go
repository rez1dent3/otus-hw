package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *Server) RequestInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	now := time.Now()
	i, err := handler(ctx, req)
	latency := time.Since(now)

	p, ok := peer.FromContext(ctx)
	if !ok {
		s.logger.Error("error get of peer information: " + p.Addr.String())
	}

	userAgent := ""
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userAgent = strings.Join(md["user-agent"], " ")
	}

	ip, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		s.logger.Error("error split host and port: " + p.Addr.String())
	}

	s.logger.Info(fmt.Sprintf(
		"%s [%s] %s %s %s %s %d %s",
		ip,
		now.Format("02/Jan/2006:15:04:05 -0700"),
		"rpc",
		info.FullMethod,
		"HTTP/2",
		status.Code(err).String(),
		latency.Microseconds(),
		userAgent,
	))

	return i, err
}
