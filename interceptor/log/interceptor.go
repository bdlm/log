// Package log contains interceptor/middleware helpers for logging.
package log

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/bdlm/log/interceptor"

	"github.com/bdlm/log"
	std "github.com/bdlm/std/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Config contains the log interceptor configuration.
type Config struct {
	// LogStreamRecvMsg if true, log out the contents of each received stream
	// message.
	LogStreamRecvMsg bool
	// LogStreamSendMsg if true, log out the contents of each sent stream
	// message.
	LogStreamSendMsg bool
	// LogUnaryReqMsg if true, log out the contents of the request
	// message/argument/parameters.
	LogUnaryReqMsg bool
	// LogAuthHeader if true, log out the contents of the authorization request
	// header.
	LogAuthHeader bool
	// LogCookies if true, log out the contents of any cookie header.
	LogCookies bool
	// IgnoreHeaders contains a list of request headers to filter from output.
	IgnoreHeaders []string
	// Default fields that should be included with all log messages. These can
	// be overwritten in your application.
	DefaultFields map[string]interface{}
}

// Interceptor implements gRPC interceptor middleware methods that logs the
// request as it comes in and the response as it goes out.
type Interceptor struct {
	cfg *Config
}

// New creates a new Validator interceptor.
func New(cfg *Config) interceptor.Middleware {
	return &Interceptor{
		cfg: cfg,
	}
}

// GetData returns the log field map from the request context.
func GetData(ctx context.Context) map[string]interface{} {
	var ok bool
	data, ok := ctx.Value(CtxKey{}).(map[string]interface{})
	if !ok {
		data = map[string]interface{}{}
	}
	return data
}

// UnaryInterceptor is a grpc interceptor middleware that logs out the request
// as it comes in, and the response as it goes out.
func (intr *Interceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Base fields
	dataFields := map[string]interface{}{}
	if nil != intr.cfg.DefaultFields {
		for k, v := range intr.cfg.DefaultFields {
			dataFields[k] = v
		}
	}
	dataFields["gateway-service"] = path.Dir(info.FullMethod)[1:]
	dataFields["gateway-method"] = path.Base(info.FullMethod)
	dataFields["gateway-request-type"] = "unary"
	dataFields["gateway-request-ip"] = getRequestIP(ctx, intr)

	// Request Payload Value
	if intr.cfg.LogUnaryReqMsg {
		if pb, ok := req.(proto.Message); ok {
			dataFields["gateway-request"] = pb
		}
	}

	// Add other fields and log the request started
	ctx = context.WithValue(ctx, CtxKey{}, dataFields)
	logRequest(ctx, intr, "request")

	// Call the handler
	resp, err := handler(ctx, req)

	// Calculate elapsed time and log the response
	// Re-extract the log fields, as they may have changed
	logResponse(ctx, intr, start, err, "response")

	// Return the response and error
	return resp, err
}

// StreamInterceptor is a grpc interceptor middleware that logs out the requests
// as they come in and the responses as they go out.
func (intr *Interceptor) StreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()

	// Get the wrapped server stream in order to access any modified context
	// from other interceptors
	wrapped := grpc_middleware.WrapServerStream(stream)
	ctx := wrapped.Context()

	// Base fields
	dataFields := map[string]interface{}{}
	if nil != intr.cfg.DefaultFields {
		for k, v := range intr.cfg.DefaultFields {
			dataFields[k] = v
		}
	}
	dataFields["gateway-service"] = path.Dir(info.FullMethod)[1:]
	dataFields["gateway-method"] = path.Base(info.FullMethod)
	dataFields["gateway-request-type"] = "stream"
	dataFields["gateway-request-ip"] = getRequestIP(ctx, intr)

	wrapped.WrappedContext = context.WithValue(ctx, CtxKey{}, dataFields)

	// Grap a log entry with just the base fields, for each streaming
	// send/receive
	streamEntry := log.WithFields(log.Fields(dataFields))

	// Add other fields and log the request started
	ctx = context.WithValue(ctx, CtxKey{}, dataFields)
	logRequest(ctx, intr, "request")

	// Call the handler
	err := handler(srv, &loggingServerStream{ServerStream: wrapped, entry: streamEntry, intr: intr})

	// Calculate elapsed time and log the response
	// Re-extract the log fields, as they may have changed
	logResponse(wrapped.Context(), intr, start, err, "response")

	// Return the error
	return err
}

// getRequestIP loads request metadata from the context.
func getRequestIP(ctx context.Context, intr *Interceptor) string {
	ip := "0.0.0.0"

	// metadata and HTTP headers.
	md := getRequestFields(ctx, intr)

	// try to find the real client IP.
	if _, ok := md["http-request-x-forwarded-for"]; ok {
		ip = md["http-request-x-forwarded-for"].([]string)[0]
	} else if _, ok := md["http-request-x-real-ip"]; ok {
		ip = md["http-request-x-real-ip"].([]string)[0]
	} else if _, ok := md["http-request-proxy-client-ip"]; ok {
		ip = md["http-request-proxy-client-ip"].([]string)[0]
	} else if _, ok := md["http-request-x-forwarded"]; ok {
		ip = md["http-request-x-forwarded"].([]string)[0]
	} else if _, ok := md["http-request-x-cluster-client-ip"]; ok {
		ip = md["http-request-x-cluster-client-ip"].([]string)[0]
	} else if _, ok := md["http-request-client-ip"]; ok {
		ip = md["http-request-client-ip"].([]string)[0]
	} else if _, ok := md["http-request-forwarded-for"]; ok {
		ip = md["http-request-forwarded-for"].([]string)[0]
	} else if _, ok := md["http-request-forwarded"]; ok {
		ip = md["http-request-forwarded"].([]string)[0]
	} else if _, ok := md["http-request-via"]; ok {
		ip = md["http-request-via"].([]string)[0]
	} else if _, ok := md["http-request-remote-addr"]; ok {
		ip = md["http-request-remote-addr"].([]string)[0]
	}

	// peer address
	if peerAddr, ok := peer.FromContext(ctx); ok {
		address := peerAddr.Addr.String()
		if address != "" &&
			!strings.HasPrefix(address, "127.0.0.1") &&
			!strings.HasPrefix(address, "localhost") {
			// strip the port and any brackets (IPv6)
			address = strings.TrimFunc(
				address[:strings.LastIndexByte(address, byte(':'))],
				func(r rune) bool {
					return '[' == r || ']' == r
				},
			)
			if ip == "0.0.0.0" {
				ip = address
			}
		}
	}

	ip = strings.Split(strings.Replace(ip, " ", "", -1), ",")[0]

	return ip
}

// getRequestFields loads request metadata from the context.
func getRequestFields(ctx context.Context, intr *Interceptor) map[string]interface{} {
	var clientID string

	requestFields := map[string]interface{}{}

	// requestID middleware.
	if requestID := middleware.GetReqID(ctx); "" != requestID {
		requestFields["gateway-request-id"] = requestID
	}

	// metadata and HTTP headers.
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			if strings.Contains(k, "cookie") {
				clientID = fmt.Sprintf("%s%s", clientID, v)
			}
			if strings.Contains(k, "authorization") {
				clientID = fmt.Sprintf("%s%s", clientID, v)
			}

			k = strings.ToLower(strings.Replace(k, ":", "_", -1))
			ignore := false
			for _, header := range intr.cfg.IgnoreHeaders {
				if k == header {
					ignore = true
					break
				}
			}
			if !ignore {
				if strings.Contains(k, "authorization") && intr.cfg.LogAuthHeader {
					requestFields[k] = v
				} else if strings.Contains(k, "cookie") && intr.cfg.LogCookies {
					requestFields[k] = v
				} else if !strings.Contains(k, "authorization") && !strings.Contains(k, "cookie") {
					requestFields[k] = v
				}
			}
		}

		// try to uniquely-ish identify clients using the user-agent and ip address
		clientID = fmt.Sprintf("%s%s", clientID, requestFields["gateway-request-ip"])
		if v, ok := requestFields["user-agent"]; ok {
			clientID = fmt.Sprintf("%s%s", clientID, v)
		}
		if "" != clientID {
			hash := sha1.New()
			_, _ = hash.Write([]byte(clientID))
			requestFields["gateway-client-id"] = base64.URLEncoding.EncodeToString(hash.Sum(nil))
		}
	}

	// peer address
	if peerAddr, ok := peer.FromContext(ctx); ok {
		address := peerAddr.Addr.String()
		if address != "" &&
			!strings.HasPrefix(address, "127.0.0.1") &&
			!strings.HasPrefix(address, "localhost") {
			// strip the port and any brackets (IPv6)
			address = strings.TrimFunc(
				address[:strings.LastIndexByte(address, byte(':'))],
				func(r rune) bool {
					return '[' == r || ']' == r
				},
			)
			requestFields["peer"] = address
		}
	}

	return requestFields
}

// logRequest adds request metadata, and then will log out the request access at
// info level.
func logRequest(ctx context.Context, intr *Interceptor, msg string) {
	var dataFields map[string]interface{}
	var ok bool

	if dataFields, ok = ctx.Value(CtxKey{}).(map[string]interface{}); !ok {
		dataFields = map[string]interface{}{}
	}

	log.WithField("request", getRequestFields(ctx, intr)).WithFields(log.Fields(dataFields)).Debug(msg)
}

// marshaller is the marshaller used for serializing protobuf messages.
var marshaller = &jsonpb.Marshaler{
	EmitDefaults: true,
	OrigName:     true,
}

// CtxKey is the key to use to lookup the log field map in the context.
type CtxKey struct{}

// logResponse calculates the elapsed time and the status code, and then
// will log out the response has finished at an appropriate level.
func logResponse(ctx context.Context, intr *Interceptor, start time.Time, err error, msg string) {
	var dataFields map[string]interface{}
	var ok bool

	if dataFields, ok = ctx.Value(CtxKey{}).(map[string]interface{}); !ok {
		dataFields = map[string]interface{}{}
	}

	// Calculate the elapsed time
	dataFields["elapsed"] = float64(time.Since(start).Nanoseconds()) / 1000000000
	dataFields["start"] = start.Format(time.RFC3339Nano)

	// Response code
	code := status.Code(err)
	dataFields["code"] = code

	// Response message
	message := ""
	if status, ok := status.FromError(err); ok {
		message = status.Message()
	}
	dataFields["msg"] = message

	// Log the response
	levelLog(log.WithField("request", getRequestFields(ctx, intr)).WithFields(log.Fields(dataFields)), DefaultCodeToLevel(code), msg)
}

// jsonpbMarshaler lets a proto interface be marshalled into json
type jsonpbMarshaler struct {
	proto.Message
}

// MarshalJSON lets jsonpbMarshaler implement json interface
func (j *jsonpbMarshaler) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	if err := marshaller.Marshal(b, j.Message); err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}
	return b.Bytes(), nil
}

// loggingServerStream wraps a ServerStream in order to log each send and
// receive.
type loggingServerStream struct {
	grpc.ServerStream
	entry *log.Entry
	intr  *Interceptor
}

// SendMsg lets loggingServerStream implement ServerStream, and will log sends.
func (l *loggingServerStream) SendMsg(m interface{}) error {
	err := l.ServerStream.SendMsg(m)
	if l.intr.cfg.LogStreamSendMsg {
		logProtoMessageAsJSON(l.entry, m, status.Code(err), "value", "StreamSend")
	}
	return err
}

// RecvMsg lets loggingServerStream implement ServerStream, and will log
// receives.
func (l *loggingServerStream) RecvMsg(m interface{}) error {
	err := l.ServerStream.RecvMsg(m)
	if l.intr.cfg.LogStreamRecvMsg {
		logProtoMessageAsJSON(l.entry, m, status.Code(err), "value", "StreamRecv")
	}
	return err
}

// logProtoMessageAsJSON logs an incoming or outgoing protobuf message as JSON.
func logProtoMessageAsJSON(
	entry *log.Entry,
	pbMsg interface{},
	code codes.Code,
	key string,
	msg string,
) {
	if p, ok := pbMsg.(proto.Message); ok {
		levelLog(entry.WithFields(log.Fields{key: &jsonpbMarshaler{p}, "code": code}), DefaultCodeToLevel(code), msg)
	} else {
		levelLog(entry.WithField("code", code), DefaultCodeToLevel(code), msg)
	}
}

// levelLog logs an entry and message at the appropriate levell
func levelLog(entry *log.Entry, level std.Level, msg string) {
	switch level {
	case log.DebugLevel:
		entry.Debug(msg)
	case log.InfoLevel:
		entry.Info(msg)
	case log.WarnLevel:
		entry.Warning(msg)
	case log.ErrorLevel:
		entry.Error(msg)
	case log.FatalLevel:
		entry.Fatal(msg)
	case log.PanicLevel:
		entry.Panic(msg)
	}
}

// DefaultCodeToLevel is the default implementation of gRPC return codes to log
// levels for server side.
func DefaultCodeToLevel(code codes.Code) std.Level {
	switch code {
	case codes.OK:
		return log.InfoLevel
	case codes.Canceled:
		return log.InfoLevel
	case codes.InvalidArgument:
		return log.InfoLevel
	case codes.NotFound:
		return log.InfoLevel
	case codes.AlreadyExists:
		return log.InfoLevel
	case codes.Unauthenticated:
		return log.InfoLevel

	case codes.DeadlineExceeded:
		return log.WarnLevel
	case codes.PermissionDenied:
		return log.WarnLevel
	case codes.ResourceExhausted:
		return log.WarnLevel
	case codes.FailedPrecondition:
		return log.WarnLevel
	case codes.Aborted:
		return log.WarnLevel
	case codes.OutOfRange:
		return log.WarnLevel
	case codes.Unavailable:
		return log.WarnLevel

	case codes.Unknown:
		return log.ErrorLevel
	case codes.Unimplemented:
		return log.ErrorLevel
	case codes.Internal:
		return log.ErrorLevel
	case codes.DataLoss:
		return log.ErrorLevel
	default:
		return log.ErrorLevel
	}
}
