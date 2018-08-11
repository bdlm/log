/*
Package log is a fork of the excellent [`sirupsen/logrus`](https://github.com/bdlm/log)
package.

log is a structured logger for Go, completely API compatible with the standard library logger.



Package-level exported logger

  package main

  import (
    "github.com/bdlm/log"
  )

  func main() {
    log.WithFields(log.Fields{
      "animal": "walrus",
      "number": 1,
      "size":   10,
    }).Info("A walrus appears")
  }

Output:
  time="2015-09-07T08:48:33Z" level=info msg="A walrus appears" animal=walrus number=1 size=10



API Compatibility

Note that it is completely api-compatible with the stdlib logger, so you can replace your `log`
imports everywhere with `"github.com/bdlm/log"` and you'll have the full flexibility of
`bdlm/log` available. You can customize it further in your code:

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
    // the logrus.Entry returned from WithFields()
    contextLogger := log.WithFields(log.Fields{
      "common": "this is a common field",
      "other": "I also should be logged always",
    })

    contextLogger.Info("I'll be logged with common and other field")
    contextLogger.Info("Me too")
  }



Managing secrets

`bdlm/log` also supports a "blacklist" of values that should not be logged. This can be used to
help prevent or mitigate leaking secrets into log files:

  import (
      "github.com/bdlm/log"
  )

  func main() {
      log.AddSecret("some-secret-text")
      log.Info("the secret is 'some-secret-text'")

      // Output: the secret is '****************'
  }



Output to multiple locations

For more advanced usage such as logging to multiple locations from the same application, you can
create an instance of the `bdlm/log` Logger:

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
    // file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
    //  if err == nil {
    //    logger.Out = file
    //  } else {
    //    logger.Info("Failed to log to file, using default stderr")
    // }

    logger.WithFields(logrus.Fields{
      "animal": "walrus",
      "size":   10,
    }).Info("A group of walrus emerges from the ocean")
  }


Output features

Color-coded output is available when attached to a TTY for development. A JSON formatter is also
available for easy parsing by logstash or Splunk:

  log.SetFormatter(&log.JSONFormatter{})

Output:
  {"caller":"main.go:37 main.main","data":{"animal":"walrus"},"host":"myhost","level":"info","msg":"A group of walrus emerges from the ocean","time":"2018-08-10T23:08:02.860Z"}
  {"caller":"main.go:61 main.main","host":"myhost","level":"warning","msg":"The group's number increased tremendously!","number":122,"omg":true,"time":"2018-08-10T23:08:02.863Z"}
  {"caller":"main.go:99 main.main","data":{"animal":"walrus"},"host":"myhost","level":"info","msg":"A giant walrus appears!","time":"2018-08-10T23:08:02.877Z"}
  {"caller":"main.go:61 main.main","data":{"animal":"walrus","host":"myhost","level":"info","msg":"Tremendously sized cow enters the ocean.","time":"2018-08-10T23:08:02.877Z"}
  {"caller":"main.go:99 main.main","host":"myhost","level":"fatal","msg":"The ice breaks!","number":100,"omg":true,"time":"2018-08-10T23:08:03.566Z"}

For a full guide visit https://github.com/bdlm/log
*/
package log
