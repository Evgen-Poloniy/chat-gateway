package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func Logger(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		requestID := uuid.NewString()

		clientIP := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		entry := logger.WithFields(logrus.Fields{
			"id":     requestID,
			"method": "gRPC",
			"path":   info.FullMethod,
			"ip":     clientIP,
		})

		entry.Info("request start")

		defer func() {
			if r := recover(); r != nil {
				err = status.Error(codes.Internal, "internal server error")

				entry = entry.WithFields(logrus.Fields{
					"status_code": codes.Internal.String(),
					"code":        "PANIC",
				})
				entry.Error(fmt.Sprintf("panic recovered: %v", r))
			}
		}()

		start := time.Now()

		resp, err = handler(ctx, req)

		latency := time.Since(start)

		statusCode := codes.OK
		var errMessage string

		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code()
				errMessage = st.Message()
			} else {
				statusCode = codes.Internal
				errMessage = err.Error()
			}
		}

		entry = logger.WithFields(logrus.Fields{
			"id":          requestID,
			"method":      "gRPC",
			"path":        info.FullMethod,
			"ip":          clientIP,
			"latency":     fmt.Sprintf("%v", latency),
			"status_code": statusCode.String(),
		})

		if err != nil {
			if statusCode == codes.InvalidArgument {
				entry = entry.WithField("code", "bad_request")
			}

			if statusCode == codes.Internal || statusCode == codes.Unknown {
				entry.Error(errMessage)
			} else {
				entry.Warn(errMessage)
			}
			return nil, err
		}

		entry.Info("request completed successfully")

		return resp, nil
	}
}
