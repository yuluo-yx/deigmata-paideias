## 1. 准备 go 应用


## 2. 加入 otel 依赖

```golang
go get "go.opentelemetry.io/otel" \
  "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric" \
  "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" \
  "go.opentelemetry.io/otel/propagation" \
  "go.opentelemetry.io/otel/sdk/metric" \
  "go.opentelemetry.io/otel/sdk/resource" \
  "go.opentelemetry.io/otel/sdk/trace" \
  "go.opentelemetry.io/otel/semconv/v1.24.0" \
  "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
```

## 3. 初始化 otel sdk 修改代码

```golang
package main

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*trace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}
```

## 4. 构建修改之后的应用程序

```sell
go mod tidy

export OTEL_RESOURCE_ATTRIBUTES="service.name=dice,service.version=0.1.0"

go run .
```

## 5. output

```json
{
        "Name": "/",
        "SpanContext": {
                "TraceID": "e19e498cb32617b5c0d8edb58fcd6b8d",
                "SpanID": "d444f5579cd68b4a",
                "TraceFlags": "01",
                "TraceState": "",
                "Remote": false
        },
        "Parent": {
                "TraceID": "00000000000000000000000000000000",
                "SpanID": "0000000000000000",
                "TraceFlags": "00",
                "TraceState": "",
                "Remote": false
        },
        "SpanKind": 2,
        "StartTime": "2024-04-18T12:30:39.2396583+08:00",
        "EndTime": "2024-04-18T12:30:39.2401623+08:00",
        "Attributes": [
                {
                        "Key": "http.method",
                        "Value": {
                                "Type": "STRING",
                                "Value": "GET"
                        }
                },
                {
                        "Key": "http.scheme",
                        "Value": {
                                "Type": "STRING",
                                "Value": "http"
                        }
                },
                {
                        "Key": "net.host.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "localhost"
                        }
                },
                {
                        "Key": "net.host.port",
                        "Value": {
                                "Type": "INT64",
                                "Value": 8080
                        }
                },
                {
                        "Key": "net.sock.peer.addr",
                        "Value": {
                                "Type": "STRING",
                                "Value": "::1"
                        }
                },
                {
                        "Key": "net.sock.peer.port",
                        "Value": {
                                "Type": "INT64",
                                "Value": 12999
                        }
                },
                {
                        "Key": "user_agent.original",
                        "Value": {
                                "Type": "STRING",
                                "Value": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
                        }
                },
                {
                        "Key": "http.target",
                        "Value": {
                                "Type": "STRING",
                                "Value": "/rolldice"
                        }
                },
                {
                        "Key": "net.protocol.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.1"
                        }
                },
                {
                        "Key": "http.route",
                        "Value": {
                                "Type": "STRING",
                                "Value": "/rolldice"
                        }
                },
                {
                        "Key": "http.wrote_bytes",
                        "Value": {
                                "Type": "INT64",
                                "Value": 2
                        }
                },
                {
                        "Key": "http.status_code",
                        "Value": {
                                "Type": "INT64",
                                "Value": 200
                        }
                }
        ],
        "Events": null,
        "Links": null,
        "Status": {
                "Code": "Unset",
                "Description": ""
        },
        "DroppedAttributes": 0,
        "DroppedEvents": 0,
        "DroppedLinks": 0,
        "ChildSpanCount": 1,
        "Resource": [
                {
                        "Key": "service.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "dice"
                        }
                },
                {
                        "Key": "service.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "0.1.0"
                        }
                },
                {
                        "Key": "telemetry.sdk.language",
                        "Value": {
                                "Type": "STRING",
                                "Value": "go"
                        }
                },
                {
                        "Key": "telemetry.sdk.name",
                        "Value": {
                                "Type": "STRING",
                                "Value": "opentelemetry"
                        }
                },
                {
                        "Key": "telemetry.sdk.version",
                        "Value": {
                                "Type": "STRING",
                                "Value": "1.25.0"
                        }
                }
        ],
        "InstrumentationLibrary": {
                "Name": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
                "Version": "0.50.0",
                "SchemaURL": ""
        }
}

```

Result:

```json
{"Resource":[{"Key":"service.name","Value":{"Type":"STRING","Value":"dice"}},{"Key":"service.version","Value":{"Type":
"STRING","Value":"0.1.0"}},{"Key":"telemetry.sdk.language","Value":{"Type":"STRING","Value":"go"}},{"Key":"telemetry.s
  dk.name","Value":{"Type":"STRING","Value":"opentelemetry"}},{"Key":"telemetry.sdk.version","Value":{"Type":"STRING","V
  alue":"1.25.0"}}],"ScopeMetrics":[{"Scope":{"Name":"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp","Ve
  rsion":"0.50.0","SchemaURL":""},"Metrics":[{"Name":"http.server.request.size","Description":"Measures the size of HTTP
  request messages.","Unit":"By","Data":{"DataPoints":[{"Attributes":[{"Key":"http.method","Value":{"Type":"STRING","Va
  lue":"GET"}},{"Key":"http.route","Value":{"Type":"STRING","Value":"/rolldice"}},{"Key":"http.scheme","Value":{"Type":"
  STRING","Value":"http"}},{"Key":"http.status_code","Value":{"Type":"INT64","Value":200}},{"Key":"net.host.name","Value
  ":{"Type":"STRING","Value":"localhost"}},{"Key":"net.host.port","Value":{"Type":"INT64","Value":8080}},{"Key":"net.pro
  tocol.name","Value":{"Type":"STRING","Value":"http"}},{"Key":"net.protocol.version","Value":{"Type":"STRING","Value":"
  1.1"}}],"StartTime":"2024-04-18T12:30:23.1387006+08:00","Time":"2024-04-18T12:30:41.146741+08:00","Value":0}],"Tempora
  lity":"CumulativeTemporality","IsMonotonic":true}},{"Name":"http.server.response.size","Description":"Measures the siz
  e of HTTP response messages.","Unit":"By","Data":{"DataPoints":[{"Attributes":[{"Key":"http.method","Value":{"Type":"S
  TRING","Value":"GET"}},{"Key":"http.route","Value":{"Type":"STRING","Value":"/rolldice"}},{"Key":"http.scheme","Value"
  :{"Type":"STRING","Value":"http"}},{"Key":"http.status_code","Value":{"Type":"INT64","Value":200}},{"Key":"net.host.na
  me","Value":{"Type":"STRING","Value":"localhost"}},{"Key":"net.host.port","Value":{"Type":"INT64","Value":8080}},{"Key
  ":"net.protocol.name","Value":{"Type":"STRING","Value":"http"}},{"Key":"net.protocol.version","Value":{"Type":"STRING"
,"Value":"1.1"}}],"StartTime":"2024-04-18T12:30:23.1387006+08:00","Time":"2024-04-18T12:30:41.146741+08:00","Value":2}
],"Temporality":"CumulativeTemporality","IsMonotonic":true}},{"Name":"http.server.duration","Description":"Measures th
e duration of inbound HTTP requests.","Unit":"ms","Data":{"DataPoints":[{"Attributes":[{"Key":"http.method","Value":{"
Type":"STRING","Value":"GET"}},{"Key":"http.route","Value":{"Type":"STRING","Value":"/rolldice"}},{"Key":"http.scheme"
,"Value":{"Type":"STRING","Value":"http"}},{"Key":"http.status_code","Value":{"Type":"INT64","Value":200}},{"Key":"net
.host.name","Value":{"Type":"STRING","Value":"localhost"}},{"Key":"net.host.port","Value":{"Type":"INT64","Value":8080
}},{"Key":"net.protocol.name","Value":{"Type":"STRING","Value":"http"}},{"Key":"net.protocol.version","Value":{"Type":
"STRING","Value":"1.1"}}],"StartTime":"2024-04-18T12:30:23.1387006+08:00","Time":"2024-04-18T12:30:41.146741+08:00","C
ount":1,"Bounds":[0,5,10,25,50,75,100,250,500,750,1000,2500,5000,7500,10000],"BucketCounts":[0,1,0,0,0,0,0,0,0,0,0,0,0
,0,0,0],"Min":0.504,"Max":0.504,"Sum":0.504}],"Temporality":"CumulativeTemporality"}}]},{"Scope":{"Name":"rolldice","V
ersion":"","SchemaURL":""},"Metrics":[{"Name":"dice.rolls","Description":"The number of rolls by roll value","Unit":"{
roll}","Data":{"DataPoints":[{"Attributes":[{"Key":"roll.value","Value":{"Type":"INT64","Value":2}}],"StartTime":"2024
-04-18T12:30:23.1387006+08:00","Time":"2024-04-18T12:30:41.146741+08:00","Value":1}],"Temporality":"CumulativeTemporality","IsMonotonic":true}}]}]}
```