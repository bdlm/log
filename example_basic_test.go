package log_test

import (
	"os"

	"github.com/bdlm/log"
)

func Example_basic() {
	var logger = log.New()
	logger.Formatter = new(log.TextFormatter)                     //default
	logger.Formatter.(*log.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Formatter.(*log.TextFormatter).DisableHostname = true  // remove timestamp from test output
	logger.Formatter.(*log.TextFormatter).DisableCaller = true    // remove caller from test output
	logger.Level = log.DebugLevel
	logger.Out = os.Stdout

	// Capture the panic result
	defer func() {
		err := recover()
		if err != nil {
			entry := err.(*log.Entry)
			logger.WithFields(log.Fields{
				"winner": entry.Data["animal"],
				"dead":   true,
			}).Error("That could have gone better...")
		}
	}()

	logger.WithFields(log.Fields{
		"animal": "bird",
		"count":  1,
	}).Debug("Oh, look, a bird...")
	logger.WithFields(log.Fields{
		"animal": "walrus",
		"count":  20,
	}).Info("A group of walrus emerges from the ocean")
	logger.WithFields(log.Fields{
		"animal": "walrus",
		"count":  100,
	}).Warn("The group's number increased tremendously!")
	logger.WithFields(log.Fields{
		"animal": "cow",
		"run":    "wait, what?",
	}).Error("Tremendously sized cow enters the ocean.")
	logger.WithFields(log.Fields{
		"animal": "walrus",
		"run":    true,
	}).Panic("The walrus are attacking!")

	// Output:
	// level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count="1"
	// level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count="20"
	// level="warning" msg="The group's number increased tremendously!" data.animal="walrus" data.count="100"
	// level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.run="wait, what?"
	// level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run="true"
	// level="error" msg="That could have gone better..." data.dead="true" data.winner="walrus"
}
