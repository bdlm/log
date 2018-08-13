# Syslog Hooks

## Usage

```go
import (
    "log/syslog"
    "github.com/bdlm/log"
    "github.com/bdlm/log/hooks/syslog"
)

func main() {
    log       := logrus.New()
    hook, err := syslog.NewSyslogHook(
        "udp",
        "localhost:514",
        syslog.LOG_INFO,
        "",
    )

    if err == nil {
        log.Hooks.Add(hook)
    }
}
```

If you want to connect to the local syslog (Ex. `/dev/log` or `/var/run/syslog` or `/var/run/log`), simply pass an empty string to the first two parameters of `NewSyslogHook`. For example:

```go
import (
    "log/syslog"
    "github.com/bdlm/log"
    syslog "github.com/bdlm/log/hooks/syslog"
)

func main() {
    log       := logrus.New()
    hook, err := syslog.NewSyslogHook(
        "",
        "",
        syslog.LOG_INFO,
        "",
    )

    if err == nil {
        log.Hooks.Add(hook)
    }
}
```
