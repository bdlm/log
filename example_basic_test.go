package log_test

import (
	"os"

	"github.com/bdlm/log"
)

func Example_basic() {
	var logger = log.New()
	logger.Formatter = new(log.JSONFormatter)
	logger.Formatter = new(log.TextFormatter)                     //default
	logger.Formatter.(*log.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Level = log.DebugLevel
	logger.Out = os.Stdout

	// file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	defer func() {
		err := recover()
		if err != nil {
			entry := err.(*log.Entry)
			logger.WithFields(log.Fields{
				"omg":         true,
				"err_animal":  entry.Data["animal"],
				"err_size":    entry.Data["size"],
				"err_level":   entry.Level,
				"err_message": entry.Message,
				"number":      100,
			}).Error("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
		}
	}()

	logger.WithFields(log.Fields{
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	logger.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logger.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	logger.WithFields(log.Fields{
		"temperature": -4,
	}).Debug("Temperature changes")

	logger.WithFields(log.Fields{
		"animal": "orca",
		"size":   9009,
	}).Panic("It's over 9000!")

	// Output:
	// level=debug msg="Started observing beach" animal=walrus number=8
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=warning msg="The group's number increased tremendously!" number=122 omg=true
	// level=debug msg="Temperature changes" temperature=-4
	// level=panic msg="It's over 9000!" animal=orca size=9009
	// level=error msg="The ice breaks!" err_animal=orca err_level=panic err_message="It's over 9000!" err_size=9009 number=100 omg=true
}
