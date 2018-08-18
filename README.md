<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty.png">
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
* verbose output including the full backtrace of logger calls.
* support for suppressing any default field.
* TTY formatting and coloring of JSON output.
* updated formatting for TTY text output.
* updated TTY color scheme.

#

`bdlm/log` is a structured logger for Go and is API compatible with the standard libaray `log` package.

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

The default log format of this package does not match the stdlib logger's default output so a compatible formatter, `STDFormatter`, is provided which includes the additional information inline:

```go
log.SetFormatter(&log.STDFormatter{})
```

Which results in a standard log output. `STDFormatter` does not have a separate TTY format:
```
2018/08/17 20:10:59 Oh, look, a bird... data.animal="bird" data.count=1 caller="main.go:39 main.main" host="myhost"
2018/08/17 20:10:59 A group of walrus emerges from the ocean data.animal="walrus" data.count=20 caller="main.go:43 main.main" host="myhost"
2018/08/17 20:10:59 The group's number increased tremendously! data.animal="walrus" data.count=100 caller="main.go:47 main.main" host="myhost"
2018/08/17 20:10:59 Tremendously sized cow enters the ocean. data.animal="cow" data.count="wait, what?" caller="main.go:51 main.main" host="myhost"
2018/08/17 20:10:59 The walrus are attacking! data.animal="walrus" data.run=true caller="main.go:55 main.main" host="myhost"
2018/08/17 20:10:59 That could have gone better... data.dead=true data.winner="walrus" caller="main.go:30 main.main.func1" host="myhost"
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
exit status 1
```

For development, color-coded output formated for humans is automatically enabled when a `tty` terminal is detected (this can be disabled with `log.SetFormatter(&log.TextFormatter{DisableTTY: true})`):

<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty.png">
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
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-json.png" width="400px">
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
time="2018-08-17T18:34:28.651-06:00" level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1 caller="main.go:37 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/log.getTrace" trace.1="formatter.go:153 github.com/bdlm/log.getData" trace.2="text_formatter.go:96 github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/log.Entry.log" trace.5="entry.go:196 github.com/bdlm/log.(*Entry).Debug" trace.6="main.go:37 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-17T18:34:28.652-06:00" level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20 caller="main.go:41 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/log.getTrace" trace.1="formatter.go:153 github.com/bdlm/log.getData" trace.2="text_formatter.go:96 github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/log.Entry.log" trace.5="entry.go:203 github.com/bdlm/log.(*Entry).Info" trace.6="main.go:41 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-17T18:34:28.652-06:00" level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100 caller="main.go:45 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/log.getTrace" trace.1="formatter.go:153 github.com/bdlm/log.getData" trace.2="text_formatter.go:96 github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/log.Entry.log" trace.5="entry.go:215 github.com/bdlm/log.(*Entry).Warn" trace.6="main.go:45 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-17T18:34:28.652-06:00" level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.count="wait, what?" caller="main.go:49 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/log.getTrace" trace.1="formatter.go:153 github.com/bdlm/log.getData" trace.2="text_formatter.go:96 github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/log.Entry.log" trace.5="entry.go:227 github.com/bdlm/log.(*Entry).Error" trace.6="main.go:49 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-17T18:34:28.652-06:00" level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true caller="main.go:53 main.main" host="myhost" trace.0="formatter.go:83 github.com/bdlm/log.getTrace" trace.1="formatter.go:153 github.com/bdlm/log.getData" trace.2="text_formatter.go:96 github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/log.Entry.log" trace.5="entry.go:242 github.com/bdlm/log.(*Entry).Panic" trace.6="main.go:53 main.main" trace.7="proc.go:198 runtime.main" trace.8="asm_amd64.s:2361 runtime.goexit"
time="2018-08-17T18:34:28.652-06:00" level="fatal" msg="That could have gone better..." data.dead=true data.winner="walrus" caller="main.go:28 main.main.func1" host="myhost" trace.0="formatter.go:83 github.com/bdlm/log.getTrace" trace.1="formatter.go:153 github.com/bdlm/log.getData" trace.2="text_formatter.go:96 github.com/bdlm/log.(*TextFormatter).Format" trace.3="entry.go:171 github.com/bdlm/log.(*Entry).write" trace.4="entry.go:147 github.com/bdlm/log.Entry.log" trace.5="entry.go:234 github.com/bdlm/log.(*Entry).Fatal" trace.6="main.go:28 main.main.func1" trace.7="asm_amd64.s:573 runtime.call32" trace.8="panic.go:502 runtime.gopanic" trace.9="entry.go:155 github.com/bdlm/log.Entry.log" trace.10="entry.go:242 github.com/bdlm/log.(*Entry).Panic" trace.11="main.go:53 main.main" trace.12="proc.go:198 runtime.main" trace.13="asm_amd64.s:2361 runtime.goexit"
```

TTY trace output:
```go
log.SetFormatter(&log.TextFormatter{
    EnableTrace: true,
    ForceTTY: true,
})
```
<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-trace.png">
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
{"caller":"main.go:38 main.main","data":{"animal":"bird","count":1},"host":"myhost","level":"debug","msg":"Oh, look, a bird...","time":"2018-08-17T18:36:17.917-06:00","trace":["formatter.go:83 github.com/bdlm/log.getTrace","formatter.go:153 github.com/bdlm/log.getData","json_formatter.go:79 github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/log.Entry.log","entry.go:196 github.com/bdlm/log.(*Entry).Debug","main.go:38 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:42 main.main","data":{"animal":"walrus","count":20},"host":"myhost","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-17T18:36:17.917-06:00","trace":["formatter.go:83 github.com/bdlm/log.getTrace","formatter.go:153 github.com/bdlm/log.getData","json_formatter.go:79 github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/log.Entry.log","entry.go:203 github.com/bdlm/log.(*Entry).Info","main.go:42 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:46 main.main","data":{"animal":"walrus","count":100},"host":"myhost","level":"warn","msg":"The group's number increased tremendously!","time":"2018-08-17T18:36:17.917-06:00","trace":["formatter.go:83 github.com/bdlm/log.getTrace","formatter.go:153 github.com/bdlm/log.getData","json_formatter.go:79 github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/log.Entry.log","entry.go:215 github.com/bdlm/log.(*Entry).Warn","main.go:46 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:50 main.main","data":{"animal":"cow","count":"wait, what?"},"host":"myhost","level":"error","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-17T18:36:17.917-06:00","trace":["formatter.go:83 github.com/bdlm/log.getTrace","formatter.go:153 github.com/bdlm/log.getData","json_formatter.go:79 github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/log.Entry.log","entry.go:227 github.com/bdlm/log.(*Entry).Error","main.go:50 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:54 main.main","data":{"animal":"walrus","run":true},"host":"myhost","level":"panic","msg":"The walrus are attacking!","time":"2018-08-17T18:36:17.917-06:00","trace":["formatter.go:83 github.com/bdlm/log.getTrace","formatter.go:153 github.com/bdlm/log.getData","json_formatter.go:79 github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/log.Entry.log","entry.go:242 github.com/bdlm/log.(*Entry).Panic","main.go:54 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
{"caller":"main.go:29 main.main.func1","data":{"dead":true,"winner":"walrus"},"host":"myhost","level":"fatal","msg":"That could have gone better...","time":"2018-08-17T18:36:17.917-06:00","trace":["formatter.go:83 github.com/bdlm/log.getTrace","formatter.go:153 github.com/bdlm/log.getData","json_formatter.go:79 github.com/bdlm/log.(*JSONFormatter).Format","entry.go:171 github.com/bdlm/log.(*Entry).write","entry.go:147 github.com/bdlm/log.Entry.log","entry.go:234 github.com/bdlm/log.(*Entry).Fatal","main.go:29 main.main.func1","asm_amd64.s:573 runtime.call32","panic.go:502 runtime.gopanic","entry.go:155 github.com/bdlm/log.Entry.log","entry.go:242 github.com/bdlm/log.(*Entry).Panic","main.go:54 main.main","proc.go:198 runtime.main","asm_amd64.s:2361 runtime.goexit"]}
```

TTY trace output:
```go
log.SetFormatter(&log.JSONFormatter{
    EnableTrace: true,
    ForceTTY: true,
})
```
<p align="center">
    <img src="https://github.com/bdlm/log/wiki/assets/images/tty-json-trace.png">
</p>
