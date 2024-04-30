package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jeffalyanak/check_https_go/check"
)

func main() {
	// Handle cli arguments
	host := flag.String("h", "", "Fully-qualified domain name to check.")
	checkString := flag.String("s", "<!DOCTYPE HTML>", "Custom string to check for in the response body.")
	userAgent := flag.String("u", "check_https_go", "Custom user-agent string.")
	verbose := flag.Bool("v", false, "More verbose output includes details of any redirects.")
	redirects := flag.Int("r", 20, "Number of redirects to follow.")
	certwarn := flag.Int("w", 10, "Number of days for which the TLS certificate must be valid before a warning state is returned.")
	certcrit := flag.Int("c", 5, "Number of days for which the TLS certificate must be valid before a critical state is returned.")
	timeoutduration := flag.Int("t", 30, "Timeout length in seconds, requests that do not finish before timeout are considered failed.")
	statusCodes := flag.String("a", "200,201,202,203,204,205,206,207,208,226", "Comma-seperated list of status codes.")

	flag.Parse()

	if *host == "" {
		fmt.Println("Please provide a fully-qualified domain name.")
		os.Exit(3)
	}

	regex := regexp.MustCompile(`^\d+(,\d+)*$`)
	if !regex.MatchString(*statusCodes) {
		fmt.Println("Status Codes must be provided as a comma-seperated string. Eg: 200,201,202")
		os.Exit(3)
	}

	// Primary var for the checks
	var h check.HTTPCheck
	h.URL = *host

	// Create mew Performance Data struct and start the timer
	var perfData check.PerfData
	perfData.StartTimer(time.Now())

	// Perform status code check, exit with additional info if error
	statuscodeResult := h.CheckStatus(*redirects, *userAgent, *timeoutduration, *statusCodes)
	if statuscodeResult.Error != nil {
		printIntro("Status Code Error", h.URL)
		fmt.Println(statuscodeResult.Error)
		if *verbose {
			printVerboseInfo(statuscodeResult.VerboseValue)
		}
		fmt.Println(perfData.Get())
		os.Exit(3)
	}

	// Check return code, exit with additional info if non-zero
	if statuscodeResult.ReturnCode != 0 {
		printIntro("Status Code Error", h.URL)
		printStatusCode(statuscodeResult.Status, statuscodeResult.Value, *statusCodes)
		if *verbose {
			printVerboseInfo(statuscodeResult.VerboseValue)
		}
		fmt.Println(perfData.Get())
		os.Exit(statuscodeResult.ReturnCode)
	}

	// Perform content check, exit with additional info if error
	contentResult := h.CheckContent(*checkString)
	if contentResult.Error != nil {
		printIntro("Web Content Error", h.URL)
		fmt.Println(statuscodeResult.Error)
		if *verbose {
			printVerboseInfo(statuscodeResult.VerboseValue +
				contentResult.VerboseValue)
		}
		fmt.Println(perfData.Get())
		os.Exit(3)
	}

	// Check return code, exit with additional info if non-zero
	if contentResult.ReturnCode != 0 {
		printIntro("Web Content Error", h.URL)
		printContentCheck(contentResult.Value)
		if *verbose {
			printVerboseInfo(statuscodeResult.VerboseValue +
				contentResult.VerboseValue)
		}
		fmt.Println(perfData.Get())
		os.Exit(contentResult.ReturnCode)
	}

	// Perform content check, exit with additional info if error
	certResult := h.CheckCertificate(*certwarn, *certcrit)
	if certResult.Error != nil {
		printIntro("TLS Certificate Error", h.URL)
		fmt.Println(statuscodeResult.Error)
		if *verbose {
			printVerboseInfo(statuscodeResult.VerboseValue +
				contentResult.VerboseValue +
				certResult.VerboseValue)
		}
		fmt.Println(perfData.Get())
		os.Exit(3)
	}

	// Check return code, exit with additional info if non-zero
	if certResult.ReturnCode != 0 {
		printIntro("TLS Certificate Error", h.URL)
		printCertCheck(certResult.Value)
		if *verbose {
			printVerboseInfo(statuscodeResult.VerboseValue +
				contentResult.VerboseValue +
				certResult.VerboseValue)
		}
		fmt.Println(perfData.Get())
		os.Exit(certResult.ReturnCode)
	}

	// Print basic info about the checks
	printIntro("OK", h.URL)
	printStatusCode(statuscodeResult.Status, statuscodeResult.Value, *statusCodes)
	printContentCheck(contentResult.Value)
	printCertCheck(certResult.Value)

	// Print verbose info if enabled
	if *verbose {
		printVerboseInfo(statuscodeResult.VerboseValue +
			contentResult.VerboseValue +
			certResult.VerboseValue)
	}

	fmt.Println(perfData.Get())
	os.Exit(0)
}

func printIntro(issue string, url string) {
	fmt.Println(issue + " — HTTPS Check for https://" + url)
}

func printContentCheck(value string) {
	fmt.Println("Content Check: " + value)
}

func printStatusCode(status int, value string, expected string) {
	fmt.Println("Status Code: " + strconv.Itoa(status) + " " + value + ", expected one of: " + expected)
}

func printCertCheck(value string) {
	fmt.Println("Cert Check: " + value)
}

func printVerboseInfo(contents string) {
	if contents != "" {
		fmt.Println("\nAdditional info:\n" + contents)
	}
}
