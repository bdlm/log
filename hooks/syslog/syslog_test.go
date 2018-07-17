package syslog

import (
	"log/syslog"
	"testing"

	"github.com/bdlm/log"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	logger := log.New()
	hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")

	if err != nil {
		t.Errorf("Unable to connect to local syslog.")
	}

	logger.Hooks.Add(hook)

	for _, level := range hook.Levels() {
		if len(logger.Hooks[level]) != 1 {
			t.Errorf("SyslogHook was not added. The length of logger.Hooks[%v]: %v", level, len(logger.Hooks[level]))
		}
	}

	logger.Info("Congratulations!")
}
