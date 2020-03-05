The `log` package provides a gRPC request interceptor for logging purposes. It is extrememly helpful in creating useful log metrics for application analytics and includes a compatible Scalyr log parser.

All streaming and unary requests and responses are logged at the `debug` level and all responses (which include the request data) are logged at the `info` level. Both request and response data intercepted via the `grpc-gateway` package autmatically add all HTTP request headers to the `"request"` key. The service and method requested are also logged. Response logs also include the elapsed time and gRPC response code.

* [Configuration](#configuration)
* [Implementation](#implementation)
* [Adding Fields](#adding-fields)
* [Log Output](#log-output)
* [Scalyr Log Parser](#scalyr-log-parser)
  * [Scalyr Metrics](#scalyr-metrics)

## Configuration

The interceptor provides several configuration options but the defaults are fine for most use cases:

```go
// Interceptor contains gRPC interceptor middleware methods that logs the
// request as it comes in and the response as it goes out.
type Interceptor struct {
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
```

## Implementation

To impliment the log interceptor in your gRPC service simply add it to the chain of request interceptors when creating your server:

```go
import log_interceptor "github.com/bdlm/log/interceptor/log"

func main() {
	logInterceptor := log_interceptor.Interceptor{}

	// Create the gRPC server.
	grpcServer := grpc.NewServer(
		// Add the log interceptor to the chain of streaming request interceptors.
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			logInterceptor.StreamInterceptor,
		)),
		// Add the log interceptor to the chain of unary request interceptors.
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			logInterceptor.UnaryInterceptor, // automatically log requests
		)),
	)

	...
}
```

## Adding Fields

Custom fields can easily be added to the response log in your gRPC service, the data map is available via the request context.

<sub>* *Any custom fields written by your service will be added to the* `data` *map and will overwrite any conflicting fields.*</sub>

```go
import log_interceptor "github.com/bdlm/log/interceptor/log"

func (r RPC) HelloWorld(ctx context.Context, msg *pb.Hello) error {
	// Fetch the field map from the request context.
	logFields := ctx.Value(log_interceptor.CtxKey{}).(map[string]interface{})

	// All maps are pointers. Writes to this map within your service calls is
	// always concurrent-safe. Nothing further is required.
	logFields["hello"] = msg.GetName()

	// ~Fin~
	return nil
}
```

## Log Output

A typical (prettified) response log message might look like this:

<sub>* *Any custom fields written by your service will be added to the* `data` *map and will overwrite any conflicting fields.*</sub>

```json
{
    "level": "info",
    "host": "c5582e3db400",
    "time": "2019-03-15T21:38:53.292Z",
    "msg": "response",
    "data": {
        "code": 0,
        "elapsed": 0.5147447,
        "gateway-method": "HelloWorld",
        "gateway-request-type": "unary",
        "gateway-service": "rp.api.example.Example",
        "grpc-service": "HelloWorld",
        "hello": "World",
        "request": {
            "content-type": [
                "application/grpc"
            ],
            "gateway-client-id": "_awsRzwqCaQG2MoKtwNWKDw8ARA=",
            "gateway-request-ip": "172.19.0.1",
            "http-request-content-type": [
                "application/json"
            ],
            "http-request-postman-token": [
                "402db792-8437-4216-97d1-a3e2979e125d"
            ],
            "http-request-user-agent": [
                "PostmanRuntime/7.6.0"
            ],
            "http-request-x-forwarded-for": [
                "172.19.0.1"
            ],
            "http-request-x-forwarded-host": [
                "api-infrastructure-example.app.returnpath.net"
            ],
            "http-request-x-forwarded-port": [
                "80"
            ],
            "http-request-x-forwarded-proto": [
                "http"
            ],
            "http-request-x-forwarded-server": [
                "b000e92aec95"
            ],
            "http-request-x-real-ip": [
                "172.19.0.1"
            ],
            "user-agent": [
                "grpc-go/1.18.0"
            ],
            "x-forwarded-for": [
                "172.19.0.1, 172.19.0.13"
            ],
            "x-forwarded-host": [
                "api-infrastructure-example.app.returnpath.net"
            ]
        },
        "start": "2019-03-15T21:38:50.8128519Z"
    },
    "caller": "log.go:347 github.com/bdlm/log-example/vendor/github.com/bdlm/log/interceptor/log.levelLog"
}
```

## Scalyr Log Parser

The provided [Scalyr log parser](https://github.com/bdlm/log/blob/master/interceptor/log/scalyr-log-parser.js) will convert these fields into Scalyr variables available to log filters and dashboards:

```
$level = 'info' $msg = 'response' $DataHello = 'World'
```

It can be annoying to always need to type the `$Data` part of the variable names. You can rename the data key to anything you like in your service:

```go
import "github.com/bdlm/log"

func main() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			"data": "_",
		},
	})
}
```

Which the Scaylr parser will translate to:

```
$level = 'info' $msg = 'response' $_Hello = 'World'
```

### Scalyr Metrics

Structured logs are extremely useful in gathering actionable metrics from your application. To that end, response timings are automatically included in response logs and can be very useful in monitoring application performance:

```json
{
    "graphWidth": 1200,
    "graphHeight": 200,
    "duration": "1w",
    "graphs": [
        {
            "label": "Request performance (seconds)",
            "plots": [
                {
                    "label": "Overall",
                    "filter": "$serverHost contains 'prod' $parser = 'api-infrastructure-example' $level = 'info' $msg = 'response' $_Grpc-service = 'HelloWorld'",
                    "facet": "mean(_Elapsed)"
                }
            ]
        }
    ]
}
```
