package log_test

import (
	"fmt"
	"os"

	errs "github.com/bdlm/errors"
	"github.com/bdlm/log/v2"
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
	// level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1
	// level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20
	// level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100
	// level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.run="wait, what?"
	// level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true
	// level="error" msg="That could have gone better..." data.dead=true data.winner="walrus"
}

func Example_JSON() {
	var logger = log.New()
	logger.Formatter = new(log.JSONFormatter)                     //default
	logger.Formatter.(*log.JSONFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableHostname = true  // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableCaller = true    // remove caller from test output
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
	// {"data":{"animal":"bird","count":1},"error":null,"level":"debug","msg":"Oh, look, a bird..."}
	// {"data":{"animal":"walrus","count":20},"error":null,"level":"info","msg":"A group of walrus emerges from the ocean"}
	// {"data":{"animal":"walrus","count":100},"error":null,"level":"warn","msg":"The group's number increased tremendously!"}
	// {"data":{"animal":"cow","run":"wait, what?"},"error":null,"level":"error","msg":"Tremendously sized cow enters the ocean."}
	// {"data":{"animal":"walrus","run":true},"error":null,"level":"panic","msg":"The walrus are attacking!"}
	// {"data":{"dead":true,"winner":"walrus"},"error":null,"level":"error","msg":"That could have gone better..."}
}

func Example_JSONTTY() {
	var logger = log.New()
	logger.Formatter = new(log.JSONFormatter)                     //default
	logger.Formatter.(*log.JSONFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableHostname = true  // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableCaller = true    // remove caller from test output
	logger.Formatter.(*log.JSONFormatter).ForceTTY = true         // remove caller from test output
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

	e := fmt.Errorf("error 1")
	e = errs.Wrap(e, "error 2")
	e = errs.Wrap(e, "error 3")
	log.WithError(e).WithField("some", "field").Info("GOT HERE")

	// Output:
	// {
	//     "[38;5;245mlevel[0m": "[38;5;245mdebug[0m",
	//     "[38;5;245mmsg[0m": "Oh, look, a bird...",
	//     "[38;5;245mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"bird"[0m,
	//         "[38;5;111mcount[0m": [38;5;180m1[0m
	//     },
	// }
	// {
	//     "[38;5;46mlevel[0m": "[38;5;46m info[0m",
	//     "[38;5;46mmsg[0m": "A group of walrus emerges from the ocean",
	//     "[38;5;46mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"walrus"[0m,
	//         "[38;5;111mcount[0m": [38;5;180m20[0m
	//     },
	// }
	// {
	//     "[38;5;226mlevel[0m": "[38;5;226m warn[0m",
	//     "[38;5;226mmsg[0m": "The group's number increased tremendously!",
	//     "[38;5;226mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"walrus"[0m,
	//         "[38;5;111mcount[0m": [38;5;180m100[0m
	//     },
	// }
	// {
	//     "[38;5;208mlevel[0m": "[38;5;208merror[0m",
	//     "[38;5;208mmsg[0m": "Tremendously sized cow enters the ocean.",
	//     "[38;5;208mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"cow"[0m,
	//         "[38;5;111mrun[0m": [38;5;180m"wait, what?"[0m
	//     },
	// }
	// {
	//     "[38;5;196mlevel[0m": "[38;5;196mpanic[0m",
	//     "[38;5;196mmsg[0m": "The walrus are attacking!",
	//     "[38;5;196mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"walrus"[0m,
	//         "[38;5;111mrun[0m": [38;5;180mtrue[0m
	//     },
	// }
	// {
	//     "[38;5;208mlevel[0m": "[38;5;208merror[0m",
	//     "[38;5;208mmsg[0m": "That could have gone better...",
	//     "[38;5;208mdata[0m": {
	//         "[38;5;111mdead[0m": [38;5;180mtrue[0m,
	//         "[38;5;111mwinner[0m": [38;5;180m"walrus"[0m
	//     },
	// }
}

func Example_basic_withError() {
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
			e := fmt.Errorf("First mistake: not running when a walrus herd \"emerged\" from the ocean")
			e = errs.Wrap(e, "Second mistake: a walrus cow is not cattle...")

			entry := err.(*log.Entry)
			logger.WithError(e).WithFields(log.Fields{
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
	// level="debug" msg="Oh, look, a bird..." data.animal="bird" data.count=1
	// level="info" msg="A group of walrus emerges from the ocean" data.animal="walrus" data.count=20
	// level="warn" msg="The group's number increased tremendously!" data.animal="walrus" data.count=100
	// level="error" msg="Tremendously sized cow enters the ocean." data.animal="cow" data.run="wait, what?"
	// level="panic" msg="The walrus are attacking!" data.animal="walrus" data.run=true
	// level="error" msg="That could have gone better..." error=Second mistake: a walrus cow is not cattle... data.dead=true data.winner="walrus"
}

func Example_JSON_withError() {
	var logger = log.New()
	logger.Formatter = new(log.JSONFormatter)                     //default
	logger.Formatter.(*log.JSONFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableHostname = true  // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableCaller = true    // remove caller from test output
	logger.Level = log.DebugLevel
	logger.Out = os.Stdout

	// Capture the panic result
	defer func() {
		err := recover()
		if err != nil {
			e := fmt.Errorf("First mistake: not running when a walrus herd \"emerged\" from the ocean")
			e = errs.Wrap(e, "Second mistake: a walrus cow is not cattle...")

			entry := err.(*log.Entry)
			logger.WithError(e).WithFields(log.Fields{
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
	// {"data":{"animal":"bird","count":1},"error":null,"level":"debug","msg":"Oh, look, a bird..."}
	// {"data":{"animal":"walrus","count":20},"error":null,"level":"info","msg":"A group of walrus emerges from the ocean"}
	// {"data":{"animal":"walrus","count":100},"error":null,"level":"warn","msg":"The group's number increased tremendously!"}
	// {"data":{"animal":"cow","run":"wait, what?"},"error":null,"level":"error","msg":"Tremendously sized cow enters the ocean."}
	// {"data":{"animal":"walrus","run":true},"error":null,"level":"panic","msg":"The walrus are attacking!"}
	// {"data":{"dead":true,"winner":"walrus"},"error":[{"caller":"#0 example_basic_test.go:280 (github.com/bdlm/log/v2_test.Example_JSON_withError.func1)","error":"Second mistake: a walrus cow is not cattle..."}],"level":"error","msg":"That could have gone better..."}
}

func Example_JSONTTY_withError() {
	var logger = log.New()
	logger.Formatter = new(log.JSONFormatter)                     //default
	logger.Formatter.(*log.JSONFormatter).DisableTimestamp = true // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableHostname = true  // remove timestamp from test output
	logger.Formatter.(*log.JSONFormatter).DisableCaller = true    // remove caller from test output
	logger.Formatter.(*log.JSONFormatter).ForceTTY = true         // remove caller from test output
	logger.Level = log.DebugLevel
	logger.Out = os.Stdout

	// Capture the panic result
	defer func() {
		err := recover()
		if err != nil {
			e := fmt.Errorf("First mistake: not running when a walrus herd \"emerged\" from the ocean")
			e = errs.Wrap(e, "Second mistake: a walrus cow is not cattle...")

			entry := err.(*log.Entry)
			logger.WithError(e).WithFields(log.Fields{
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

	e := fmt.Errorf("error 1")
	e = errs.Wrap(e, "error 2")
	e = errs.Wrap(e, "error 3")
	log.WithError(e).WithField("some", "field").Info("GOT HERE")

	// Output:
	// {
	//     "[38;5;245mlevel[0m": "[38;5;245mdebug[0m",
	//     "[38;5;245mmsg[0m": "Oh, look, a bird...",
	//     "[38;5;245mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"bird"[0m,
	//         "[38;5;111mcount[0m": [38;5;180m1[0m
	//     },
	// }
	// {
	//     "[38;5;46mlevel[0m": "[38;5;46m info[0m",
	//     "[38;5;46mmsg[0m": "A group of walrus emerges from the ocean",
	//     "[38;5;46mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"walrus"[0m,
	//         "[38;5;111mcount[0m": [38;5;180m20[0m
	//     },
	// }
	// {
	//     "[38;5;226mlevel[0m": "[38;5;226m warn[0m",
	//     "[38;5;226mmsg[0m": "The group's number increased tremendously!",
	//     "[38;5;226mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"walrus"[0m,
	//         "[38;5;111mcount[0m": [38;5;180m100[0m
	//     },
	// }
	// {
	//     "[38;5;208mlevel[0m": "[38;5;208merror[0m",
	//     "[38;5;208mmsg[0m": "Tremendously sized cow enters the ocean.",
	//     "[38;5;208mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"cow"[0m,
	//         "[38;5;111mrun[0m": [38;5;180m"wait, what?"[0m
	//     },
	// }
	// {
	//     "[38;5;196mlevel[0m": "[38;5;196mpanic[0m",
	//     "[38;5;196mmsg[0m": "The walrus are attacking!",
	//     "[38;5;196mdata[0m": {
	//         "[38;5;111manimal[0m": [38;5;180m"walrus"[0m,
	//         "[38;5;111mrun[0m": [38;5;180mtrue[0m
	//     },
	// }
	// {
	//     "[38;5;208mlevel[0m": "[38;5;208merror[0m",
	//     "[38;5;208mmsg[0m": "That could have gone better...",
	//     "[38;5;208merror[0m": "[38;5;208m[
	//     {
	//         "caller": "#0 example_basic_test.go:335 (github.com/bdlm/log/v2_test.Example_JSONTTY_withError.func1)",
	//         "error": "Second mistake: a walrus cow is not cattle..."
	//     }
	// ][0m",
	//     "[38;5;208mdata[0m": {
	//         "[38;5;111mdead[0m": [38;5;180mtrue[0m,
	//         "[38;5;111mwinner[0m": [38;5;180m"walrus"[0m
	//     },
	// }
}
