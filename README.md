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

`bdlm/log` is a fork of the excellent [`sirupsen/logrus`](https://github.com/bdlm/log) package that adds support for sanitizing strings from logs to prevent accidental output of sensitive data.

# log

`bdlm/log` is a structured logger for Go and is API compatible with the standard libaray `log` package.

## Formats

By default, `bdlm/log` uses a basic text format:
```javascript
time="2018-08-10T18:19:09.424-06:00" level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count="1" caller="main.go:19 main.main"
time="2018-08-10T18:19:09.424-06:00" level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count="20" caller="main.go:23 main.main"
time="2018-08-10T18:19:09.424-06:00" level="warning" msg="The group's number increased tremendously!" data.animal="walrus" data.count="100" caller="main.go:27 main.main"
time="2018-08-10T18:19:09.424-06:00" level="warning" msg="A giant walrus monster appears!" data.animal="walrus" data.run="true" caller="main.go:31 main.main"
time="2018-08-10T18:19:09.424-06:00" level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.run="wait, what?" caller="main.go:35 main.main"
time="2018-08-10T18:19:09.424-06:00" level="fatal" msg="The walrus are attacking!" data.animal="walrus" data.panic="true" caller="main.go:39 main.main"
```

Color-coded output is used when a TTY terminal is detected for development:

<img src="https://github.com/bdlm/log/wiki/assets/images/tty.png" width="600px">

JSON formatting is also available with `log.SetFormatter(&log.JSONFormatter{})` for easy parsing by logstash or similar:

```json
{"caller":"main.go:17 main.main","data":{"animal":"bird","count":1},"host":"","level":"debug","msg":"Oh, look, a bird...","time":"2018-08-10T18:17:13.723-06:00"}
{"caller":"main.go:21 main.main","data":{"animal":"walrus","count":20},"host":"","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-10T18:17:13.723-06:00"}
{"caller":"main.go:25 main.main","data":{"animal":"walrus","count":100},"host":"","level":"warning","msg":"The group's number increased tremendously!","time":"2018-08-10T18:17:13.723-06:00"}
{"caller":"main.go:29 main.main","data":{"animal":"walrus","run":true},"host":"","level":"warning","msg":"A giant walrus monster appears!","time":"2018-08-10T18:17:13.723-06:00"}
{"caller":"main.go:33 main.main","data":{"animal":"cow","run":"wait, what?"},"host":"","level":"error","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-10T18:17:13.723-06:00"}
{"caller":"main.go:37 main.main","data":{"animal":"walrus","panic":true},"host":"","level":"fatal","msg":"The walrus are attacking!","time":"2018-08-10T18:17:13.723-06:00"}
```

## Examples

The simplest way to use `bdlm/log` is simply the package-level exported logger:

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

Note that it is completely api-compatible with the stdlib logger, so you can replace your `log`
imports everywhere with `"github.com/bdlm/log"` and you'll have the full flexibility of
`bdlm/log` available. You can customize it further in your code:

```go
package main

import (
  "os"
  log "github.com/bdlm/log"
)

func init() {
  // Log as JSON instead of the default ASCII formatter.
  log.SetFormatter(&log.JSONFormatter{})

  // Output to stdout instead of the default stderr. Can be any io.Writer, see
  // below for a File example.
  log.SetOutput(os.Stdout)

  // Only log the warning severity or above.
  log.SetLevel(log.WarnLevel)
}

func main() {
  log.WithFields(log.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")

  log.WithFields(log.Fields{
    "omg":    true,
    "number": 122,
  }).Warn("The group's number increased tremendously!")

  log.WithFields(log.Fields{
    "omg":    true,
    "number": 100,
  }).Fatal("The ice breaks!")

  // A common pattern is to re-use fields between logging statements by re-using
  // the log.Entry returned from WithFields()
  contextLogger := log.WithFields(log.Fields{
    "common": "this is a common field",
    "other": "I also should be logged always",
  })

  contextLogger.Info("I'll be logged with common and other field")
  contextLogger.Info("Me too")
}
```

`bdlm/log` supports a "blacklist" of values that should not be logged. This can be used to help prevent or mitigate leaking secrets into log files:

```go
import (
    "github.com/bdlm/log"
)

func main() {
    log.AddSecret("some-secret-text")
    log.Info("the secret is 'some-secret-text'")

    // Output: the secret is '****************'
}
```

For more advanced usage such as logging to multiple locations from the same application, you can also create an instance of the `bdlm/log` Logger:

```go
package main

import (
  "os"
  "github.com/bdlm/log"
)

// Create a new instance of the logger. You can have any number of instances.
var logger = log.New()

func main() {
  // The API for setting attributes is a little different than the package level
  // exported logger. See Godoc.
  logger.Out = os.Stdout

  // You could set this to any `io.Writer` such as a file
  // file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY, 0666)
  //  if err == nil {
  //    logger.Out = file
  //  } else {
  //    logger.Info("Failed to log to file, using default stderr")
  // }

  logger.WithFields(log.Fields{
    "animal": "walrus",
    "size":   10,
  }).Info("A group of walrus emerges from the ocean")
}
```

## Fields

`bdlm/log` encourages careful, structured logging through logging fields instead of long, unparseable error messages. For example, instead of: `log.Fatalf("Failed to send event %s to topic %s with key %d")`, you should log the much more discoverable:

```go
log.WithFields(log.Fields{
  "event": event,
  "topic": topic,
  "key": key,
}).Fatal("Failed to send event")
```

We've found this API forces you to think about logging in a way that produces much more useful logging messages. We've been in countless situations where just a single added field to a log statement that was already there would've saved us hours. The `WithFields` call is optional.

In general, with `bdlm/log` using any of the `printf`-family functions should be seen as a hint you should add a field, however, you can still use the `printf`-family functions with `bdlm/log`.

### Default Fields

Often it's helpful to have fields _always_ attached to log statements in an application or parts of one. For example, you may want to always log the `request_id` and `user_ip` in the context of a request. Instead of writing `log.WithFields(log.Fields{"request_id": request_id, "user_ip": user_ip})` on every line, you can create a `log.Entry` to pass around instead:

```go
requestLogger := log.WithFields(log.Fields{"request_id": request_id, "user_ip": user_ip})
requestLogger.Info("something happened on that request") # will log request_id and user_ip
requestLogger.Warn("something not great happened")
```

## Hooks

You can add hooks for logging levels. For example to send errors to an exception tracking service on `Error`, `Fatal` and `Panic`, info to StatsD or log to multiple places simultaneously, e.g. syslog.

`bdlm/log` comes with [built-in hooks](hooks/). Add those, or your custom hook, in `init`:

```go
import (
  log "github.com/bdlm/log"
  log_syslog "github.com/bdlm/log/hooks/syslog"
  "log/syslog"
)

func init() {

  // Use the Airbrake hook to report errors that have Error severity or above to
  // an exception tracker. You can create custom hooks, see the Hooks section.
  log.AddHook(airbrake.NewHook(123, "xyz", "production"))

  hook, err := log_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
  if err != nil {
    log.Error("Unable to connect to local syslog daemon")
  } else {
    log.AddHook(hook)
  }
}
```
Note: Syslog hook also support connecting to local syslog (Ex. "/dev/log" or "/var/run/syslog" or "/var/run/log"). For the detail, please check the [syslog hook README](hooks/syslog/README.md).

A list of currently known of service hook can be found in this wiki [page](https://github.com/bdlm/log/wiki/Hooks)


## Level logging

`bdlm/log` has six logging levels: Debug, Info, Warning, Error, Fatal and Panic.

```go
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
```

You can set the logging level on a `Logger`, then it will only log entries with that severity or anything above it:

```go
// Will log anything that is info or above (warn, error, fatal, panic). Default.
log.SetLevel(log.InfoLevel)
```

It may be useful to set `log.Level = log.DebugLevel` in a debug or verbose environment if your application has that.

## Entries

Besides the fields added with `WithField` or `WithFields` some fields are automatically added to all logging events:

1. `time`. The timestamp when the entry was created.
1. `msg`. The logging message passed to `{Info,Warn,Error,Fatal,Panic}` after the `AddFields` call. E.g. `Failed to send event.`
1. `level`. The logging level. E.g. `info`.

## Environments

`bdlm/log` has no notion of environment.

If you wish for hooks and formatters to only be used in specific environments, you should handle that yourself. For example, if your application has a global variable `Environment`, which is a string representation of the environment you could do:

```go
import (
  log "github.com/bdlm/log"
)

init() {
  // do something here to set environment depending on an environment variable
  // or command-line flag
  if Environment == "production" {
    log.SetFormatter(&log.JSONFormatter{})
  } else {
    // The TextFormatter is default, you don't actually have to do this.
    log.SetFormatter(&log.TextFormatter{})
  }
}
```

This configuration is how `bdlm/log` was intended to be used, but JSON in production is mostly only useful if you do log aggregation with tools like Splunk or Logstash.

## Formatters

The built-in logging formatters are:

* `log.TextFormatter`. Logs the event in colors if stdout is a tty, otherwise without colors.
  * *Note:* to force colored output when there is no TTY, set the `ForceColors` field to `true`.  To force no colored output even if there is a TTY  set the `DisableColors` field to `true`. For Windows, see [github.com/mattn/go-colorable](https://github.com/mattn/go-colorable).
  * When colors are enabled, levels are truncated to 4 characters by default. To disable truncation set the `DisableLevelTruncation` field to `true`.
  * All options are listed in the [generated docs](https://godoc.org/github.com/bdlm/log#TextFormatter).
* `log.JSONFormatter`. Logs fields as JSON.
  * All options are listed in the [generated docs](https://godoc.org/github.com/bdlm/log#JSONFormatter).

Third party logging formatters:

* [`FluentdFormatter`](https://github.com/joonix/log). Formats entries that can be parsed by Kubernetes and Google * [`zalgo`](https://github.com/aybabtme/logzalgo). Invoking the P͉̫o̳̼̊w̖͈̰͎e̬͔̭͂r͚̼̹̲ ̫͓͉̳͈ō̠͕͖̚f̝͍̠ ͕̲̞͖͑Z̖̫̤̫ͪa͉̬͈̗l͖͎g̳̥o̰̥̅!̣͔̲̻͊̄ ̙̘̦̹̦.

You can define your formatter by implementing the `Formatter` interface, requiring a `Format` method. `Format` takes an `*Entry`. `entry.Data` is a `Fields` type (`map[string]interface{}`) with all your fields as well as the default ones (see Entries section above):

```go
type MyJSONFormatter struct {
}

log.SetFormatter(new(MyJSONFormatter))

func (f *MyJSONFormatter) Format(entry *Entry) ([]byte, error) {
  // Note this doesn't include Time, Level and Message which are available on
  // the Entry. Consult `godoc` on information about those fields or read the
  // source of the official loggers.
  serialized, err := json.Marshal(entry.Data)
    if err != nil {
      return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
    }
  return append(serialized, '\n'), nil
}
```

## Logger as an `io.Writer`

`bdlm/log` can be transformed into an `io.Writer`. That writer is the end of an `io.Pipe` and it is your responsibility to close it.

```go
w := logger.Writer()
defer w.Close()

srv := http.Server{
    // create a stdlib log.Logger that writes to
    // log.Logger.
    ErrorLog: log.New(w, "", 0),
}
```

Each line written to that writer will be printed the usual way, using formatters and hooks. The level for those entries is `info`.

This means that we can override the standard library logger easily:

```go
logger := log.New()
logger.Formatter = &log.JSONFormatter{}

// Use `bdlm/log` for standard log output
// Note that `log` here references stdlib's log
log.SetOutput(logger.Writer())
```

## Rotation

Log rotation is not provided with `bdlm/log`. Log rotation should be done by an external program (like `logrotate(8)`) that can compress and delete old log entries. It should not be a feature of the application-level logger.

## Testing

`bdlm/log` has a built in facility for asserting the presence of log messages. This is implemented through the `test` hook and provides:

* decorators for existing logger (`test.NewLocal` and `test.NewGlobal`) which basically just add the `test` hook
* a test logger (`test.NewNullLogger`) that just records log messages (and does not output any):

```go
import(
  "github.com/bdlm/log"
  "github.com/bdlm/log/hooks/test"
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestSomething(t*testing.T){
  logger, hook := test.NewNullLogger()
  logger.Error("Helloerror")

  assert.Equal(t, 1, len(hook.Entries))
  assert.Equal(t, log.ErrorLevel, hook.LastEntry().Level)
  assert.Equal(t, "Helloerror", hook.LastEntry().Message)

  hook.Reset()
  assert.Nil(t, hook.LastEntry())
}
```

## Fatal handlers

`bdlm/log` can register one or more functions that will be called when any `fatal` level message is logged. The registered handlers will be executed before `bdlm/log` performs a `os.Exit(1)`. This behavior may be helpful if callers need to gracefully shutdown. Unlike a `panic("Something went wrong...")` call which can be intercepted with a deferred `recover` a call to `os.Exit(1)` can not be intercepted.

```go
...
handler := func() {
  // gracefully shutdown something...
}
log.RegisterExitHandler(handler)
...
```

## Thread safety

By default, Logger is protected by a mutex for concurrent writes. The mutex is held when calling hooks and writing logs. If you are sure such locking is not needed, you can call logger.SetNoLock() to disable the locking.

Situation when locking is not needed includes:

* You have no hooks registered, or hooks calling is already thread-safe.

* Writing to logger.Out is already thread-safe, for example:
  1. logger.Out is protected by locks.
  1. logger.Out is a os.File handler opened with `O_APPEND` flag, and every write is smaller than 4k. (This allow multi-thread/multi-process writing)

  (Refer to http://www.notthewizard.com/2014/06/17/are-files-appends-really-atomic/)
