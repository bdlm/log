# v0.1.13

* TTY color scheme update

# v0.1.12: Revert "move Fields type to bdlm/std (#13)" (#14)

* This reverts commit da2feacffefce803820e8c090306bffb59d3f08c.

# v0.1.11 move Fields type to bdlm/std (#13)


# v0.1.10

* adds a stdlib compatible formatter

# v0.1.9

* don't remove empty fields
* don't escape log messages in text TTY output

# v0.1.8

* remove message truncation in text TTY output

# v0.1.7

* Documentation updates and minor cleanup.

# v0.1.6 Implement the std.Logger interface

* implements the github.com/bdlm/std:Logger interface
* adds support for a verbose trace logging mode.

# v0.1.5 implement various PRs listed on sirupsen/logrus

* sirupsen/logrus/pull/664
* sirupsen/logrus/pull/647
* sirupsen/logrus/pull/687
* sirupsen/logrus/pull/685
* sirupsen/logrus/pull/788 (existed previously)

This also updates the string escape logic, all values are now JSON escaped fixes an issue with internal properties being included in JSON format adds new fields to unit tests (data and caller) minor cleanup of text templates

# v0.1.4

* updated tty formatting
* added tty formatting to JSON output

# v0.1.3

* TTY format updates
* Minor cleanup

# v0.1.2: cleanup goreportcard errors (#5)

* cleanup 'ineffassign' errors
* reduce cyclomatic complexity
* update documentation

# v0.1.1: update build (#3)

* Update TTY formatting

# v0.1.0 cleanup README

