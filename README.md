[![GitHub version](https://img.shields.io/github/v/release/jeffalyanak/check_https_go)](https://github.com/jeffalyanak/check_https_go/releases/latest)
[![License](https://img.shields.io/github/license/jeffalyanak/check_https_go)](https://github.com/jeffalyanak/check_https_go/blob/master/LICENSE)
[![Donate](https://img.shields.io/badge/donate--green)](https://jeff.alyanak.ca/donate)

# Golang Icinga/Nagios HTTPS Checker

Icinga/Nagios plugin, checks that a site returns a 2xx-series `status code`, `returns content`, and has a `valid certificate`.

User configurable `warning` and `critical` levels for the number of days left in the certificate validity period.

## Installation and requirements

The pre-compiled binaries available on the [releases page](https://github.com/jeffalyanak/check_https_go/releases) are self-contained and have no dependancies to run.

If you wish to compile it yourself, you'll need to install `go` and `make`. It's been tested on:

* Golang 1.14.1
* GNU Make 4.2.1

It'll probably build just fine on many other versions. To build, simply run `make`:

```bash
make           # Make all builds
make linux     # Make linux binary
make windows   # Make windows binary
make macos     # Make macOS binary
```

You can run `make help` for additional options.

## Usage

```bash
usage:
  required
    -h string
        Fully-qualified domain name to check.
  optional
    -w int
        Number of days for which the TLS certificate must be valid before a warning state is returned. (default 10)
    -c int
        Number of days for which the TLS certificate must be valid before a critical state is returned. (default 5)
    -t
        Timeout length in seconds, requests that do not finish before timeout are considered failed. (default 30)
```

## License

Golang Icinga/Nagios HTTPS Checker is licensed under the terms of the GNU General Public License Version 3.