All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

- **Major**: backwards incompatible package updates
- **Minor**: feature additions
- **Patch**: bug fixes, API route/ingress updates, DNS updates, backward compatible protobuf model chaanges

# v2.0.0 - 2020-05-01
#### Added
* `go.mod`
* `github.com/bdlm/std` error interface

#### Changed
*

#### Removed
* gRPC request interceptor.

# v0.1.20

#### Added
* gRPC request interceptor.

# v0.1.19
#### Added
* Expose the `LevelString` function for use in custom formatters.

# v0.1.18
#### Changed
* properly encode json values in TTY output

# v0.1.17
#### Changed
* fix type display logic
* update pr template
* update readme

# v0.1.16
#### Changed
* update TTY color scheme and layout
* cleanup template logic
* template fixes for v1.9 and v1.10

# v0.1.15
#### Added
* caller level adjustment control with SetCallerLevel(level uint)

# v0.1.13
#### Changed
* TTY color scheme

# v0.1.12: Revert "move Fields type to bdlm/std (#13)" (#14)
#### Changed
* This reverts commit da2feacffefce803820e8c090306bffb59d3f08c.

# v0.1.11
#### Changed
* move Fields type to bdlm/std (#13)


# v0.1.10
#### Added
* adds a stdlib compatible formatter

# v0.1.9
#### Changed
* don't remove empty fields
* don't escape log messages in text TTY output

# v0.1.8
#### Changed
* remove message truncation in text TTY output

# v0.1.7
#### Changed
* Documentation updates and minor cleanup.

# v0.1.6 Implement the std.Logger interface
#### Added
* implements the github.com/bdlm/std:Logger interface
* adds support for a verbose trace logging mode.

# v0.1.5
#### Added
implement various PRs listed on sirupsen/logrus

* sirupsen/logrus/pull/664
* sirupsen/logrus/pull/647
* sirupsen/logrus/pull/687
* sirupsen/logrus/pull/685
* sirupsen/logrus/pull/788 (existed previously)

This also updates the string escape logic, all values are now JSON escaped fixes an issue with internal properties being included in JSON format adds new fields to unit tests (data and caller) minor cleanup of text templates

# v0.1.4
#### Added
* added tty formatting to JSON output
#### Changed
* updated tty formatting

# v0.1.3
#### Changed
* TTY format updates
* Minor cleanup

# v0.1.2
#### Changed
cleanup goreportcard errors (#5)
* cleanup 'ineffassign' errors
* reduce cyclomatic complexity
* update documentation

# v0.1.1: update build (#3)
#### Changed
* Update TTY formatting

# v0.1.0
#### Changed
* cleanup README
