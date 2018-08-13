<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty.png" width="750px">
</p>

# log

<p align="center">
	<a href="https://github.com/bdlm/log/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT"></a>
	<a href="https://github.com/mkenney/software-guides/blob/master/STABILITY-BADGES.md#release-candidate"><img src="https://img.shields.io/badge/stability-pre--release-48c9b0.svg" alt="Release Candidate"></a>
	<a href="https://travis-ci.org/bdlm/log"><img src="https://travis-ci.org/bdlm/log.svg?branch=master" alt="Build status"></a>
	<a href="https://codecov.io/gh/bdlm/log"><img src="https://img.shields.io/codecov/c/github/bdlm/log/master.svg" alt="Coverage status"></a>
	<a href="https://goreportcard.com/report/github.com/bdlm/log"><img src="https://goreportcard.com/badge/github.com/bdlm/log" alt="Go Report Card"></a>
	<a href="https://github.com/bdlm/log/issues"><img src="https://img.shields.io/github/issues-raw/bdlm/log.svg" alt="Github issues"></a>
	<a href="https://github.com/bdlm/log/pulls"><img src="https://img.shields.io/github/issues-pr/bdlm/log.svg" alt="Github pull requests"></a>
	<a href="https://godoc.org/github.com/bdlm/log"><img src="https://godoc.org/github.com/bdlm/log?status.svg" alt="GoDoc"></a>
</p>

`bdlm/log` is a fork of the excellent [`sirupsen/logrus`](https://github.com/sirupsen/logrus) package. This package adds:

* support for sanitizing strings from log output to aid in preventing leaking sensitive data.
* additional default fields `host` and `caller`.
* support for suppressing any default field.
* TTY formatting and coloring of JSON output.
* updated text formatting for TTY output.
* updated TTY color scheme.

#

`bdlm/log` is a structured logger for Go and is API compatible with the standard libaray `log` package.

## Examples

Please see the [user documentation](https://github.com/bdlm/log/wiki) for full details.

### Simple usage

The simplest way to use `bdlm/log` is with the exported package logger:

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

Note that it is fully api-compatible with the stdlib logger, so you can replace your `log` imports everywhere or using a strangler pattern with `"github.com/bdlm/log"` and add the full logging flexibility to your service without impacting existing code.

## Formats

By default, `bdlm/log` uses a basic text format:
```sh
time="2018-08-12T20:36:34.577-06:00" level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1 caller="main.go:42 main.main" host="myhost"
time="2018-08-12T20:36:34.578-06:00" level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20 caller="main.go:46 main.main" host="myhost"
time="2018-08-12T20:36:34.578-06:00" level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100 caller="main.go:50 main.main" host="myhost"
time="2018-08-12T20:36:34.578-06:00" level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.count="wait, what?" caller="main.go:54 main.main" host="myhost"
time="2018-08-12T20:36:34.579-06:00" level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true caller="main.go:58 main.main" host="myhost"
time="2018-08-12T20:36:34.579-06:00" level="fatal" msg="That could have gone better..." data.dead=true data.winner="walrus" caller="main.go:33 main.main.func1" host="myhost"
```

For development, color-coded output formated for humans is automatically enabled when a `tty` terminal is detected (this can be disabled with `log.SetFormatter(&log.TextFormatter{DisableTTY: true})`):

<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty.png" width="750px">
</p>

JSON formatting is also available with `log.SetFormatter(&log.JSONFormatter{})` for easy parsing by logstash or similar:

```json
{"caller":"main.go:42 main.main","data":{"animal":"bird","count":1},"host":"myhost","level":"debug","msg":"Oh, look, a bird...","time":"2018-08-12T20:38:03.997-06:00"}
{"caller":"main.go:46 main.main","data":{"animal":"walrus","count":20},"host":"myhost","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-12T20:38:03.998-06:00"}
{"caller":"main.go:50 main.main","data":{"animal":"walrus","count":100},"host":"myhost","level":"warn","msg":"The group's number increased tremendously!","time":"2018-08-12T20:38:03.998-06:00"}
{"caller":"main.go:54 main.main","data":{"animal":"cow","count":"wait, what?"},"host":"myhost","level":"error","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-12T20:38:03.998-06:00"}
{"caller":"main.go:58 main.main","data":{"animal":"walrus","run":true},"host":"myhost","level":"panic","msg":"The walrus are attacking!","time":"2018-08-12T20:38:03.999-06:00"}
{"caller":"main.go:33 main.main.func1","data":{"dead":true,"winner":"walrus"},"host":"myhost","level":"fatal","msg":"That could have gone better...","time":"2018-08-12T20:38:03.999-06:00"}
```

The JSON formatter also makes adjustments by default when a `tty` terminal is detected and can be disabled similarly with `log.SetFormatter(&log.JSONFormatter{DisableTTY: true})`:

<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-json.png" width="400px">
</p>

## Trace-level logging

The default formatters also have a `trace` mode that is disabled by default. Rather than acting as an additional log level, it is instead a verbose mode that includes the full backtrace of the call that triggered the log write. Here are the above examples with trace enabled:

### TextFormat

To enable trace output:
```go
log.SetFormatter(&log.TextFormatter{EnableTrace: true})
```

Non-TTY trace output:
```sh
time="2018-08-12T20:40:59.258-06:00" level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1 caller="main.go:42 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace" trace.1="formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData" trace.2="text_formatter.go:93 github.com/bdlm/test/vendor/github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.5="entry.go:196 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Debug" trace.6="main.go:42 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-12T20:40:59.258-06:00" level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20 caller="main.go:46 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace" trace.1="formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData" trace.2="text_formatter.go:93 github.com/bdlm/test/vendor/github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.5="entry.go:203 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Info" trace.6="main.go:46 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-12T20:40:59.258-06:00" level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100 caller="main.go:50 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace" trace.1="formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData" trace.2="text_formatter.go:93 github.com/bdlm/test/vendor/github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.5="entry.go:215 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Warn" trace.6="main.go:50 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-12T20:40:59.259-06:00" level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.count="wait, what?" caller="main.go:54 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace" trace.1="formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData" trace.2="text_formatter.go:93 github.com/bdlm/test/vendor/github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.5="entry.go:227 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Error" trace.6="main.go:54 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-12T20:40:59.259-06:00" level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true caller="main.go:58 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace" trace.1="formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData" trace.2="text_formatter.go:93 github.com/bdlm/test/vendor/github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.5="entry.go:242 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Panic" trace.6="main.go:58 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-12T20:40:59.259-06:00" level="fatal" msg="That could have gone better..." data.dead=true data.winner="walrus" caller="main.go:33 main.main.func1" host="myhost" trace.0="formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace" trace.1="formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData" trace.2="text_formatter.go:93 github.com/bdlm/test/vendor/github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.5="entry.go:234 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Fatal" trace.6="main.go:33 main.main.func1" trace.7="asm_amd64.s:573 runtime.call32" trace.8="panic.go:505 runtime.gopanic" trace.9="entry.go:155 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log" trace.10="entry.go:242 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Panic" trace.11="main.go:58 main.main" trace.12="proc.go:198 runtime.main" trace.13="asm_amd64.s:2361 runtime.goexit"
```

TTY trace output:
<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-trace.png" width="750px">
</p>

### JSONFormat

To enable trace output:
```go
log.SetFormatter(&log.JSONFormatter{EnableTrace: true})
```

Non-TTY trace output:
```json
{"caller":"main.go:42 main.main","data":{"animal":"bird","count":1},"host":"myhost","level":"debug","msg":"Oh, look, a bird...","time":"2018-08-12T20:43:17.410-06:00","trace":["formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace","formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData","json_formatter.go:76 github.com/bdlm/test/vendor/github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:196 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Debug","main.go:42 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:46 main.main","data":{"animal":"walrus","count":20},"host":"myhost","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-12T20:43:17.411-06:00","trace":["formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace","formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData","json_formatter.go:76 github.com/bdlm/test/vendor/github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:203 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Info","main.go:46 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:50 main.main","data":{"animal":"walrus","count":100},"host":"myhost","level":"warn","msg":"The group's number increased tremendously!","time":"2018-08-12T20:43:17.411-06:00","trace":["formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace","formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData","json_formatter.go:76 github.com/bdlm/test/vendor/github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:215 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Warn","main.go:50 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:54 main.main","data":{"animal":"cow","count":"wait, what?"},"host":"myhost","level":"error","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-12T20:43:17.411-06:00","trace":["formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace","formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData","json_formatter.go:76 github.com/bdlm/test/vendor/github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:227 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Error","main.go:54 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:58 main.main","data":{"animal":"walrus","run":true},"host":"myhost","level":"panic","msg":"The walrus are attacking!","time":"2018-08-12T20:43:17.411-06:00","trace":["formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace","formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData","json_formatter.go:76 github.com/bdlm/test/vendor/github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:242 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Panic","main.go:58 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:33 main.main.func1","data":{"dead":true,"winner":"walrus"},"host":"myhost","level":"fatal","msg":"That could have gone better...","time":"2018-08-12T20:43:17.411-06:00","trace":["formatter.go:83 github.com/bdlm/test/vendor/github.com/bdlm/log.getTrace","formatter.go:156 github.com/bdlm/test/vendor/github.com/bdlm/log.getData","json_formatter.go:76 github.com/bdlm/test/vendor/github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:234 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Fatal","main.go:33 main.main.func1","asm_amd64.s:573 runtime.call32","panic.go:505 runtime.gopanic","entry.go:155 github.com/bdlm/test/vendor/github.com/bdlm/log.Entry.log","entry.go:242 github.com/bdlm/test/vendor/github.com/bdlm/log.(*Entry).Panic","main.go:58 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
```

TTY trace output:
<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-json-trace.png" width="750px">
</p>
