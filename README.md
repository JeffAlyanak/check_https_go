[![GitHub version](https://img.shields.io/github/v/release/jeffalyanak/check_https_go)](https://github.com/jeffalyanak/check_https_go/releases/latest)
[![License](https://img.shields.io/github/license/jeffalyanak/check_https_go)](https://github.com/jeffalyanak/check_https_go/blob/master/LICENSE)
[![Donate](https://img.shields.io/badge/donate--green)](https://jeff.alyanak.ca/donate)

# Golang Icinga/Nagios HTTPS Checker

Icinga/Nagios plugin, checks that a site returns an `expected status code`, `returns expected content`, and has a `valid certificate`.

User configurable `warning` and `critical` levels for the number of days left in the certificate validity period.

## Installation and requirements

The pre-compiled binaries available on the [releases page](https://github.com/jeffalyanak/check_https_go/releases) are self-contained and have no dependancies to run.

If you wish to compile it yourself, you'll need to install `go` and `make`. It's been tested on:

* Golang 1.20.3
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
    -a string
        Comma-seperated list of status codes. (default "200,201,202,203,204,205,206,207,208,226")
    -c int
            Number of days for which the TLS certificate must be valid before a critical state is returned. (default 5)
    -r int
            Number of redirects to follow. (default 20)
    -s string
            Custom string to check for in the response body. (default "<!DOCTYPE HTML>")
    -t int
            Timeout length in seconds, requests that do not finish before timeout are considered failed. (default 30)
    -u string
            Custom user-agent string. (default "check_https_go")
    -v    More verbose output includes details of any redirects.
    -w int
            Number of days for which the TLS certificate must be valid before a warning state is returned. (default 10)
```

## Version history

1.4—Add parameter for configuring status code check.
1.3—Add custom string to check the content returned by the page.
1.2—Follow Redirect feature now handles 307s.
1.1—Added configurable timeout. Default remains at 30, but you can now increase or decrease this.

## License

Golang Icinga/Nagios HTTPS Checker is licensed under the terms of the GNU General Public License Version 3.
