# log

<p align="center">
	<a href="https://github.com/bdlm/log/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT"></a>
	<a href="https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#mature"><img src="https://img.shields.io/badge/stability-mature-008000.svg" alt="Mature"></a>
	<a href="https://travis-ci.org/bdlm/log"><img src="https://travis-ci.org/bdlm/log.svg?branch=master" alt="Build status"></a>
	<a href="https://codecov.io/gh/bdlm/log"><img src="https://img.shields.io/codecov/c/github/bdlm/log/master.svg" alt="Coverage status"></a>
	<a href="https://goreportcard.com/report/github.com/bdlm/log"><img src="https://goreportcard.com/badge/github.com/bdlm/log" alt="Go Report Card"></a>
	<a href="https://github.com/bdlm/log/issues"><img src="https://img.shields.io/github/issues-raw/bdlm/log.svg" alt="Github issues"></a>
	<a href="https://github.com/bdlm/log/pulls"><img src="https://img.shields.io/github/issues-pr/bdlm/log.svg" alt="Github pull requests"></a>
	<a href="https://godoc.org/github.com/bdlm/log"><img src="https://godoc.org/github.com/bdlm/log?status.svg" alt="GoDoc"></a>
</p>

`bdlm/log` is a fork of the excellent [`sirupsen/logrus`](https://github.com/sirupsen/logrus) package.

<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-header.png" width="75%">
</p>

This package adds:

* support for sanitizing strings from log output to aid in preventing leaking sensitive data.
* additional default fields `host` and `caller`.
* verbose output including the full backtrace of logger calls.
* support for suppressing any default field.
* TTY formatting and coloring of JSON output.
* updated formatting for TTY text output.
* updated default TTY color scheme and color customization.

#

`bdlm/log` is a structured logger for Go and is fully API compatible with both the standard libaray `log` package as well as the [`sirupsen/logrus`](https://github.com/sirupsen/logrus) package.

## Examples

#### [User documentation](https://github.com/bdlm/log/wiki)

### Simple usage

The simplest way to use `bdlm/log` is with the exported package-level logger:

```go
package main

import (
    log "github.com/bdlm/log"
)

func main() {
    log.WithFields(log.Fields{
        "animal": "walrus",
    }).Info("A walrus appears")
}
```

## Compatibility

Note that `bdlm/log` is fully api-compatible with the stdlib logger, so you can replace your `log` imports everywhere or using a strangler pattern with `"github.com/bdlm/log"` and add the full logging flexibility to your service without impacting existing code.

The default log format of this package does not match the stdlib logger's default output so a compatible formatter, `StdFormatter`, is provided which includes the additional information inline:

```go
log.SetFormatter(&log.StdFormatter{})
```

Which results in the following log output. `StdFormatter` does not have a TTY feature:
```
2018/08/17 20:17:45 Oh, look, a bird... level="debug" data.animal="bird" data.count=1 caller="main.go:39 main.main" host="myhost"
2018/08/17 20:17:45 A group of walrus emerges from the ocean level="info" data.animal="walrus" data.count=20 caller="main.go:43 main.main" host="myhost"
2018/08/17 20:17:45 The group's number increased tremendously! level="warn" data.animal="walrus" data.count=100 caller="main.go:47 main.main" host="myhost"
2018/08/17 20:17:45 Tremendously sized cow enters the ocean. level="error" data.animal="cow" data.count="wait, what?" caller="main.go:51 main.main" host="myhost"
2018/08/17 20:17:45 The walrus are attacking! level="panic" data.animal="walrus" data.run=true caller="main.go:55 main.main" host="myhost"
2018/08/17 20:17:45 That could have gone better... level="fatal" data.dead=true data.winner="walrus" caller="main.go:30 main.main.func1" host="myhost"
```

## Log Formatters

By default, `bdlm/log` uses a basic text format:
```sh
time="2018-08-17T18:28:07.385-06:00" level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1 caller="main.go:34 main.main" host="myhost"
time="2018-08-17T18:28:07.385-06:00" level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20 caller="main.go:38 main.main" host="myhost"
time="2018-08-17T18:28:07.385-06:00" level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100 caller="main.go:42 main.main" host="myhost"
time="2018-08-17T18:28:07.385-06:00" level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.count="wait, what?" caller="main.go:46 main.main" host="myhost"
time="2018-08-17T18:28:07.385-06:00" level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true caller="main.go:50 main.main" host="myhost"
time="2018-08-17T18:28:07.385-06:00" level="fatal" msg="That could have gone better..." data.dead=true data.winner="walrus" caller="main.go:25 main.main.func1" host="myhost"
```

For development, color-coded output formated for humans is automatically enabled when a `tty` terminal is detected (this can be disabled with `log.SetFormatter(&log.TextFormatter{DisableTTY: true})`):

<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty.png" width="50%">
</p>

JSON formatting is also available with `log.SetFormatter(&log.JSONFormatter{})` for easy parsing by logstash or similar:

```json
{"caller":"main.go:36 main.main","data":{"animal":"bird","count":1},"host":"myhost","level":"debug","msg":"Oh, look, a bird...","time":"2018-08-17T18:32:30.786-06:00"}
{"caller":"main.go:40 main.main","data":{"animal":"walrus","count":20},"host":"myhost","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-17T18:32:30.786-06:00"}
{"caller":"main.go:44 main.main","data":{"animal":"walrus","count":100},"host":"myhost","level":"warn","msg":"The group's number increased tremendously!","time":"2018-08-17T18:32:30.786-06:00"}
{"caller":"main.go:48 main.main","data":{"animal":"cow","count":"wait, what?"},"host":"myhost","level":"error","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-17T18:32:30.786-06:00"}
{"caller":"main.go:52 main.main","data":{"animal":"walrus","run":true},"host":"myhost","level":"panic","msg":"The walrus are attacking!","time":"2018-08-17T18:32:30.786-06:00"}
{"caller":"main.go:27 main.main.func1","data":{"dead":true,"winner":"walrus"},"host":"myhost","level":"fatal","msg":"That could have gone better...","time":"2018-08-17T18:32:30.787-06:00"}
```

The JSON formatter also makes adjustments by default when a `tty` terminal is detected and can be disabled similarly with `log.SetFormatter(&log.JSONFormatter{DisableTTY: true})`:

<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-json.png" width="50%">
</p>

## Backtrace data

The standard formatters also have a `trace` mode that is disabled by default. Rather than acting as an additional log level, it is instead a verbose mode that includes the full backtrace of the call that triggered the log write. To enable trace output, set `EnableTrace` to `true`.

Here are the above examples with trace enabled:

### TextFormat

Non-TTY trace output:
```go
log.SetFormatter(&log.TextFormatter{
    DisableTTY: true,
    EnableTrace: true,
})
```

```sh
time="2018-08-18T00:20:36.468-06:00" level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1 host="myhost" trace.0="main.go:38 main.main"
time="2018-08-18T00:20:36.469-06:00" level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20 host="myhost" trace.0="main.go:42 main.main"
time="2018-08-18T00:20:36.469-06:00" level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100 host="myhost" trace.0="main.go:46 main.main"
time="2018-08-18T00:20:36.469-06:00" level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.count="wait, what?" host="myhost" trace.0="main.go:50 main.main"
time="2018-08-18T00:20:36.469-06:00" level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true host="myhost" trace.0="main.go:54 main.main"
time="2018-08-18T00:20:36.469-06:00" level="fatal" msg="That could have gone better..." data.dead=true data.winner="walrus" host="myhost" trace.0="main.go:30 main.main.func1" trace.1="asm_amd64.s:573 runtime.call32" trace.2="panic.go:502 runtime.gopanic" trace.3="main.go:54 main.main"
```

TTY trace output:
```go
log.SetFormatter(&log.TextFormatter{
    EnableTrace: true,
    ForceTTY: true,
})
```
<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-trace.png" width="50%">
</p>

### JSONFormat

To enable trace output:
```go
log.SetFormatter(&log.JSONFormatter{
    DisableTTY: true,
    EnableTrace: true,
})
```

Non-TTY trace output:
```json
{"caller":"main.go:38 main.main","data":{"animal":"bird","count":1},"host":"myhost","level":"debug","msg":"Oh, look, a bird...","time":"2018-08-18T00:22:16.057-06:00","trace":["main.go:38 main.main"]}
{"caller":"main.go:42 main.main","data":{"animal":"walrus","count":20},"host":"myhost","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-18T00:22:16.058-06:00","trace":["main.go:42 main.main"]}
{"caller":"main.go:46 main.main","data":{"animal":"walrus","count":100},"host":"myhost","level":"warn","msg":"The group's number increased tremendously!","time":"2018-08-18T00:22:16.058-06:00","trace":["main.go:46 main.main"]}
{"caller":"main.go:50 main.main","data":{"animal":"cow","count":"wait, what?"},"host":"myhost","level":"error","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-18T00:22:16.058-06:00","trace":["main.go:50 main.main"]}
{"caller":"main.go:54 main.main","data":{"animal":"walrus","run":true},"host":"myhost","level":"panic","msg":"The walrus are attacking!","time":"2018-08-18T00:22:16.058-06:00","trace":["main.go:54 main.main"]}
{"caller":"main.go:30 main.main.func1","data":{"dead":true,"winner":"walrus"},"host":"myhost","level":"fatal","msg":"That could have gone better...","time":"2018-08-18T00:22:16.058-06:00","trace":["main.go:30 main.main.func1","asm_amd64.s:573 runtime.call32","panic.go:502 runtime.gopanic","main.go:54 main.main"]}
```

TTY trace output:
```go
log.SetFormatter(&log.JSONFormatter{
    EnableTrace: true,
    ForceTTY: true,
})
```
<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-json-trace.png" width="50%">
</p>
